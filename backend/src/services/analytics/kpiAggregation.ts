/**
 * KPI Aggregation Service
 * 
 * Calculates and aggregates Key Performance Indicators for politicians.
 * Generates comprehensive scorecards across multiple dimensions.
 */

import { Knex } from 'knex';
import db from '../../database/connection';

interface KPIMetrics {
  integrityScore: number;
  honestyScore: number;
  consistencyScore: number;
  transparencyScore: number;
  effectivenessScore: number;
  bipartisanScore: number;
  attendanceRate: number;
  votesCast: number;
  votesMissed: number;
  billsSponsored: number;
  billsCosponsored: number;
  amendmentsProposed: number;
  speechesGiven: number;
  socialPostsCount: number;
  totalEngagement: number;
  averageToxicity: number;
  controversialPosts: number;
  leadershipScore: number;
  influenceScore: number;
  billsPassed: number;
  successfulAmendments: number;
  constituentAlignment: number;
  districtApprovalRating: number;
}

export class KPIAggregationService {
  private db: Knex;

  constructor() {
    this.db = db;
  }

  /**
   * Calculate all KPIs for a politician for a given time period
   */
  async calculateKPIs(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date,
    periodType: 'weekly' | 'monthly' | 'quarterly' | 'annual' = 'monthly'
  ): Promise<KPIMetrics> {
    const [
      coreMetrics,
      engagementMetrics,
      socialMetrics,
      influenceMetrics,
      constituentMetrics
    ] = await Promise.all([
      this.calculateCoreMetrics(politicianId, periodStart, periodEnd),
      this.calculateEngagementMetrics(politicianId, periodStart, periodEnd),
      this.calculateSocialMetrics(politicianId, periodStart, periodEnd),
      this.calculateInfluenceMetrics(politicianId, periodStart, periodEnd),
      this.calculateConstituentMetrics(politicianId, periodStart, periodEnd)
    ]);

    const kpis: KPIMetrics = {
      ...coreMetrics,
      ...engagementMetrics,
      ...socialMetrics,
      ...influenceMetrics,
      ...constituentMetrics
    };

    // Store KPIs in database
    await this.storeKPIs(politicianId, periodStart, periodEnd, periodType, kpis);

    return kpis;
  }

  /**
   * Calculate core integrity metrics
   */
  private async calculateCoreMetrics(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date
  ): Promise<Partial<KPIMetrics>> {
    // Get consistency score
    const consistency = await this.db('consistency_analysis')
      .where('politician_id', politicianId)
      .whereBetween('analysis_period_end', [periodStart, periodEnd])
      .avg('consistency_score as avg_consistency')
      .first();

    // Get promise keep rate (honesty)
    const promises = await this.db('campaign_promises')
      .where('politician_id', politicianId)
      .select(
        this.db.raw(`
          COUNT(*) as total,
          COUNT(CASE WHEN status = 'kept' THEN 1 END) as kept
        `)
      )
      .first();

    const honestyScore = promises && promises.total > 0
      ? (promises.kept / promises.total) * 100
      : 50;

    // Get transparency metrics (participation + communication)
    const transparency = await this.calculateTransparencyScore(
      politicianId,
      periodStart,
      periodEnd
    );

    // Get bipartisan score
    const bipartisan = await this.calculateBipartisanScore(
      politicianId,
      periodStart,
      periodEnd
    );

    // Overall integrity is weighted average
    const integrityScore = Math.round(
      (consistency?.avg_consistency || 50) * 0.3 +
      honestyScore * 0.3 +
      transparency * 0.25 +
      bipartisan * 0.15
    );

    return {
      integrityScore,
      honestyScore: Math.round(honestyScore),
      consistencyScore: Math.round(consistency?.avg_consistency || 50),
      transparencyScore: Math.round(transparency),
      bipartisanScore: Math.round(bipartisan)
    };
  }

  /**
   * Calculate engagement metrics
   */
  private async calculateEngagementMetrics(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date
  ): Promise<Partial<KPIMetrics>> {
    // Voting participation
    const voting = await this.db('votes')
      .where('politician_id', politicianId)
      .whereBetween('vote_date', [periodStart, periodEnd])
      .select(
        this.db.raw(`
          COUNT(*) as total_votes,
          COUNT(CASE WHEN vote_position IN ('yes', 'no') THEN 1 END) as votes_cast,
          COUNT(CASE WHEN vote_position = 'not_voting' THEN 1 END) as votes_missed
        `)
      )
      .first();

    const attendanceRate = voting && voting.total_votes > 0
      ? (voting.votes_cast / voting.total_votes) * 100
      : 0;

    // Bill sponsorship
    const bills = await this.db('bills')
      .where('sponsor_id', politicianId)
      .whereBetween('introduced_date', [periodStart, periodEnd])
      .count('* as count')
      .first();

    // Co-sponsorship (simplified - would need junction table)
    const cosponsored = 0; // Placeholder

    return {
      attendanceRate: Math.round(attendanceRate * 100) / 100,
      votesCast: voting?.votes_cast || 0,
      votesMissed: voting?.votes_missed || 0,
      billsSponsored: bills?.count || 0,
      billsCosponsored: cosponsored,
      amendmentsProposed: 0, // Would need amendments table
      speechesGiven: 0 // Would need speeches table
    };
  }

  /**
   * Calculate social media metrics
   */
  private async calculateSocialMetrics(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date
  ): Promise<Partial<KPIMetrics>> {
    const social = await this.db('social_media_posts')
      .where('politician_id', politicianId)
      .whereBetween('posted_at', [periodStart, periodEnd])
      .select(
        this.db.raw(`
          COUNT(*) as post_count,
          SUM(likes_count + shares_count + comments_count) as total_engagement
        `)
      )
      .first();

    // Get NLP toxicity metrics
    const nlp = await this.db('nlp_analysis')
      .where('politician_id', politicianId)
      .where('content_type', 'social_post')
      .whereBetween('analyzed_at', [periodStart, periodEnd])
      .select(
        this.db.raw(`
          AVG(toxicity_score) as avg_toxicity,
          COUNT(CASE WHEN toxicity_score > 0.7 THEN 1 END) as controversial
        `)
      )
      .first();

    return {
      socialPostsCount: social?.post_count || 0,
      totalEngagement: social?.total_engagement || 0,
      averageToxicity: nlp?.avg_toxicity || 0,
      controversialPosts: nlp?.controversial || 0
    };
  }

  /**
   * Calculate influence metrics
   */
  private async calculateInfluenceMetrics(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date
  ): Promise<Partial<KPIMetrics>> {
    // Bills passed
    const passed = await this.db('bills')
      .where('sponsor_id', politicianId)
      .whereIn('status', ['signed', 'enacted'])
      .whereBetween('introduced_date', [periodStart, periodEnd])
      .count('* as count')
      .first();

    // Leadership score based on committee positions, bill success, etc.
    const leadershipScore = await this.calculateLeadershipScore(politicianId);

    // Influence score based on bill passage rate, bipartisan support, etc.
    const influenceScore = await this.calculateInfluenceScore(
      politicianId,
      periodStart,
      periodEnd
    );

    // Effectiveness combines leadership and influence
    const effectivenessScore = Math.round((leadershipScore + influenceScore) / 2);

    return {
      leadershipScore: Math.round(leadershipScore),
      influenceScore: Math.round(influenceScore),
      effectivenessScore,
      billsPassed: passed?.count || 0,
      successfulAmendments: 0 // Placeholder
    };
  }

  /**
   * Calculate constituent alignment metrics
   */
  private async calculateConstituentMetrics(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date
  ): Promise<Partial<KPIMetrics>> {
    // This would require polling data and district demographics
    // For now, return placeholders
    return {
      constituentAlignment: 50, // Placeholder
      districtApprovalRating: 50 // Placeholder
    };
  }

  /**
   * Calculate transparency score
   */
  private async calculateTransparencyScore(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date
  ): Promise<number> {
    // Based on attendance + communication frequency
    const voting = await this.db('votes')
      .where('politician_id', politicianId)
      .whereBetween('vote_date', [periodStart, periodEnd])
      .select(
        this.db.raw(`
          COUNT(CASE WHEN vote_position IN ('yes', 'no') THEN 1 END)::numeric / 
          NULLIF(COUNT(*), 0) as participation
        `)
      )
      .first();

    const social = await this.db('social_media_posts')
      .where('politician_id', politicianId)
      .whereBetween('posted_at', [periodStart, periodEnd])
      .count('* as count')
      .first();

    const participation = voting?.participation || 0;
    const communication = Math.min(1, (social?.count || 0) / 30); // Normalize to monthly posts

    return (participation * 70 + communication * 30) * 100;
  }

  /**
   * Calculate bipartisan score
   */
  private async calculateBipartisanScore(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date
  ): Promise<number> {
    const votes = await this.db('votes')
      .where('politician_id', politicianId)
      .whereBetween('vote_date', [periodStart, periodEnd])
      .select(
        this.db.raw(`
          COUNT(*) as total,
          COUNT(CASE WHEN party_line_vote = false THEN 1 END) as cross_party
        `)
      )
      .first();

    if (!votes || votes.total === 0) return 50;

    return (votes.cross_party / votes.total) * 100;
  }

  /**
   * Calculate leadership score
   */
  private async calculateLeadershipScore(politicianId: string): Promise<number> {
    const politician = await this.db('politicians')
      .where('id', politicianId)
      .first();

    let score = 50; // Base score

    // Add points for leadership roles
    const leadershipRoles = JSON.parse(politician?.leadership_roles || '[]');
    score += leadershipRoles.length * 10;

    // Add points for committee leadership
    const committees = JSON.parse(politician?.committees || '[]');
    const chairPositions = committees.filter((c: any) => 
      c.position === 'chair' || c.position === 'ranking'
    );
    score += chairPositions.length * 15;

    return Math.min(100, score);
  }

  /**
   * Calculate influence score
   */
  private async calculateInfluenceScore(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date
  ): Promise<number> {
    // Bill success rate
    const bills = await this.db('bills')
      .where('sponsor_id', politicianId)
      .whereBetween('introduced_date', [periodStart, periodEnd])
      .select(
        this.db.raw(`
          COUNT(*) as total,
          COUNT(CASE WHEN status IN ('signed', 'enacted') THEN 1 END) as passed,
          AVG(bipartisan_support) as avg_bipartisan
        `)
      )
      .first();

    if (!bills || bills.total === 0) return 50;

    const successRate = (bills.passed / bills.total) * 100;
    const bipartisanSupport = bills.avg_bipartisan || 0;

    return (successRate * 0.6 + bipartisanSupport * 0.4);
  }

  /**
   * Store KPIs in database
   */
  private async storeKPIs(
    politicianId: string,
    periodStart: Date,
    periodEnd: Date,
    periodType: string,
    kpis: KPIMetrics
  ): Promise<void> {
    await this.db('politician_kpis')
      .insert({
        politician_id: politicianId,
        period_start: periodStart,
        period_end: periodEnd,
        period_type: periodType,
        integrity_score: kpis.integrityScore,
        honesty_score: kpis.honestyScore,
        consistency_score: kpis.consistencyScore,
        transparency_score: kpis.transparencyScore,
        effectiveness_score: kpis.effectivenessScore,
        bipartisan_score: kpis.bipartisanScore,
        votes_cast: kpis.votesCast,
        votes_missed: kpis.votesMissed,
        attendance_rate: kpis.attendanceRate,
        bills_sponsored: kpis.billsSponsored,
        bills_cosponsored: kpis.billsCosponsored,
        amendments_proposed: kpis.amendmentsProposed,
        speeches_given: kpis.speechesGiven,
        social_posts_count: kpis.socialPostsCount,
        total_engagement: kpis.totalEngagement,
        average_toxicity: kpis.averageToxicity,
        controversial_posts: kpis.controversialPosts,
        leadership_score: kpis.leadershipScore,
        influence_score: kpis.influenceScore,
        bills_passed: kpis.billsPassed,
        successful_amendments: kpis.successfulAmendments,
        constituent_alignment: kpis.constituentAlignment,
        district_approval_rating: kpis.districtApprovalRating
      })
      .onConflict(['politician_id', 'period_start', 'period_type'])
      .merge();
  }

  /**
   * Calculate KPIs for all active politicians
   */
  async calculateAllPoliticiansKPIs(
    periodStart: Date,
    periodEnd: Date,
    periodType: 'weekly' | 'monthly' | 'quarterly' | 'annual' = 'monthly'
  ): Promise<void> {
    const politicians = await this.db('politicians')
      .where('in_office', true)
      .select('id');

    console.log(`Calculating KPIs for ${politicians.length} politicians...`);

    for (const politician of politicians) {
      try {
        await this.calculateKPIs(
          politician.id,
          periodStart,
          periodEnd,
          periodType
        );
        console.log(`Calculated KPIs for politician ${politician.id}`);
      } catch (error) {
        console.error(`Error calculating KPIs for ${politician.id}:`, error);
      }
    }

    console.log('KPI calculation completed');
  }

  /**
   * Get KPI trends for a politician
   */
  async getKPITrends(
    politicianId: string,
    months: number = 12
  ): Promise<any[]> {
    const startDate = new Date();
    startDate.setMonth(startDate.getMonth() - months);

    return await this.db('politician_kpis')
      .where('politician_id', politicianId)
      .where('period_start', '>=', startDate)
      .where('period_type', 'monthly')
      .orderBy('period_start', 'asc');
  }

  /**
   * Get comparative rankings
   */
  async getRankings(
    metric: keyof KPIMetrics,
    chamber?: string,
    party?: string,
    limit: number = 50
  ): Promise<any[]> {
    let query = this.db('politician_kpis as pk')
      .join('politicians as p', 'pk.politician_id', 'p.id')
      .select(
        'p.id',
        'p.full_name',
        'p.party',
        'p.state',
        'p.chamber',
        `pk.${metric}`
      )
      .where('p.in_office', true)
      .orderBy(`pk.${metric}`, 'desc')
      .limit(limit);

    if (chamber) {
      query = query.where('p.chamber', chamber);
    }

    if (party) {
      query = query.where('p.party', party);
    }

    // Get latest KPIs
    query = query.whereRaw(`
      pk.period_start = (
        SELECT MAX(period_start) 
        FROM politician_kpis 
        WHERE politician_id = pk.politician_id
      )
    `);

    return await query;
  }
}

// Export singleton instance
export default new KPIAggregationService();
