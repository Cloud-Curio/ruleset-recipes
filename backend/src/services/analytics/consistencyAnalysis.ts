/**
 * Consistency Analysis Service
 * 
 * Analyzes consistency between politicians' public statements and voting records.
 * Identifies contradictions and tracks position changes over time.
 */

import { Knex } from 'knex';
import db from '../../database/connection';
import nlpAnalysis from './nlpAnalysis';

interface ConsistencyResult {
  consistencyScore: number;
  statementVoteAlignment: number;
  inconsistencies: Inconsistency[];
  positionChanges: PositionChange[];
  supportingEvidence: Evidence[];
}

interface Inconsistency {
  topic: string;
  statement: string;
  statementDate: Date;
  statementPosition: 'support' | 'oppose' | 'neutral';
  vote: string;
  voteDate: Date;
  votePosition: 'yes' | 'no';
  alignmentScore: number;
  severity: 'low' | 'medium' | 'high';
}

interface PositionChange {
  topic: string;
  previousPosition: 'support' | 'oppose' | 'neutral';
  newPosition: 'support' | 'oppose' | 'neutral';
  changeDate: Date;
  context: string;
}

interface Evidence {
  type: 'vote' | 'statement' | 'bill';
  date: Date;
  description: string;
  position: string;
}

export class ConsistencyAnalysisService {
  private db: Knex;

  constructor() {
    this.db = db;
  }

  /**
   * Analyze consistency for a politician across all policy topics
   */
  async analyzeConsistency(
    politicianId: string,
    startDate: Date,
    endDate: Date
  ): Promise<Map<string, ConsistencyResult>> {
    const results = new Map<string, ConsistencyResult>();

    // Get all policy topics for this politician
    const topics = await this.getPolicyTopics(politicianId, startDate, endDate);

    for (const topic of topics) {
      const result = await this.analyzeTopicConsistency(
        politicianId,
        topic,
        startDate,
        endDate
      );
      results.set(topic, result);
    }

    // Store results in database
    for (const [topic, result] of results.entries()) {
      await this.storeConsistencyAnalysis(
        politicianId,
        topic,
        startDate,
        endDate,
        result
      );
    }

    return results;
  }

  /**
   * Analyze consistency for a specific policy topic
   */
  private async analyzeTopicConsistency(
    politicianId: string,
    topic: string,
    startDate: Date,
    endDate: Date
  ): Promise<ConsistencyResult> {
    // Get statements about this topic
    const statements = await this.getTopicStatements(
      politicianId,
      topic,
      startDate,
      endDate
    );

    // Get votes related to this topic
    const votes = await this.getTopicVotes(
      politicianId,
      topic,
      startDate,
      endDate
    );

    // Compare statements with votes
    const inconsistencies = await this.findInconsistencies(
      statements,
      votes,
      topic
    );

    // Detect position changes
    const positionChanges = this.detectPositionChanges(
      statements,
      votes
    );

    // Gather supporting evidence
    const supportingEvidence = this.gatherEvidence(
      statements,
      votes
    );

    // Calculate consistency scores
    const consistencyScore = this.calculateConsistencyScore(
      statements.length,
      votes.length,
      inconsistencies.length
    );

    const statementVoteAlignment = this.calculateAlignmentScore(
      statements,
      votes,
      inconsistencies
    );

    return {
      consistencyScore,
      statementVoteAlignment,
      inconsistencies,
      positionChanges,
      supportingEvidence
    };
  }

  /**
   * Get policy topics for a politician
   */
  private async getPolicyTopics(
    politicianId: string,
    startDate: Date,
    endDate: Date
  ): Promise<string[]> {
    // Get topics from NLP analysis
    const nlpTopics = await this.db('nlp_analysis')
      .select('topics')
      .where('politician_id', politicianId)
      .whereBetween('analyzed_at', [startDate, endDate]);

    // Get topics from votes
    const voteTopics = await this.db('votes as v')
      .join('bills as b', 'v.bill_id', 'b.id')
      .select('b.policy_areas')
      .where('v.politician_id', politicianId)
      .whereBetween('v.vote_date', [startDate, endDate]);

    // Combine and deduplicate topics
    const allTopics = new Set<string>();
    
    nlpTopics.forEach(row => {
      const topics = JSON.parse(row.topics || '[]');
      topics.forEach((t: string) => allTopics.add(t));
    });

    voteTopics.forEach(row => {
      const topics = JSON.parse(row.policy_areas || '[]');
      topics.forEach((t: string) => allTopics.add(t));
    });

    return Array.from(allTopics);
  }

  /**
   * Get statements about a specific topic
   */
  private async getTopicStatements(
    politicianId: string,
    topic: string,
    startDate: Date,
    endDate: Date
  ): Promise<any[]> {
    const statements = await this.db('nlp_analysis as nlp')
      .join('social_media_posts as smp', function() {
        this.on('nlp.content_id', '=', 'smp.id')
          .andOn(this.db.raw("nlp.content_type = 'social_post'"));
      })
      .select(
        'nlp.id',
        'nlp.policy_positions',
        'nlp.sentiment_score',
        'nlp.analyzed_text',
        'smp.posted_at as date',
        'smp.content'
      )
      .where('nlp.politician_id', politicianId)
      .whereBetween('smp.posted_at', [startDate, endDate])
      .whereRaw("nlp.topics::text ILIKE ?", [`%${topic}%`]);

    return statements.map(s => ({
      ...s,
      policy_positions: JSON.parse(s.policy_positions || '[]')
    }));
  }

  /**
   * Get votes related to a specific topic
   */
  private async getTopicVotes(
    politicianId: string,
    topic: string,
    startDate: Date,
    endDate: Date
  ): Promise<any[]> {
    return await this.db('votes as v')
      .join('bills as b', 'v.bill_id', 'b.id')
      .select(
        'v.id',
        'v.vote_position',
        'v.vote_date as date',
        'b.title',
        'b.policy_areas',
        'b.congress_bill_id'
      )
      .where('v.politician_id', politicianId)
      .whereBetween('v.vote_date', [startDate, endDate])
      .whereRaw("b.policy_areas::text ILIKE ?", [`%${topic}%`]);
  }

  /**
   * Find inconsistencies between statements and votes
   */
  private async findInconsistencies(
    statements: any[],
    votes: any[],
    topic: string
  ): Promise<Inconsistency[]> {
    const inconsistencies: Inconsistency[] = [];

    for (const statement of statements) {
      const stmtPosition = statement.policy_positions.find(
        (p: any) => p.topic === topic
      );

      if (!stmtPosition) continue;

      // Find votes within 90 days of statement
      const relevantVotes = votes.filter(vote => {
        const daysDiff = Math.abs(
          (new Date(vote.date).getTime() - new Date(statement.date).getTime()) / 
          (1000 * 60 * 60 * 24)
        );
        return daysDiff <= 90;
      });

      for (const vote of relevantVotes) {
        // Check alignment
        const aligned = this.checkAlignment(
          stmtPosition.stance,
          vote.vote_position
        );

        if (!aligned) {
          const alignmentScore = this.calculateSingleAlignment(
            stmtPosition.stance,
            vote.vote_position
          );

          inconsistencies.push({
            topic,
            statement: statement.analyzed_text || statement.content,
            statementDate: new Date(statement.date),
            statementPosition: stmtPosition.stance,
            vote: vote.title,
            voteDate: new Date(vote.date),
            votePosition: vote.vote_position,
            alignmentScore,
            severity: this.determineSeverity(alignmentScore)
          });
        }
      }
    }

    return inconsistencies;
  }

  /**
   * Check if statement stance aligns with vote position
   */
  private checkAlignment(
    stance: 'support' | 'oppose' | 'neutral',
    votePosition: 'yes' | 'no'
  ): boolean {
    if (stance === 'neutral') return true;
    if (stance === 'support' && votePosition === 'yes') return true;
    if (stance === 'oppose' && votePosition === 'no') return true;
    return false;
  }

  /**
   * Calculate alignment score between statement and vote
   */
  private calculateSingleAlignment(
    stance: string,
    votePosition: string
  ): number {
    if (stance === 'neutral') return 0.5;
    if ((stance === 'support' && votePosition === 'yes') ||
        (stance === 'oppose' && votePosition === 'no')) {
      return 1.0;
    }
    return 0.0;
  }

  /**
   * Determine severity of inconsistency
   */
  private determineSeverity(alignmentScore: number): 'low' | 'medium' | 'high' {
    if (alignmentScore >= 0.7) return 'low';
    if (alignmentScore >= 0.3) return 'medium';
    return 'high';
  }

  /**
   * Detect position changes over time
   */
  private detectPositionChanges(
    statements: any[],
    votes: any[]
  ): PositionChange[] {
    const changes: PositionChange[] = [];
    
    // Sort by date
    const timeline = [
      ...statements.map(s => ({ type: 'statement', ...s })),
      ...votes.map(v => ({ type: 'vote', ...v }))
    ].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

    // Track position over time
    let lastPosition: 'support' | 'oppose' | 'neutral' | null = null;
    let lastDate: Date | null = null;

    for (const item of timeline) {
      let currentPosition: 'support' | 'oppose' | 'neutral' | null = null;

      if (item.type === 'statement') {
        const pos = item.policy_positions[0];
        if (pos) currentPosition = pos.stance;
      } else {
        currentPosition = item.vote_position === 'yes' ? 'support' : 'oppose';
      }

      if (currentPosition && lastPosition && currentPosition !== lastPosition) {
        changes.push({
          topic: item.topic || 'unknown',
          previousPosition: lastPosition,
          newPosition: currentPosition,
          changeDate: new Date(item.date),
          context: item.title || item.analyzed_text || 'position change detected'
        });
      }

      if (currentPosition) {
        lastPosition = currentPosition;
        lastDate = new Date(item.date);
      }
    }

    return changes;
  }

  /**
   * Gather supporting evidence
   */
  private gatherEvidence(statements: any[], votes: any[]): Evidence[] {
    const evidence: Evidence[] = [];

    statements.forEach(stmt => {
      evidence.push({
        type: 'statement',
        date: new Date(stmt.date),
        description: stmt.analyzed_text || stmt.content,
        position: stmt.policy_positions[0]?.stance || 'neutral'
      });
    });

    votes.forEach(vote => {
      evidence.push({
        type: 'vote',
        date: new Date(vote.date),
        description: vote.title,
        position: vote.vote_position
      });
    });

    return evidence.sort((a, b) => a.date.getTime() - b.date.getTime());
  }

  /**
   * Calculate overall consistency score
   */
  private calculateConsistencyScore(
    statementCount: number,
    voteCount: number,
    inconsistencyCount: number
  ): number {
    const totalActions = statementCount + voteCount;
    if (totalActions === 0) return 100;

    const consistencyRate = 1 - (inconsistencyCount / totalActions);
    return Math.round(consistencyRate * 100);
  }

  /**
   * Calculate statement-vote alignment score
   */
  private calculateAlignmentScore(
    statements: any[],
    votes: any[],
    inconsistencies: Inconsistency[]
  ): number {
    const totalPairs = statements.length * votes.length;
    if (totalPairs === 0) return 100;

    const alignedPairs = totalPairs - inconsistencies.length;
    return Math.round((alignedPairs / totalPairs) * 100);
  }

  /**
   * Store consistency analysis results
   */
  private async storeConsistencyAnalysis(
    politicianId: string,
    topic: string,
    startDate: Date,
    endDate: Date,
    result: ConsistencyResult
  ): Promise<void> {
    await this.db('consistency_analysis')
      .insert({
        politician_id: politicianId,
        policy_topic: topic,
        analysis_period_start: startDate,
        analysis_period_end: endDate,
        consistency_score: result.consistencyScore,
        statement_vote_alignment: result.statementVoteAlignment,
        total_votes_analyzed: result.supportingEvidence.filter(e => e.type === 'vote').length,
        total_statements_analyzed: result.supportingEvidence.filter(e => e.type === 'statement').length,
        inconsistencies: JSON.stringify(result.inconsistencies),
        position_changes: JSON.stringify(result.positionChanges),
        supporting_evidence: JSON.stringify(result.supportingEvidence),
        summary: this.generateSummary(result)
      })
      .onConflict(['politician_id', 'policy_topic', 'analysis_period_start'])
      .merge();
  }

  /**
   * Generate text summary of consistency analysis
   */
  private generateSummary(result: ConsistencyResult): string {
    const parts: string[] = [];

    parts.push(`Consistency Score: ${result.consistencyScore}%`);
    parts.push(`Statement-Vote Alignment: ${result.statementVoteAlignment}%`);
    
    if (result.inconsistencies.length > 0) {
      parts.push(`Found ${result.inconsistencies.length} inconsistencies`);
      const highSeverity = result.inconsistencies.filter(i => i.severity === 'high').length;
      if (highSeverity > 0) {
        parts.push(`${highSeverity} high severity inconsistencies detected`);
      }
    }

    if (result.positionChanges.length > 0) {
      parts.push(`${result.positionChanges.length} position changes detected`);
    }

    return parts.join('. ') + '.';
  }

  /**
   * Get consistency report for politician
   */
  async getConsistencyReport(politicianId: string): Promise<any> {
    return await this.db('consistency_analysis')
      .where('politician_id', politicianId)
      .orderBy('analysis_period_end', 'desc');
  }
}

// Export singleton instance
export default new ConsistencyAnalysisService();
