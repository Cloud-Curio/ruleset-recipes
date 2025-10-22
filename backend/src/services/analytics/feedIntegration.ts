/**
 * Feed Integration Service
 * 
 * Combines data from multiple sources (voting records, social media, bills, etc.)
 * into a unified feed for display on politician profiles ("Fakebook").
 */

import { Knex } from 'knex';
import db from '../../database/connection';

interface FeedItem {
  id: string;
  type: 'vote' | 'bill' | 'social_post' | 'statement' | 'alert' | 'achievement';
  timestamp: Date;
  politician: {
    id: string;
    name: string;
    party: string;
    state: string;
    avatar: string;
  };
  content: any;
  metadata: {
    engagementCount?: number;
    sentiment?: number;
    toxicity?: number;
    importance?: number;
  };
  analytics?: {
    consistencyScore?: number;
    alignmentWithPreviousPositions?: number;
    publicReaction?: number;
  };
}

interface FeedFilter {
  politicianIds?: string[];
  states?: string[];
  parties?: string[];
  types?: string[];
  topics?: string[];
  startDate?: Date;
  endDate?: Date;
  minImportance?: number;
  onlyControversial?: boolean;
}

export class FeedIntegrationService {
  private db: Knex;

  constructor() {
    this.db = db;
  }

  /**
   * Get unified feed for a politician or set of politicians
   */
  async getFeed(
    filter: FeedFilter,
    page: number = 1,
    pageSize: number = 50
  ): Promise<{ items: FeedItem[]; total: number; hasMore: boolean }> {
    const offset = (page - 1) * pageSize;

    // Fetch different feed item types in parallel
    const [votes, bills, socialPosts, alerts] = await Promise.all([
      this.getVoteFeedItems(filter, offset, pageSize),
      this.getBillFeedItems(filter, offset, pageSize),
      this.getSocialPostFeedItems(filter, offset, pageSize),
      this.getAlertFeedItems(filter, offset, pageSize)
    ]);

    // Combine and sort by timestamp
    const allItems = [...votes, ...bills, ...socialPosts, ...alerts]
      .sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime())
      .slice(0, pageSize);

    const total = await this.getTotalFeedItems(filter);

    return {
      items: allItems,
      total,
      hasMore: offset + pageSize < total
    };
  }

  /**
   * Get vote-related feed items
   */
  private async getVoteFeedItems(
    filter: FeedFilter,
    offset: number,
    limit: number
  ): Promise<FeedItem[]> {
    let query = this.db('votes as v')
      .join('politicians as p', 'v.politician_id', 'p.id')
      .leftJoin('bills as b', 'v.bill_id', 'b.id')
      .select(
        'v.id',
        'v.vote_position',
        'v.vote_date as timestamp',
        'v.question',
        'v.result',
        'v.party_line_vote',
        'b.title as bill_title',
        'b.congress_bill_id',
        'b.policy_areas',
        'b.controversy_score',
        'p.id as politician_id',
        'p.full_name',
        'p.party',
        'p.state',
        'p.image_url'
      );

    query = this.applyFilters(query, filter, 'v', 'p');

    if (filter.types && !filter.types.includes('vote')) {
      return [];
    }

    const results = await query
      .orderBy('v.vote_date', 'desc')
      .limit(limit)
      .offset(offset);

    return results.map(row => ({
      id: `vote-${row.id}`,
      type: 'vote' as const,
      timestamp: new Date(row.timestamp),
      politician: {
        id: row.politician_id,
        name: row.full_name,
        party: row.party,
        state: row.state,
        avatar: row.image_url
      },
      content: {
        votePosition: row.vote_position,
        question: row.question,
        result: row.result,
        billTitle: row.bill_title,
        billId: row.congress_bill_id,
        partyLineVote: row.party_line_vote,
        policyAreas: JSON.parse(row.policy_areas || '[]')
      },
      metadata: {
        importance: row.controversy_score ? row.controversy_score / 100 : 0.5
      },
      analytics: {
        alignmentWithPreviousPositions: await this.calculateVoteAlignment(
          row.politician_id,
          row.id
        )
      }
    }));
  }

  /**
   * Get bill-related feed items
   */
  private async getBillFeedItems(
    filter: FeedFilter,
    offset: number,
    limit: number
  ): Promise<FeedItem[]> {
    let query = this.db('bills as b')
      .join('politicians as p', 'b.sponsor_id', 'p.id')
      .select(
        'b.id',
        'b.congress_bill_id',
        'b.title',
        'b.short_title',
        'b.status',
        'b.policy_areas',
        'b.introduced_date as timestamp',
        'b.summary',
        'b.controversy_score',
        'b.bipartisan_support',
        'p.id as politician_id',
        'p.full_name',
        'p.party',
        'p.state',
        'p.image_url'
      );

    query = this.applyFilters(query, filter, 'b', 'p');

    if (filter.types && !filter.types.includes('bill')) {
      return [];
    }

    const results = await query
      .orderBy('b.introduced_date', 'desc')
      .limit(limit)
      .offset(offset);

    return results.map(row => ({
      id: `bill-${row.id}`,
      type: 'bill' as const,
      timestamp: new Date(row.timestamp),
      politician: {
        id: row.politician_id,
        name: row.full_name,
        party: row.party,
        state: row.state,
        avatar: row.image_url
      },
      content: {
        billId: row.congress_bill_id,
        title: row.title,
        shortTitle: row.short_title,
        status: row.status,
        summary: row.summary,
        policyAreas: JSON.parse(row.policy_areas || '[]'),
        bipartisanSupport: row.bipartisan_support
      },
      metadata: {
        importance: row.controversy_score ? row.controversy_score / 100 : 0.5
      }
    }));
  }

  /**
   * Get social media post feed items
   */
  private async getSocialPostFeedItems(
    filter: FeedFilter,
    offset: number,
    limit: number
  ): Promise<FeedItem[]> {
    let query = this.db('social_media_posts as smp')
      .join('politicians as p', 'smp.politician_id', 'p.id')
      .leftJoin('nlp_analysis as nlp', function() {
        this.on('nlp.content_id', '=', 'smp.id')
          .andOn(this.db.raw("nlp.content_type = 'social_post'"));
      })
      .select(
        'smp.id',
        'smp.platform',
        'smp.content',
        'smp.posted_at as timestamp',
        'smp.likes_count',
        'smp.shares_count',
        'smp.comments_count',
        'smp.post_url',
        'nlp.sentiment_score',
        'nlp.toxicity_score',
        'nlp.contains_hate_speech',
        'nlp.topics',
        'p.id as politician_id',
        'p.full_name',
        'p.party',
        'p.state',
        'p.image_url'
      )
      .where('smp.is_deleted', false);

    query = this.applyFilters(query, filter, 'smp', 'p');

    if (filter.types && !filter.types.includes('social_post')) {
      return [];
    }

    // Filter controversial posts if requested
    if (filter.onlyControversial) {
      query = query.where('nlp.toxicity_score', '>', 0.6);
    }

    const results = await query
      .orderBy('smp.posted_at', 'desc')
      .limit(limit)
      .offset(offset);

    return results.map(row => ({
      id: `post-${row.id}`,
      type: 'social_post' as const,
      timestamp: new Date(row.timestamp),
      politician: {
        id: row.politician_id,
        name: row.full_name,
        party: row.party,
        state: row.state,
        avatar: row.image_url
      },
      content: {
        platform: row.platform,
        text: row.content,
        url: row.post_url,
        topics: JSON.parse(row.topics || '[]')
      },
      metadata: {
        engagementCount: (row.likes_count || 0) + (row.shares_count || 0) + (row.comments_count || 0),
        sentiment: row.sentiment_score,
        toxicity: row.toxicity_score,
        importance: row.contains_hate_speech ? 1.0 : (row.toxicity_score || 0)
      }
    }));
  }

  /**
   * Get alert feed items
   */
  private async getAlertFeedItems(
    filter: FeedFilter,
    offset: number,
    limit: number
  ): Promise<FeedItem[]> {
    let query = this.db('politician_alerts as pa')
      .join('politicians as p', 'pa.politician_id', 'p.id')
      .select(
        'pa.id',
        'pa.alert_type',
        'pa.severity',
        'pa.category',
        'pa.title',
        'pa.description',
        'pa.alert_date as timestamp',
        'pa.evidence',
        'pa.is_verified',
        'p.id as politician_id',
        'p.full_name',
        'p.party',
        'p.state',
        'p.image_url'
      )
      .where('pa.is_dismissed', false);

    query = this.applyFilters(query, filter, 'pa', 'p');

    if (filter.types && !filter.types.includes('alert')) {
      return [];
    }

    const results = await query
      .orderBy('pa.alert_date', 'desc')
      .limit(limit)
      .offset(offset);

    return results.map(row => ({
      id: `alert-${row.id}`,
      type: 'alert' as const,
      timestamp: new Date(row.timestamp),
      politician: {
        id: row.politician_id,
        name: row.full_name,
        party: row.party,
        state: row.state,
        avatar: row.image_url
      },
      content: {
        alertType: row.alert_type,
        severity: row.severity,
        category: row.category,
        title: row.title,
        description: row.description,
        evidence: JSON.parse(row.evidence || '{}'),
        isVerified: row.is_verified
      },
      metadata: {
        importance: row.severity === 'critical' ? 1.0 :
                   row.severity === 'high' ? 0.8 :
                   row.severity === 'medium' ? 0.6 : 0.4
      }
    }));
  }

  /**
   * Apply common filters to query
   */
  private applyFilters(
    query: any,
    filter: FeedFilter,
    tableAlias: string,
    politicianAlias: string
  ): any {
    if (filter.politicianIds && filter.politicianIds.length > 0) {
      query = query.whereIn(`${politicianAlias}.id`, filter.politicianIds);
    }

    if (filter.states && filter.states.length > 0) {
      query = query.whereIn(`${politicianAlias}.state`, filter.states);
    }

    if (filter.parties && filter.parties.length > 0) {
      query = query.whereIn(`${politicianAlias}.party`, filter.parties);
    }

    if (filter.startDate) {
      const timestampCol = this.getTimestampColumn(tableAlias);
      query = query.where(timestampCol, '>=', filter.startDate);
    }

    if (filter.endDate) {
      const timestampCol = this.getTimestampColumn(tableAlias);
      query = query.where(timestampCol, '<=', filter.endDate);
    }

    return query;
  }

  /**
   * Get timestamp column name for table
   */
  private getTimestampColumn(tableAlias: string): string {
    const columnMap: Record<string, string> = {
      'v': 'v.vote_date',
      'b': 'b.introduced_date',
      'smp': 'smp.posted_at',
      'pa': 'pa.alert_date'
    };
    return columnMap[tableAlias] || `${tableAlias}.created_at`;
  }

  /**
   * Calculate vote alignment with previous positions
   */
  private async calculateVoteAlignment(
    politicianId: string,
    voteId: string
  ): Promise<number> {
    // Simplified - would use consistency analysis
    return 0.75; // Placeholder
  }

  /**
   * Get total count of feed items
   */
  private async getTotalFeedItems(filter: FeedFilter): Promise<number> {
    // Simplified count - would need to count across all types
    return 1000; // Placeholder
  }

  /**
   * Get personalized feed for a user
   */
  async getPersonalizedFeed(
    userId: string,
    page: number = 1,
    pageSize: number = 50
  ): Promise<{ items: FeedItem[]; total: number; hasMore: boolean }> {
    // Get user preferences (would come from user_preferences table)
    const preferences = await this.getUserPreferences(userId);

    const filter: FeedFilter = {
      politicianIds: preferences.followedPoliticians,
      states: preferences.interestedStates,
      topics: preferences.interestedTopics,
      minImportance: preferences.minImportance || 0.3
    };

    return this.getFeed(filter, page, pageSize);
  }

  /**
   * Get user preferences (placeholder)
   */
  private async getUserPreferences(userId: string): Promise<any> {
    // Would fetch from database
    return {
      followedPoliticians: [],
      interestedStates: [],
      interestedTopics: [],
      minImportance: 0.5
    };
  }

  /**
   * Get trending items (most engagement/importance)
   */
  async getTrendingFeed(
    timeRange: 'day' | 'week' | 'month' = 'day',
    limit: number = 20
  ): Promise<FeedItem[]> {
    const now = new Date();
    const startDate = new Date();
    
    switch (timeRange) {
      case 'day':
        startDate.setDate(startDate.getDate() - 1);
        break;
      case 'week':
        startDate.setDate(startDate.getDate() - 7);
        break;
      case 'month':
        startDate.setMonth(startDate.getMonth() - 1);
        break;
    }

    const filter: FeedFilter = {
      startDate,
      endDate: now,
      minImportance: 0.7
    };

    const feed = await this.getFeed(filter, 1, limit);
    
    // Sort by importance and engagement
    return feed.items.sort((a, b) => {
      const scoreA = (a.metadata.importance || 0) + 
                    ((a.metadata.engagementCount || 0) / 10000);
      const scoreB = (b.metadata.importance || 0) + 
                    ((b.metadata.engagementCount || 0) / 10000);
      return scoreB - scoreA;
    });
  }

  /**
   * Get comparative feed for multiple politicians
   */
  async getComparativeFeed(
    politicianIds: string[],
    page: number = 1,
    pageSize: number = 50
  ): Promise<{ items: FeedItem[]; total: number; hasMore: boolean }> {
    const filter: FeedFilter = {
      politicianIds
    };

    return this.getFeed(filter, page, pageSize);
  }
}

// Export singleton instance
export default new FeedIntegrationService();
