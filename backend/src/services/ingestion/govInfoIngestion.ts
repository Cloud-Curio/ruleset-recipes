/**
 * GovInfo.gov API Data Ingestion Service
 * 
 * This service handles comprehensive data ingestion from the GovInfo.gov API,
 * including Congressional Record, Federal Register, compiled bills, and bulk data.
 * 
 * API Documentation: https://api.govinfo.gov/docs/
 */

import axios, { AxiosInstance } from 'axios';
import { Knex } from 'knex';
import db from '../../database/connection';
import * as cheerio from 'cheerio';
import { parseStringPromise } from 'xml2js';

interface GovInfoApiConfig {
  apiKey: string;
  baseUrl: string;
  requestDelay: number;
}

interface CongressionalRecordItem {
  packageId: string;
  title: string;
  dateIssued: string;
  docClass: string;
  congress?: number;
  session?: number;
}

interface FederalRegisterDocument {
  documentNumber: string;
  title: string;
  publicationDate: string;
  agencyNames: string[];
  documentType: string;
}

export class GovInfoIngestionService {
  private client: AxiosInstance;
  private db: Knex;
  private config: GovInfoApiConfig;

  constructor(apiKey: string) {
    this.config = {
      apiKey: apiKey || process.env.GOVINFO_API_KEY || '',
      baseUrl: 'https://api.govinfo.gov',
      requestDelay: 100, // Rate limiting
    };

    this.client = axios.create({
      baseURL: this.config.baseUrl,
      params: {
        api_key: this.config.apiKey,
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
   * Ingest Congressional Record for a date range
   */
  async ingestCongressionalRecord(
    startDate: Date,
    endDate: Date = new Date()
  ): Promise<void> {
    console.log(`Ingesting Congressional Record from ${startDate.toISOString()} to ${endDate.toISOString()}...`);

    try {
      const response = await this.client.get('/collections/CREC', {
        params: {
          offsetMark: '*',
          pageSize: 100,
          startDate: this.formatDate(startDate),
          endDate: this.formatDate(endDate),
        },
      });

      const packages = response.data.packages || [];
      console.log(`Found ${packages.length} Congressional Record items`);

      for (const pkg of packages) {
        await this.processCongressionalRecord(pkg);
        await this.delay(this.config.requestDelay);
      }

      console.log('Congressional Record ingestion completed');
    } catch (error) {
      console.error('Error ingesting Congressional Record:', error);
      throw error;
    }
  }

  /**
   * Process individual Congressional Record item
   */
  private async processCongressionalRecord(packageData: any): Promise<void> {
    try {
      const packageId = packageData.packageId;
      
      // Fetch full package details
      const detailResponse = await this.client.get(`/packages/${packageId}/summary`);
      const detail = detailResponse.data;

      // Extract text content
      let textContent = '';
      try {
        const textResponse = await this.client.get(`/packages/${packageId}/htm`, {
          responseType: 'text',
        });
        const $ = cheerio.load(textResponse.data);
        textContent = $('body').text().trim();
      } catch (error) {
        console.warn(`Could not fetch text for ${packageId}:`, error);
      }

      const record = {
        package_id: packageId,
        title: detail.title,
        date_issued: detail.dateIssued,
        doc_class: detail.docClass,
        collection: 'CREC',
        congress: detail.congress,
        session: detail.session,
        pages: detail.pages,
        government_author: JSON.stringify(detail.governmentAuthor || []),
        content_text: textContent.substring(0, 50000), // Limit to 50k chars
        download_url: detail.download?.txtLink,
        pdf_url: detail.download?.pdfLink,
        last_modified: detail.lastModified,
        last_updated: new Date(),
      };

      // Store in congressional_records table (would need to create this table)
      await this.db('congressional_records')
        .insert(record)
        .onConflict('package_id')
        .merge();

      console.log(`Processed Congressional Record: ${packageId}`);
    } catch (error) {
      console.error(`Error processing Congressional Record ${packageData.packageId}:`, error);
    }
  }

  /**
   * Ingest Federal Register documents
   */
  async ingestFederalRegister(
    startDate: Date,
    endDate: Date = new Date()
  ): Promise<void> {
    console.log(`Ingesting Federal Register from ${startDate.toISOString()} to ${endDate.toISOString()}...`);

    try {
      const response = await this.client.get('/collections/FR', {
        params: {
          offsetMark: '*',
          pageSize: 100,
          startDate: this.formatDate(startDate),
          endDate: this.formatDate(endDate),
        },
      });

      const packages = response.data.packages || [];
      console.log(`Found ${packages.length} Federal Register documents`);

      for (const pkg of packages) {
        await this.processFederalRegister(pkg);
        await this.delay(this.config.requestDelay);
      }

      console.log('Federal Register ingestion completed');
    } catch (error) {
      console.error('Error ingesting Federal Register:', error);
      throw error;
    }
  }

  /**
   * Process individual Federal Register document
   */
  private async processFederalRegister(packageData: any): Promise<void> {
    try {
      const packageId = packageData.packageId;
      
      const detailResponse = await this.client.get(`/packages/${packageId}/summary`);
      const detail = detailResponse.data;

      // Extract metadata
      const record = {
        package_id: packageId,
        title: detail.title,
        publication_date: detail.dateIssued,
        volume: detail.volume,
        issue_number: detail.issue,
        collection: 'FR',
        doc_class: detail.docClass,
        pages: detail.pages,
        agencies: JSON.stringify(detail.governmentAuthor || []),
        subjects: JSON.stringify(detail.subjects || []),
        download_url: detail.download?.txtLink,
        pdf_url: detail.download?.pdfLink,
        last_modified: detail.lastModified,
        last_updated: new Date(),
      };

      await this.db('federal_register_documents')
        .insert(record)
        .onConflict('package_id')
        .merge();

      console.log(`Processed Federal Register: ${packageId}`);
    } catch (error) {
      console.error(`Error processing Federal Register ${packageData.packageId}:`, error);
    }
  }

  /**
   * Ingest compiled bills (bills as enacted)
   */
  async ingestCompiledBills(congress: number = 118): Promise<void> {
    console.log(`Ingesting compiled bills for Congress ${congress}...`);

    try {
      const response = await this.client.get('/collections/BILLS', {
        params: {
          congress: congress,
          offsetMark: '*',
          pageSize: 100,
        },
      });

      const packages = response.data.packages || [];
      console.log(`Found ${packages.length} compiled bills`);

      for (const pkg of packages) {
        await this.processCompiledBill(pkg);
        await this.delay(this.config.requestDelay);
      }

      console.log('Compiled bills ingestion completed');
    } catch (error) {
      console.error('Error ingesting compiled bills:', error);
      throw error;
    }
  }

  /**
   * Process individual compiled bill
   */
  private async processCompiledBill(packageData: any): Promise<void> {
    try {
      const packageId = packageData.packageId;
      
      const detailResponse = await this.client.get(`/packages/${packageId}/summary`);
      const detail = detailResponse.data;

      // Extract full bill text
      let billText = '';
      try {
        const textResponse = await this.client.get(`/packages/${packageId}/htm`, {
          responseType: 'text',
        });
        const $ = cheerio.load(textResponse.data);
        billText = $('body').text().trim();
      } catch (error) {
        console.warn(`Could not fetch text for bill ${packageId}`);
      }

      // Parse package ID to extract bill info
      // Format: BILLS-118hr1500ih
      const billMatch = packageId.match(/BILLS-(\d+)([a-z]+)(\d+)([a-z]+)/i);
      
      if (billMatch) {
        const [, congress, type, number, version] = billMatch;

        const compiledBill = {
          package_id: packageId,
          congress: parseInt(congress),
          bill_type: type.toUpperCase(),
          bill_number: parseInt(number),
          version: version,
          title: detail.title,
          date_issued: detail.dateIssued,
          full_text: billText.substring(0, 100000), // Limit to 100k chars
          pages: detail.pages,
          download_url: detail.download?.txtLink,
          pdf_url: detail.download?.pdfLink,
          xml_url: detail.download?.xmlLink,
          last_modified: detail.lastModified,
          last_updated: new Date(),
        };

        await this.db('compiled_bills')
          .insert(compiledBill)
          .onConflict('package_id')
          .merge();

        console.log(`Processed compiled bill: ${packageId}`);
      }
    } catch (error) {
      console.error(`Error processing compiled bill ${packageData.packageId}:`, error);
    }
  }

  /**
   * Ingest committee reports
   */
  async ingestCommitteeReports(congress: number = 118): Promise<void> {
    console.log(`Ingesting committee reports for Congress ${congress}...`);

    try {
      const response = await this.client.get('/collections/CRPT', {
        params: {
          congress: congress,
          offsetMark: '*',
          pageSize: 100,
        },
      });

      const packages = response.data.packages || [];
      console.log(`Found ${packages.length} committee reports`);

      for (const pkg of packages) {
        await this.processCommitteeReport(pkg);
        await this.delay(this.config.requestDelay);
      }

      console.log('Committee reports ingestion completed');
    } catch (error) {
      console.error('Error ingesting committee reports:', error);
      throw error;
    }
  }

  /**
   * Process individual committee report
   */
  private async processCommitteeReport(packageData: any): Promise<void> {
    try {
      const packageId = packageData.packageId;
      
      const detailResponse = await this.client.get(`/packages/${packageId}/summary`);
      const detail = detailResponse.data;

      const report = {
        package_id: packageId,
        congress: detail.congress,
        chamber: detail.chamber,
        report_type: detail.reportType,
        report_number: detail.reportNumber,
        title: detail.title,
        date_issued: detail.dateIssued,
        committee: detail.committee,
        pages: detail.pages,
        download_url: detail.download?.txtLink,
        pdf_url: detail.download?.pdfLink,
        last_modified: detail.lastModified,
        last_updated: new Date(),
      };

      await this.db('committee_reports')
        .insert(report)
        .onConflict('package_id')
        .merge();

      console.log(`Processed committee report: ${packageId}`);
    } catch (error) {
      console.error(`Error processing committee report ${packageData.packageId}:`, error);
    }
  }

  /**
   * Ingest hearing transcripts
   */
  async ingestHearings(congress: number = 118): Promise<void> {
    console.log(`Ingesting hearings for Congress ${congress}...`);

    try {
      const response = await this.client.get('/collections/CHRG', {
        params: {
          congress: congress,
          offsetMark: '*',
          pageSize: 100,
        },
      });

      const packages = response.data.packages || [];
      console.log(`Found ${packages.length} hearings`);

      for (const pkg of packages) {
        await this.processHearing(pkg);
        await this.delay(this.config.requestDelay);
      }

      console.log('Hearings ingestion completed');
    } catch (error) {
      console.error('Error ingesting hearings:', error);
      throw error;
    }
  }

  /**
   * Process individual hearing transcript
   */
  private async processHearing(packageData: any): Promise<void> {
    try {
      const packageId = packageData.packageId;
      
      const detailResponse = await this.client.get(`/packages/${packageId}/summary`);
      const detail = detailResponse.data;

      // Extract transcript text
      let transcript = '';
      try {
        const textResponse = await this.client.get(`/packages/${packageId}/htm`, {
          responseType: 'text',
        });
        const $ = cheerio.load(textResponse.data);
        transcript = $('body').text().trim();
      } catch (error) {
        console.warn(`Could not fetch transcript for hearing ${packageId}`);
      }

      const hearing = {
        package_id: packageId,
        congress: detail.congress,
        chamber: detail.chamber,
        title: detail.title,
        date_issued: detail.dateIssued,
        committee: detail.committee,
        jacketNumber: detail.jacketNumber,
        transcript_text: transcript.substring(0, 100000), // Limit size
        pages: detail.pages,
        download_url: detail.download?.txtLink,
        pdf_url: detail.download?.pdfLink,
        last_modified: detail.lastModified,
        last_updated: new Date(),
      };

      await this.db('hearings')
        .insert(hearing)
        .onConflict('package_id')
        .merge();

      console.log(`Processed hearing: ${packageId}`);
    } catch (error) {
      console.error(`Error processing hearing ${packageData.packageId}:`, error);
    }
  }

  /**
   * Download bulk data files
   */
  async downloadBulkData(collection: string, year?: number): Promise<void> {
    console.log(`Downloading bulk data for collection: ${collection}...`);

    try {
      // Get list of available bulk data files
      const response = await this.client.get(`/collections/${collection}/bulk`, {
        params: year ? { year } : {},
      });

      const files = response.data.files || [];
      console.log(`Found ${files.length} bulk data files`);

      // Process each bulk data file
      for (const file of files) {
        console.log(`Processing bulk file: ${file.link}`);
        // Implementation would download and process the bulk XML/JSON files
        // This could be stored or processed in batches
        await this.delay(this.config.requestDelay);
      }

      console.log('Bulk data download completed');
    } catch (error) {
      console.error('Error downloading bulk data:', error);
      throw error;
    }
  }

  /**
   * Format date for API requests
   */
  private formatDate(date: Date): string {
    return date.toISOString().split('T')[0];
  }

  /**
   * Run comprehensive ingestion
   */
  async runFullIngestion(congress: number = 118): Promise<void> {
    console.log(`Starting full GovInfo ingestion for Congress ${congress}...`);

    const startDate = new Date(2023, 0, 1); // Adjust as needed
    const endDate = new Date();

    try {
      await this.ingestCongressionalRecord(startDate, endDate);
      await this.ingestFederalRegister(startDate, endDate);
      await this.ingestCompiledBills(congress);
      await this.ingestCommitteeReports(congress);
      await this.ingestHearings(congress);

      console.log('Full GovInfo ingestion completed successfully!');
    } catch (error) {
      console.error('Error during full ingestion:', error);
      throw error;
    }
  }
}

// Export singleton instance
export default new GovInfoIngestionService(process.env.GOVINFO_API_KEY || '');
