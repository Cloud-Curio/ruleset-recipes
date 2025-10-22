/**
 * OpenStates API Data Ingestion Service
 * 
 * This service handles comprehensive data ingestion from the OpenStates API,
 * including state legislators, bills, votes, and committees.
 * 
 * API Documentation: https://docs.openstates.org/api-v3/
 */

import axios, { AxiosInstance } from 'axios';
import { Knex } from 'knex';
import db from '../../database/connection';

interface OpenStatesApiConfig {
  apiKey: string;
  baseUrl: string;
  requestDelay: number;
}

interface StateLegislator {
  id: string;
  name: string;
  party: string;
  jurisdiction: string;
  currentRole: any;
}

interface StateBill {
  id: string;
  identifier: string;
  title: string;
  classification: string[];
  jurisdiction: string;
}

export class OpenStatesIngestionService {
  private client: AxiosInstance;
  private db: Knex;
  private config: OpenStatesApiConfig;

  // All US states and territories
  private jurisdictions = [
    'al', 'ak', 'az', 'ar', 'ca', 'co', 'ct', 'de', 'fl', 'ga',
    'hi', 'id', 'il', 'in', 'ia', 'ks', 'ky', 'la', 'me', 'md',
    'ma', 'mi', 'mn', 'ms', 'mo', 'mt', 'ne', 'nv', 'nh', 'nj',
    'nm', 'ny', 'nc', 'nd', 'oh', 'ok', 'or', 'pa', 'ri', 'sc',
    'sd', 'tn', 'tx', 'ut', 'vt', 'va', 'wa', 'wv', 'wi', 'wy',
    'dc', 'pr', // DC and Puerto Rico
  ];

  constructor(apiKey: string) {
    this.config = {
      apiKey: apiKey || process.env.OPENSTATES_API_KEY || '',
      baseUrl: 'https://v3.openstates.org',
      requestDelay: 200, // Be conservative with rate limiting
    };

    this.client = axios.create({
      baseURL: this.config.baseUrl,
      headers: {
        'X-API-KEY': this.config.apiKey,
      },
      timeout: 30000,
    });

    this.db = db;
  }

  /**
   * Delay between API requests
   */
  private async delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * Paginated request handler
   */
  private async paginatedRequest<T>(
    endpoint: string,
    params: Record<string, any> = {}
  ): Promise<T[]> {
    const results: T[] = [];
    let page = 1;
    const perPage = 100;

    while (true) {
      try {
        const response = await this.client.get(endpoint, {
          params: { ...params, page, per_page: perPage },
        });

        const data = response.data;
        const items = data.results || data;

        if (!items || items.length === 0) break;

        results.push(...items);

        // Check if there are more pages
        if (!data.pagination || !data.pagination.next) break;
        if (items.length < perPage) break;

        page++;
        await this.delay(this.config.requestDelay);
      } catch (error) {
        console.error(`Error fetching ${endpoint}:`, error);
        throw error;
      }
    }

    return results;
  }

  /**
   * Ingest all state legislators
   */
  async ingestStateLegislators(jurisdiction?: string): Promise<void> {
    const jurisdictionsToProcess = jurisdiction ? [jurisdiction] : this.jurisdictions;

    console.log(`Starting state legislator ingestion for ${jurisdictionsToProcess.length} jurisdictions...`);

    for (const jur of jurisdictionsToProcess) {
      try {
        console.log(`Fetching legislators for ${jur.toUpperCase()}...`);
        
        const legislators = await this.paginatedRequest<any>('/people', {
          jurisdiction: jur,
          per_page: 100,
        });

        console.log(`Found ${legislators.length} legislators in ${jur.toUpperCase()}`);

        for (const legislator of legislators) {
          await this.processStateLegislator(legislator, jur);
          await this.delay(this.config.requestDelay);
        }
      } catch (error) {
        console.error(`Error ingesting legislators for ${jur}:`, error);
      }
    }

    console.log('State legislator ingestion completed');
  }

  /**
   * Process and store individual state legislator
   */
  private async processStateLegislator(legislatorData: any, jurisdiction: string): Promise<void> {
    try {
      // Fetch detailed legislator information
      const detailResponse = await this.client.get(`/people/${legislatorData.id}`);
      const detail = detailResponse.data;

      const politician = {
        openstates_id: detail.id,
        first_name: detail.given_name || detail.name.split(' ')[0],
        last_name: detail.family_name || detail.name.split(' ').slice(-1)[0],
        full_name: detail.name,
        party: detail.current_role?.party || detail.party?.[0]?.name,
        state: jurisdiction.toUpperCase(),
        district: detail.current_role?.district,
        chamber: detail.current_role?.chamber,
        office: detail.current_role?.title,
        image_url: detail.image,
        email: detail.email,
        capitol_address: detail.capitol_address?.address,
        capitol_phone: detail.capitol_address?.voice,
        website: detail.links?.[0]?.url,
        in_office: detail.current_role ? true : false,
        extras: JSON.stringify(detail.extras || {}),
        sources: JSON.stringify(detail.sources || []),
        last_updated: new Date(),
      };

      // Upsert politician
      await this.db('politicians')
        .insert(politician)
        .onConflict('openstates_id')
        .merge();

      console.log(`Processed state legislator: ${politician.full_name} (${jurisdiction.toUpperCase()})`);
    } catch (error) {
      console.error(`Error processing legislator ${legislatorData.id}:`, error);
    }
  }

  /**
   * Ingest state bills
   */
  async ingestStateBills(
    jurisdiction: string,
    session?: string,
    updatedSince?: Date
  ): Promise<void> {
    console.log(`Ingesting bills for ${jurisdiction.toUpperCase()}...`);

    try {
      const params: any = {
        jurisdiction: jurisdiction,
      };

      if (session) params.session = session;
      if (updatedSince) {
        params.updated_since = updatedSince.toISOString().split('T')[0];
      }

      const bills = await this.paginatedRequest<any>('/bills', params);
      console.log(`Found ${bills.length} bills for ${jurisdiction.toUpperCase()}`);

      for (const bill of bills) {
        await this.processStateBill(bill);
        await this.delay(this.config.requestDelay);
      }

      console.log(`Bill ingestion completed for ${jurisdiction.toUpperCase()}`);
    } catch (error) {
      console.error(`Error ingesting bills for ${jurisdiction}:`, error);
      throw error;
    }
  }

  /**
   * Process and store individual state bill
   */
  private async processStateBill(billData: any): Promise<void> {
    try {
      // Fetch detailed bill information
      const detailResponse = await this.client.get(`/bills/${billData.id}`);
      const detail = detailResponse.data;

      const bill = {
        openstates_id: detail.id,
        identifier: detail.identifier,
        title: detail.title,
        classification: JSON.stringify(detail.classification || []),
        subject: JSON.stringify(detail.subject || []),
        jurisdiction: detail.jurisdiction,
        session: detail.session,
        from_organization: detail.from_organization,
        sponsorships: JSON.stringify(detail.sponsorships || []),
        abstracts: JSON.stringify(detail.abstracts || []),
        other_titles: JSON.stringify(detail.other_titles || []),
        actions: JSON.stringify(detail.actions || []),
        latest_action_date: detail.actions?.[0]?.date,
        latest_action_description: detail.actions?.[0]?.description,
        created_at: detail.created_at,
        updated_at: detail.updated_at,
        first_action_date: detail.first_action_date,
        latest_passage_date: detail.latest_passage_date,
        versions: JSON.stringify(detail.versions || []),
        documents: JSON.stringify(detail.documents || []),
        sources: JSON.stringify(detail.sources || []),
        last_updated: new Date(),
      };

      // Upsert bill
      const [billId] = await this.db('state_bills')
        .insert(bill)
        .onConflict('openstates_id')
        .merge()
        .returning('id');

      // Process votes if available
      if (detail.votes && detail.votes.length > 0) {
        for (const vote of detail.votes) {
          await this.processStateVote(vote, billId.id || billId);
        }
      }

      console.log(`Processed state bill: ${bill.identifier}`);
    } catch (error) {
      console.error(`Error processing bill ${billData.id}:`, error);
    }
  }

  /**
   * Process and store state vote
   */
  private async processStateVote(voteData: any, billId: string): Promise<void> {
    try {
      const vote = {
        openstates_id: voteData.id,
        bill_id: billId,
        identifier: voteData.identifier,
        motion_text: voteData.motion_text,
        motion_classification: JSON.stringify(voteData.motion_classification || []),
        start_date: voteData.start_date,
        result: voteData.result,
        organization: voteData.organization,
        yes_count: voteData.counts?.find((c: any) => c.option === 'yes')?.value || 0,
        no_count: voteData.counts?.find((c: any) => c.option === 'no')?.value || 0,
        other_count: voteData.counts?.filter((c: any) => c.option !== 'yes' && c.option !== 'no')
          .reduce((sum: number, c: any) => sum + c.value, 0) || 0,
        votes: JSON.stringify(voteData.votes || []),
        sources: JSON.stringify(voteData.sources || []),
        last_updated: new Date(),
      };

      await this.db('state_votes')
        .insert(vote)
        .onConflict('openstates_id')
        .merge();

      console.log(`Processed state vote: ${vote.identifier || vote.openstates_id}`);
    } catch (error) {
      console.error(`Error processing vote ${voteData.id}:`, error);
    }
  }

  /**
   * Ingest state committees
   */
  async ingestStateCommittees(jurisdiction: string): Promise<void> {
    console.log(`Ingesting committees for ${jurisdiction.toUpperCase()}...`);

    try {
      // OpenStates v3 uses organizations endpoint for committees
      const committees = await this.paginatedRequest<any>('/organizations', {
        jurisdiction: jurisdiction,
        classification: 'committee',
      });

      console.log(`Found ${committees.length} committees for ${jurisdiction.toUpperCase()}`);

      for (const committee of committees) {
        await this.processStateCommittee(committee, jurisdiction);
        await this.delay(this.config.requestDelay);
      }

      console.log(`Committee ingestion completed for ${jurisdiction.toUpperCase()}`);
    } catch (error) {
      console.error(`Error ingesting committees for ${jurisdiction}:`, error);
      throw error;
    }
  }

  /**
   * Process and store state committee
   */
  private async processStateCommittee(committeeData: any, jurisdiction: string): Promise<void> {
    try {
      const committee = {
        openstates_id: committeeData.id,
        name: committeeData.name,
        classification: committeeData.classification,
        jurisdiction: jurisdiction,
        chamber: committeeData.chamber,
        parent_id: committeeData.parent_id,
        memberships: JSON.stringify(committeeData.memberships || []),
        sources: JSON.stringify(committeeData.sources || []),
        last_updated: new Date(),
      };

      await this.db('state_committees')
        .insert(committee)
        .onConflict('openstates_id')
        .merge();

      console.log(`Processed state committee: ${committee.name}`);
    } catch (error) {
      console.error(`Error processing committee ${committeeData.id}:`, error);
    }
  }

  /**
   * Ingest all data for a specific state
   */
  async ingestFullState(jurisdiction: string, session?: string): Promise<void> {
    console.log(`Starting full ingestion for ${jurisdiction.toUpperCase()}...`);

    try {
      await this.ingestStateLegislators(jurisdiction);
      await this.ingestStateCommittees(jurisdiction);
      await this.ingestStateBills(jurisdiction, session);

      console.log(`Full ingestion completed for ${jurisdiction.toUpperCase()}`);
    } catch (error) {
      console.error(`Error during full state ingestion for ${jurisdiction}:`, error);
      throw error;
    }
  }

  /**
   * Run comprehensive ingestion for all states
   */
  async runFullIngestion(): Promise<void> {
    console.log('Starting full OpenStates ingestion for all jurisdictions...');

    for (const jurisdiction of this.jurisdictions) {
      try {
        await this.ingestFullState(jurisdiction);
        // Add a longer delay between states to be respectful of API limits
        await this.delay(1000);
      } catch (error) {
        console.error(`Error ingesting ${jurisdiction}:`, error);
        // Continue with next jurisdiction even if one fails
      }
    }

    console.log('Full OpenStates ingestion completed!');
  }

  /**
   * Incremental update - fetch only recent updates
   */
  async incrementalUpdate(jurisdiction?: string, daysSince: number = 1): Promise<void> {
    const jurisdictionsToUpdate = jurisdiction ? [jurisdiction] : this.jurisdictions;
    const since = new Date(Date.now() - daysSince * 24 * 60 * 60 * 1000);

    console.log(`Running incremental update since ${since.toISOString()}...`);

    for (const jur of jurisdictionsToUpdate) {
      try {
        console.log(`Updating ${jur.toUpperCase()}...`);
        await this.ingestStateBills(jur, undefined, since);
        await this.delay(1000);
      } catch (error) {
        console.error(`Error updating ${jur}:`, error);
      }
    }

    console.log('Incremental update completed');
  }

  /**
   * Get current sessions for all jurisdictions
   */
  async fetchCurrentSessions(): Promise<Map<string, string>> {
    const sessions = new Map<string, string>();

    for (const jurisdiction of this.jurisdictions) {
      try {
        const response = await this.client.get(`/jurisdictions/${jurisdiction}`);
        const currentSession = response.data.legislative_sessions?.find(
          (s: any) => s.active
        );
        if (currentSession) {
          sessions.set(jurisdiction, currentSession.identifier);
        }
        await this.delay(this.config.requestDelay);
      } catch (error) {
        console.error(`Error fetching sessions for ${jurisdiction}:`, error);
      }
    }

    return sessions;
  }
}

// Export singleton instance
export default new OpenStatesIngestionService(process.env.OPENSTATES_API_KEY || '');
