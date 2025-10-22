/**
 * Congress.gov API Data Ingestion Service
 * 
 * This service handles comprehensive data ingestion from the Congress.gov API,
 * including members, bills, votes, committees, and amendments.
 * 
 * API Documentation: https://api.congress.gov/
 */

import axios, { AxiosInstance } from 'axios';
import { Knex } from 'knex';
import db from '../../database/connection';

interface CongressApiConfig {
  apiKey: string;
  baseUrl: string;
  requestDelay: number; // milliseconds between requests
}

interface Member {
  bioguideId: string;
  name: string;
  party: string;
  state: string;
  district?: string;
  chamber: 'house' | 'senate';
  // Additional fields...
}

interface Bill {
  billId: string;
  congress: number;
  billType: string;
  billNumber: number;
  title: string;
  // Additional fields...
}

interface Vote {
  rollCallNumber: number;
  chamber: 'house' | 'senate';
  congress: number;
  session: number;
  date: string;
  // Additional fields...
}

export class CongressGovIngestionService {
  private client: AxiosInstance;
  private db: Knex;
  private config: CongressApiConfig;

  constructor(apiKey: string) {
    this.config = {
      apiKey: apiKey || process.env.CONGRESS_API_KEY || '',
      baseUrl: 'https://api.congress.gov/v3',
      requestDelay: 100, // Rate limiting: 10 requests/second max
    };

    this.client = axios.create({
      baseURL: this.config.baseUrl,
      params: {
        api_key: this.config.apiKey,
        format: 'json',
      },
      timeout: 30000,
    });

    this.db = db;
  }

  /**
   * Delay between API requests to respect rate limits
   */
  private async delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * Paginated API request handler
   */
  private async paginatedRequest<T>(
    endpoint: string,
    params: Record<string, any> = {}
  ): Promise<T[]> {
    const results: T[] = [];
    let offset = 0;
    const limit = 250; // Max per request

    while (true) {
      try {
        const response = await this.client.get(endpoint, {
          params: { ...params, offset, limit },
        });

        const data = response.data;
        const items = data.members || data.bills || data.votes || data.amendments || [];

        if (items.length === 0) break;

        results.push(...items);
        
        // Check if there are more results
        if (items.length < limit) break;
        
        offset += limit;
        await this.delay(this.config.requestDelay);
      } catch (error) {
        console.error(`Error fetching ${endpoint}:`, error);
        throw error;
      }
    }

    return results;
  }

  /**
   * Ingest all members of Congress (current and historical)
   */
  async ingestMembers(congress: number = 118): Promise<void> {
    console.log(`Starting member ingestion for Congress ${congress}...`);

    try {
      // Fetch House members
      console.log('Fetching House members...');
      const houseMembers = await this.paginatedRequest<any>(
        `/member/congress/${congress}/house`
      );

      // Fetch Senate members
      console.log('Fetching Senate members...');
      const senateMembers = await this.paginatedRequest<any>(
        `/member/congress/${congress}/senate`
      );

      const allMembers = [...houseMembers, ...senateMembers];
      console.log(`Fetched ${allMembers.length} members total`);

      // Process and store each member
      for (const member of allMembers) {
        await this.processMember(member);
        await this.delay(this.config.requestDelay);
      }

      console.log('Member ingestion completed');
    } catch (error) {
      console.error('Error ingesting members:', error);
      throw error;
    }
  }

  /**
   * Process and store individual member data
   */
  private async processMember(memberData: any): Promise<void> {
    try {
      // Fetch detailed member information
      const detailResponse = await this.client.get(
        `/member/${memberData.bioguideId}`
      );
      const detail = detailResponse.data.member;

      const politician = {
        bioguide_id: detail.bioguideId,
        first_name: detail.firstName || detail.name?.split(' ')[0],
        last_name: detail.lastName || detail.name?.split(' ').slice(-1)[0],
        middle_name: detail.middleName,
        suffix: detail.suffix,
        full_name: detail.directOrderName || detail.name,
        party: detail.partyName || detail.party,
        state: detail.state,
        district: detail.district,
        chamber: detail.chamber?.toLowerCase(),
        office: detail.currentMember?.title,
        terms_served: detail.termsInOffice || 0,
        in_office: detail.currentMember || false,
        term_start: detail.termBeginYear,
        term_end: detail.termEndYear,
        image_url: detail.imageUrl || detail.depiction?.imageUrl,
        website: detail.officialWebsiteUrl,
        twitter_account: detail.socialMedia?.find((s: any) => s.type === 'Twitter')?.id,
        facebook_account: detail.socialMedia?.find((s: any) => s.type === 'Facebook')?.id,
        youtube_account: detail.socialMedia?.find((s: any) => s.type === 'YouTube')?.id,
        last_updated: new Date(),
      };

      // Upsert politician
      await this.db('politicians')
        .insert(politician)
        .onConflict('bioguide_id')
        .merge();

      console.log(`Processed member: ${politician.full_name}`);
    } catch (error) {
      console.error(`Error processing member ${memberData.bioguideId}:`, error);
    }
  }

  /**
   * Ingest bills for a specific Congress
   */
  async ingestBills(congress: number = 118): Promise<void> {
    console.log(`Starting bill ingestion for Congress ${congress}...`);

    try {
      const bills = await this.paginatedRequest<any>(`/bill/${congress}`);
      console.log(`Fetched ${bills.length} bills`);

      for (const bill of bills) {
        await this.processBill(bill);
        await this.delay(this.config.requestDelay);
      }

      console.log('Bill ingestion completed');
    } catch (error) {
      console.error('Error ingesting bills:', error);
      throw error;
    }
  }

  /**
   * Process and store individual bill data
   */
  private async processBill(billData: any): Promise<void> {
    try {
      // Fetch detailed bill information
      const detailResponse = await this.client.get(
        `/bill/${billData.congress}/${billData.type}/${billData.number}`
      );
      const detail = detailResponse.data.bill;

      const bill = {
        congress_gov_id: `${detail.congress}-${detail.type}-${detail.number}`,
        bill_type: detail.type,
        bill_number: detail.number,
        congress: detail.congress,
        title: detail.title,
        short_title: detail.titles?.find((t: any) => t.titleType === 'Short')?.title,
        official_title: detail.titles?.find((t: any) => t.titleType === 'Official')?.title,
        introduced_date: detail.introducedDate,
        latest_action_date: detail.latestAction?.actionDate,
        latest_action_text: detail.latestAction?.text,
        sponsor_bioguide_id: detail.sponsors?.[0]?.bioguideId,
        policy_area: detail.policyArea?.name,
        subjects: JSON.stringify(detail.subjects?.legislativeSubjects || []),
        summary: detail.summaries?.[0]?.text,
        status: this.mapBillStatus(detail),
        chamber: detail.originChamber?.toLowerCase(),
        last_updated: new Date(),
      };

      // Upsert bill
      const [billId] = await this.db('bills')
        .insert(bill)
        .onConflict('congress_gov_id')
        .merge()
        .returning('id');

      // Process co-sponsors
      if (detail.cosponsors) {
        await this.processCosponsors(billId.id || billId, detail.cosponsors);
      }

      console.log(`Processed bill: ${bill.congress_gov_id}`);
    } catch (error) {
      console.error(`Error processing bill ${billData.number}:`, error);
    }
  }

  /**
   * Map bill status from various fields
   */
  private mapBillStatus(bill: any): string {
    if (bill.becamePublicLaw) return 'enacted';
    if (bill.vetoDate) return 'vetoed';
    const latestAction = bill.latestAction?.text?.toLowerCase() || '';
    if (latestAction.includes('passed')) return 'passed';
    if (latestAction.includes('failed')) return 'failed';
    if (latestAction.includes('committee')) return 'in_committee';
    return 'introduced';
  }

  /**
   * Process bill co-sponsors
   */
  private async processCosponsors(billId: string, cosponsors: any[]): Promise<void> {
    // This would insert into a bill_cosponsors junction table
    // Implementation depends on your schema
  }

  /**
   * Ingest voting records for a specific Congress and chamber
   */
  async ingestVotes(congress: number = 118, chamber: 'house' | 'senate'): Promise<void> {
    console.log(`Starting vote ingestion for ${chamber} in Congress ${congress}...`);

    try {
      const votes = await this.paginatedRequest<any>(`/vote/${congress}/${chamber}`);
      console.log(`Fetched ${votes.length} votes`);

      for (const vote of votes) {
        await this.processVote(vote, chamber);
        await this.delay(this.config.requestDelay);
      }

      console.log('Vote ingestion completed');
    } catch (error) {
      console.error('Error ingesting votes:', error);
      throw error;
    }
  }

  /**
   * Process and store individual vote data
   */
  private async processVote(voteData: any, chamber: string): Promise<void> {
    try {
      const detailResponse = await this.client.get(
        `/vote/${voteData.congress}/${chamber}/${voteData.rollCall}`
      );
      const detail = detailResponse.data.vote;

      const vote = {
        congress: detail.congress,
        chamber: chamber,
        session: detail.session,
        roll_call_number: detail.rollCall,
        vote_date: detail.date,
        vote_time: detail.time,
        question: detail.question,
        description: detail.description,
        bill_id: detail.bill?.congress && detail.bill?.type && detail.bill?.number
          ? `${detail.bill.congress}-${detail.bill.type}-${detail.bill.number}`
          : null,
        result: detail.result,
        vote_type: detail.voteType,
        yea_count: detail.yeas || 0,
        nay_count: detail.nays || 0,
        present_count: detail.present || 0,
        not_voting_count: detail.notVoting || 0,
        last_updated: new Date(),
      };

      // Upsert vote
      const [voteId] = await this.db('votes')
        .insert(vote)
        .onConflict(['congress', 'chamber', 'roll_call_number'])
        .merge()
        .returning('id');

      // Process individual member votes
      if (detail.members) {
        await this.processMemberVotes(voteId.id || voteId, detail.members);
      }

      console.log(`Processed vote: ${chamber} roll ${vote.roll_call_number}`);
    } catch (error) {
      console.error(`Error processing vote ${voteData.rollCall}:`, error);
    }
  }

  /**
   * Process individual member votes
   */
  private async processMemberVotes(voteId: string, members: any[]): Promise<void> {
    // This would insert into a member_votes junction table
    // Implementation depends on your schema
  }

  /**
   * Ingest committee data
   */
  async ingestCommittees(congress: number = 118): Promise<void> {
    console.log(`Starting committee ingestion for Congress ${congress}...`);

    try {
      // Fetch House committees
      const houseCommittees = await this.paginatedRequest<any>(
        `/committee/${congress}/house`
      );

      // Fetch Senate committees
      const senateCommittees = await this.paginatedRequest<any>(
        `/committee/${congress}/senate`
      );

      const allCommittees = [...houseCommittees, ...senateCommittees];
      console.log(`Fetched ${allCommittees.length} committees`);

      for (const committee of allCommittees) {
        await this.processCommittee(committee);
        await this.delay(this.config.requestDelay);
      }

      console.log('Committee ingestion completed');
    } catch (error) {
      console.error('Error ingesting committees:', error);
      throw error;
    }
  }

  /**
   * Process and store committee data
   */
  private async processCommittee(committeeData: any): Promise<void> {
    // Implementation for storing committee data
    // Would require a committees table
  }

  /**
   * Run comprehensive ingestion for current Congress
   */
  async runFullIngestion(congress: number = 118): Promise<void> {
    console.log(`Starting full ingestion for Congress ${congress}...`);

    try {
      // Ingest in order of dependencies
      await this.ingestMembers(congress);
      await this.ingestCommittees(congress);
      await this.ingestBills(congress);
      await this.ingestVotes(congress, 'house');
      await this.ingestVotes(congress, 'senate');

      console.log('Full ingestion completed successfully!');
    } catch (error) {
      console.error('Error during full ingestion:', error);
      throw error;
    }
  }

  /**
   * Incremental update - fetch only new/updated records
   */
  async incrementalUpdate(congress: number = 118, since?: Date): Promise<void> {
    const sinceDate = since || new Date(Date.now() - 24 * 60 * 60 * 1000); // Last 24 hours
    console.log(`Running incremental update since ${sinceDate.toISOString()}...`);

    // Implementation for fetching only updated records
    // Would use fromDateTime parameter in API calls
  }
}

// Export singleton instance
export default new CongressGovIngestionService(process.env.CONGRESS_API_KEY || '');
