# Implementation Guide

## Overview

This guide provides technical implementation details for translating research methodologies into working code. It covers architecture, code structure, best practices, and specific implementation patterns for political data analysis.

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Frontend Layer                     │
│  - React/Next.js UI                                 │
│  - Data Visualization (Charts, Networks)           │
│  - User Interactions                                │
└─────────────────┬───────────────────────────────────┘
                  │
                  │ REST API / GraphQL
                  │
┌─────────────────▼───────────────────────────────────┐
│                   API Layer                          │
│  - Express.js Routes                                │
│  - Request Validation                               │
│  - Authentication/Authorization                     │
└─────────────────┬───────────────────────────────────┘
                  │
                  │
┌─────────────────▼───────────────────────────────────┐
│                 Business Logic Layer                 │
│  - Analytics Services                               │
│  - Political Figure Analysis                        │
│  - Legislation Analysis                             │
│  - Data Processing                                  │
└─────────────────┬───────────────────────────────────┘
                  │
                  │
┌─────────────────▼───────────────────────────────────┐
│                   Data Layer                         │
│  - PostgreSQL (Primary Data)                        │
│  - Redis (Cache, Sessions)                          │
│  - File Storage (Documents, PDFs)                   │
└─────────────────┬───────────────────────────────────┘
                  │
                  │
┌─────────────────▼───────────────────────────────────┐
│              External Data Sources                   │
│  - Congress.gov API                                 │
│  - GovInfo.gov API                                  │
│  - OpenStates API                                   │
└─────────────────────────────────────────────────────┘
```

## Project Structure

```
backend/
├── src/
│   ├── index.ts                    # Application entry point
│   ├── config/                     # Configuration
│   │   ├── database.ts
│   │   ├── redis.ts
│   │   └── api-keys.ts
│   ├── routes/                     # API routes
│   │   ├── politicians.ts
│   │   ├── bills.ts
│   │   ├── votes.ts
│   │   └── analytics.ts
│   ├── controllers/                # Request handlers
│   │   ├── politician.controller.ts
│   │   ├── bill.controller.ts
│   │   └── analytics.controller.ts
│   ├── services/                   # Business logic
│   │   ├── data-collection/
│   │   │   ├── congress-api.service.ts
│   │   │   ├── govinfo-api.service.ts
│   │   │   └── openstates-api.service.ts
│   │   ├── analysis/
│   │   │   ├── voting-analysis.service.ts
│   │   │   ├── similarity-analysis.service.ts
│   │   │   ├── bill-classification.service.ts
│   │   │   └── nlp-analysis.service.ts
│   │   └── cache/
│   │       └── redis.service.ts
│   ├── models/                     # Database models
│   │   ├── Politician.ts
│   │   ├── Bill.ts
│   │   ├── Vote.ts
│   │   └── VotingRecord.ts
│   ├── jobs/                       # Background jobs
│   │   ├── data-collection.job.ts
│   │   ├── analysis.job.ts
│   │   └── cache-warmup.job.ts
│   ├── analytics/                  # Analytics modules
│   │   ├── metrics/
│   │   │   ├── bipartisanship.ts
│   │   │   ├── productivity.ts
│   │   │   └── influence.ts
│   │   ├── ml/
│   │   │   ├── bill-classifier.ts
│   │   │   ├── similarity.ts
│   │   │   └── clustering.ts
│   │   └── nlp/
│   │       ├── summarization.ts
│   │       ├── classification.ts
│   │       └── entity-extraction.ts
│   ├── database/                   # Database management
│   │   ├── connection.ts
│   │   ├── migrations/
│   │   └── seeds/
│   ├── middleware/                 # Express middleware
│   │   ├── auth.ts
│   │   ├── validation.ts
│   │   └── error-handler.ts
│   └── utils/                      # Utilities
│       ├── logger.ts
│       ├── validators.ts
│       └── helpers.ts
└── tests/                          # Tests
    ├── unit/
    ├── integration/
    └── e2e/
```

## Core Implementation Patterns

### 1. Data Collection Service

```typescript
// src/services/data-collection/congress-api.service.ts

import axios, { AxiosInstance } from 'axios';
import { BillData, MemberData, VoteData } from '../types';

export class CongressAPIService {
  private client: AxiosInstance;
  private lastRequestTime: number = 0;
  private minInterval: number = 1000; // 1 second between requests

  constructor(apiKey: string) {
    this.client = axios.create({
      baseURL: 'https://api.congress.gov/v3',
      headers: {
        'X-Api-Key': apiKey
      }
    });
  }

  private async rateLimit(): Promise<void> {
    const now = Date.now();
    const elapsed = now - this.lastRequestTime;
    
    if (elapsed < this.minInterval) {
      await new Promise(resolve => 
        setTimeout(resolve, this.minInterval - elapsed)
      );
    }
    
    this.lastRequestTime = Date.now();
  }

  async getBill(
    congress: number, 
    billType: string, 
    billNumber: number
  ): Promise<BillData> {
    await this.rateLimit();
    
    try {
      const response = await this.client.get(
        `/bill/${congress}/${billType}/${billNumber}`
      );
      
      return this.transformBillData(response.data);
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 429) {
        // Rate limited - wait and retry
        const retryAfter = parseInt(
          error.response.headers['retry-after'] || '60'
        );
        await new Promise(resolve => setTimeout(resolve, retryAfter * 1000));
        return this.getBill(congress, billType, billNumber);
      }
      throw error;
    }
  }

  async getBillsForCongress(
    congress: number,
    options: {
      limit?: number;
      offset?: number;
    } = {}
  ): Promise<BillData[]> {
    const { limit = 250, offset = 0 } = options;
    
    await this.rateLimit();
    
    const response = await this.client.get(`/bill/${congress}`, {
      params: { limit, offset }
    });
    
    return response.data.bills.map(this.transformBillData);
  }

  async getMember(bioguideId: string): Promise<MemberData> {
    await this.rateLimit();
    
    const response = await this.client.get(`/member/${bioguideId}`);
    return this.transformMemberData(response.data);
  }

  private transformBillData(rawData: any): BillData {
    return {
      congress: rawData.congress,
      billType: rawData.type,
      billNumber: rawData.number,
      title: rawData.title,
      introducedDate: new Date(rawData.introducedDate),
      sponsor: rawData.sponsors?.[0],
      cosponsors: rawData.cosponsors || [],
      subjects: rawData.subjects || [],
      status: rawData.latestAction?.text,
      // Add more fields as needed
    };
  }

  private transformMemberData(rawData: any): MemberData {
    // Transform raw API data to internal format
    return {
      bioguideId: rawData.bioguideId,
      name: rawData.name,
      party: rawData.partyName,
      state: rawData.state,
      // Add more fields
    };
  }
}
```

### 2. Voting Analysis Service

```typescript
// src/services/analysis/voting-analysis.service.ts

import { VotingRecord, PoliticianMetrics } from '../types';
import { calculateMean, calculateStdDev } from '../utils/stats';

export class VotingAnalysisService {
  
  /**
   * Calculate bipartisan voting score for a politician
   */
  async calculateBipartisanScore(
    politicianId: string,
    timeRange?: { start: Date; end: Date }
  ): Promise<number> {
    // 1. Fetch voting records
    const votes = await this.getVotingRecords(politicianId, timeRange);
    
    // 2. Identify party-line votes
    const partyLineVotes = await this.identifyPartyLineVotes(votes);
    
    // 3. Count cross-party votes
    const crossPartyVotes = partyLineVotes.filter(vote => 
      this.isVotingWithOpposition(vote, politicianId)
    );
    
    // 4. Calculate percentage
    if (partyLineVotes.length === 0) return 0;
    
    return (crossPartyVotes.length / partyLineVotes.length) * 100;
  }

  /**
   * Calculate voting participation rate
   */
  async calculateParticipationRate(
    politicianId: string,
    timeRange?: { start: Date; end: Date }
  ): Promise<{
    rate: number;
    votesCast: number;
    totalVotes: number;
  }> {
    const politician = await this.getPolitician(politicianId);
    
    // Get all votes during time period
    const allVotes = await this.getAllVotes(timeRange);
    
    // Get politician's votes
    const politicianVotes = await this.getVotingRecords(
      politicianId, 
      timeRange
    );
    
    const votesCast = politicianVotes.filter(
      v => v.position !== 'Not Voting'
    ).length;
    
    return {
      rate: (votesCast / allVotes.length) * 100,
      votesCast,
      totalVotes: allVotes.length
    };
  }

  /**
   * Calculate cosine similarity between two politicians' voting patterns
   */
  calculateVotingSimilarity(
    votes1: VotingRecord[],
    votes2: VotingRecord[]
  ): number {
    // Create voting vectors (1 for Yea, -1 for Nay, 0 for absent/present)
    const vector1 = this.createVotingVector(votes1);
    const vector2 = this.createVotingVector(votes2);
    
    // Calculate cosine similarity
    return this.cosineSimilarity(vector1, vector2);
  }

  /**
   * Analyze voting patterns by policy area
   */
  async analyzeVotingByPolicyArea(
    politicianId: string
  ): Promise<Map<string, {
    total: number;
    yea: number;
    nay: number;
    yeaPercentage: number;
  }>> {
    const votes = await this.getVotingRecords(politicianId);
    const votesByArea = new Map();
    
    for (const vote of votes) {
      const bill = await this.getBill(vote.billId);
      const policyArea = bill.policyArea || 'Unclassified';
      
      if (!votesByArea.has(policyArea)) {
        votesByArea.set(policyArea, {
          total: 0,
          yea: 0,
          nay: 0,
          yeaPercentage: 0
        });
      }
      
      const stats = votesByArea.get(policyArea);
      stats.total++;
      
      if (vote.position === 'Yea') stats.yea++;
      if (vote.position === 'Nay') stats.nay++;
      
      stats.yeaPercentage = (stats.yea / stats.total) * 100;
    }
    
    return votesByArea;
  }

  private createVotingVector(votes: VotingRecord[]): number[] {
    return votes.map(vote => {
      switch (vote.position) {
        case 'Yea': return 1;
        case 'Nay': return -1;
        default: return 0;
      }
    });
  }

  private cosineSimilarity(vec1: number[], vec2: number[]): number {
    // Ensure same length
    const length = Math.min(vec1.length, vec2.length);
    vec1 = vec1.slice(0, length);
    vec2 = vec2.slice(0, length);
    
    // Calculate dot product
    const dotProduct = vec1.reduce((sum, val, i) => sum + val * vec2[i], 0);
    
    // Calculate magnitudes
    const mag1 = Math.sqrt(vec1.reduce((sum, val) => sum + val * val, 0));
    const mag2 = Math.sqrt(vec2.reduce((sum, val) => sum + val * val, 0));
    
    // Avoid division by zero
    if (mag1 === 0 || mag2 === 0) return 0;
    
    return dotProduct / (mag1 * mag2);
  }

  private async identifyPartyLineVotes(
    votes: VotingRecord[]
  ): Promise<VotingRecord[]> {
    const partyLineVotes = [];
    
    for (const vote of votes) {
      // Get all votes for this roll call
      const allVotes = await this.getVotesForRollCall(vote.rollCallId);
      
      // Check if >80% of party voted same way
      const isPartyLine = this.isPartyLineVote(allVotes, vote.party);
      
      if (isPartyLine) {
        partyLineVotes.push(vote);
      }
    }
    
    return partyLineVotes;
  }

  private isPartyLineVote(
    allVotes: VotingRecord[], 
    party: string
  ): boolean {
    const partyVotes = allVotes.filter(v => v.party === party);
    
    if (partyVotes.length === 0) return false;
    
    // Count most common position
    const positions = partyVotes.map(v => v.position);
    const positionCounts = new Map();
    
    for (const pos of positions) {
      positionCounts.set(pos, (positionCounts.get(pos) || 0) + 1);
    }
    
    const maxCount = Math.max(...positionCounts.values());
    const threshold = partyVotes.length * 0.8;
    
    return maxCount >= threshold;
  }

  private isVotingWithOpposition(
    vote: VotingRecord,
    politicianId: string
  ): boolean {
    // Determine if politician voted with opposition party majority
    // Implementation depends on data structure
    return false; // Placeholder
  }

  // Placeholder methods - implement based on data layer
  private async getVotingRecords(
    politicianId: string,
    timeRange?: { start: Date; end: Date }
  ): Promise<VotingRecord[]> {
    // Fetch from database
    return [];
  }

  private async getPolitician(politicianId: string) {
    // Fetch from database
    return null;
  }

  private async getAllVotes(timeRange?: { start: Date; end: Date }) {
    // Fetch from database
    return [];
  }

  private async getBill(billId: string) {
    // Fetch from database
    return null;
  }

  private async getVotesForRollCall(rollCallId: string) {
    // Fetch from database
    return [];
  }
}
```

### 3. Bill Classification Service

```typescript
// src/services/analysis/bill-classification.service.ts

import natural from 'natural';
import { TfIdf } from 'natural';

export class BillClassificationService {
  private tfidf: TfIdf;
  private classifier: any; // Use appropriate type for your ML library

  constructor() {
    this.tfidf = new TfIdf();
  }

  /**
   * Classify bill by policy area
   */
  async classifyBill(billText: string): Promise<{
    category: string;
    confidence: number;
  }> {
    // 1. Preprocess text
    const processedText = this.preprocessText(billText);
    
    // 2. Extract features
    const features = this.extractFeatures(processedText);
    
    // 3. Classify
    const result = await this.classify(features);
    
    return result;
  }

  /**
   * Extract key topics from bill text
   */
  async extractTopics(
    billText: string,
    numTopics: number = 5
  ): Promise<string[]> {
    // Use TF-IDF to find most important terms
    this.tfidf.addDocument(billText);
    
    const terms: { term: string; score: number }[] = [];
    
    this.tfidf.listTerms(0).forEach((item: any) => {
      terms.push({
        term: item.term,
        score: item.tfidf
      });
    });
    
    // Return top N terms
    return terms
      .sort((a, b) => b.score - a.score)
      .slice(0, numTopics)
      .map(t => t.term);
  }

  /**
   * Calculate bill complexity score
   */
  calculateComplexity(billText: string): {
    wordCount: number;
    readabilityScore: number;
    gradeLevel: number;
    complexityScore: number;
  } {
    const tokenizer = new natural.WordTokenizer();
    const sentenceTokenizer = new natural.SentenceTokenizer();
    
    const words = tokenizer.tokenize(billText);
    const sentences = sentenceTokenizer.tokenize(billText);
    
    const wordCount = words?.length || 0;
    const sentenceCount = sentences?.length || 0;
    
    // Calculate Flesch Reading Ease
    const syllables = this.countSyllables(words || []);
    const readabilityScore = this.fleschReadingEase(
      wordCount,
      sentenceCount,
      syllables
    );
    
    // Calculate Flesch-Kincaid Grade Level
    const gradeLevel = this.fleschKincaidGrade(
      wordCount,
      sentenceCount,
      syllables
    );
    
    // Custom complexity score
    const complexityScore = this.calculateCustomComplexity(billText);
    
    return {
      wordCount,
      readabilityScore,
      gradeLevel,
      complexityScore
    };
  }

  private preprocessText(text: string): string {
    // Convert to lowercase
    text = text.toLowerCase();
    
    // Remove special characters
    text = text.replace(/[^a-z\s]/g, ' ');
    
    // Remove extra whitespace
    text = text.replace(/\s+/g, ' ').trim();
    
    return text;
  }

  private extractFeatures(text: string): number[] {
    // Extract TF-IDF features
    // This is a simplified version - use actual ML library
    return [];
  }

  private async classify(features: number[]): Promise<{
    category: string;
    confidence: number;
  }> {
    // Use trained classifier
    // This is a placeholder
    return {
      category: 'Healthcare',
      confidence: 0.85
    };
  }

  private countSyllables(words: string[]): number {
    let syllableCount = 0;
    
    for (const word of words) {
      // Simple syllable counting (can be improved)
      const vowels = word.match(/[aeiouy]+/g);
      syllableCount += vowels ? vowels.length : 1;
    }
    
    return syllableCount;
  }

  private fleschReadingEase(
    words: number,
    sentences: number,
    syllables: number
  ): number {
    if (sentences === 0 || words === 0) return 0;
    
    return 206.835 - 1.015 * (words / sentences) 
           - 84.6 * (syllables / words);
  }

  private fleschKincaidGrade(
    words: number,
    sentences: number,
    syllables: number
  ): number {
    if (sentences === 0 || words === 0) return 0;
    
    return 0.39 * (words / sentences) 
           + 11.8 * (syllables / words) - 15.59;
  }

  private calculateCustomComplexity(billText: string): number {
    // Count various complexity indicators
    const crossReferences = (billText.match(/\d+ U\.S\.C\./g) || []).length;
    const amendments = (billText.match(/is amended/gi) || []).length;
    const exceptions = (billText.match(/notwithstanding/gi) || []).length;
    
    // Weighted complexity score
    const score = (crossReferences * 2) + (amendments * 1.5) + (exceptions * 1);
    
    // Normalize to 0-100 scale
    return Math.min(100, score / 10);
  }
}
```

### 4. Caching Strategy

```typescript
// src/services/cache/redis.service.ts

import Redis from 'ioredis';

export class CacheService {
  private redis: Redis;

  constructor(redisUrl: string) {
    this.redis = new Redis(redisUrl);
  }

  /**
   * Cache politician metrics
   */
  async cachePoliticianMetrics(
    politicianId: string,
    metrics: any,
    ttl: number = 3600 // 1 hour default
  ): Promise<void> {
    const key = `politician:${politicianId}:metrics`;
    await this.redis.setex(key, ttl, JSON.stringify(metrics));
  }

  /**
   * Get cached politician metrics
   */
  async getCachedPoliticianMetrics(
    politicianId: string
  ): Promise<any | null> {
    const key = `politician:${politicianId}:metrics`;
    const cached = await this.redis.get(key);
    
    return cached ? JSON.parse(cached) : null;
  }

  /**
   * Cache similarity calculation
   */
  async cacheSimilarity(
    politicianId1: string,
    politicianId2: string,
    similarity: number,
    ttl: number = 86400 // 24 hours
  ): Promise<void> {
    // Use sorted IDs as key to avoid duplicates
    const [id1, id2] = [politicianId1, politicianId2].sort();
    const key = `similarity:${id1}:${id2}`;
    
    await this.redis.setex(key, ttl, similarity.toString());
  }

  /**
   * Get cached similarity
   */
  async getCachedSimilarity(
    politicianId1: string,
    politicianId2: string
  ): Promise<number | null> {
    const [id1, id2] = [politicianId1, politicianId2].sort();
    const key = `similarity:${id1}:${id2}`;
    
    const cached = await this.redis.get(key);
    return cached ? parseFloat(cached) : null;
  }

  /**
   * Invalidate cache for a politician
   */
  async invalidatePoliticianCache(politicianId: string): Promise<void> {
    const pattern = `politician:${politicianId}:*`;
    const keys = await this.redis.keys(pattern);
    
    if (keys.length > 0) {
      await this.redis.del(...keys);
    }
  }
}
```

### 5. Background Jobs

```typescript
// src/jobs/data-collection.job.ts

import Bull from 'bull';
import { CongressAPIService } from '../services/data-collection/congress-api.service';

export function setupDataCollectionJobs(redisUrl: string) {
  const queue = new Bull('data-collection', redisUrl);

  // Daily bill collection
  queue.add(
    'collect-new-bills',
    {},
    {
      repeat: {
        cron: '0 2 * * *' // 2 AM daily
      }
    }
  );

  // Process job
  queue.process('collect-new-bills', async (job) => {
    const apiService = new CongressAPIService(process.env.CONGRESS_API_KEY!);
    const currentCongress = getCurrentCongress();
    
    job.progress(0);
    
    // Get bills updated in last 24 hours
    const bills = await apiService.getBillsForCongress(currentCongress);
    
    job.progress(50);
    
    // Save to database
    for (const bill of bills) {
      await saveBill(bill);
    }
    
    job.progress(100);
    
    return { billsProcessed: bills.length };
  });

  return queue;
}

function getCurrentCongress(): number {
  // Calculate current congress number
  const startYear = 1789;
  const currentYear = new Date().getFullYear();
  return Math.floor((currentYear - startYear) / 2) + 1;
}

async function saveBill(bill: any): Promise<void> {
  // Save to database
}
```

## API Endpoints

### Example REST API Routes

```typescript
// src/routes/politicians.ts

import express from 'express';
import { VotingAnalysisService } from '../services/analysis/voting-analysis.service';
import { CacheService } from '../services/cache/redis.service';

const router = express.Router();
const votingAnalysis = new VotingAnalysisService();
const cache = new CacheService(process.env.REDIS_URL!);

/**
 * GET /api/politicians/:id/metrics
 * Get comprehensive metrics for a politician
 */
router.get('/:id/metrics', async (req, res) => {
  try {
    const { id } = req.params;
    
    // Check cache first
    let metrics = await cache.getCachedPoliticianMetrics(id);
    
    if (!metrics) {
      // Calculate metrics
      metrics = {
        bipartisanScore: await votingAnalysis.calculateBipartisanScore(id),
        participationRate: await votingAnalysis.calculateParticipationRate(id),
        votingByPolicyArea: await votingAnalysis.analyzeVotingByPolicyArea(id),
        // Add more metrics
      };
      
      // Cache results
      await cache.cachePoliticianMetrics(id, metrics);
    }
    
    res.json(metrics);
  } catch (error) {
    res.status(500).json({ error: 'Failed to calculate metrics' });
  }
});

/**
 * GET /api/politicians/:id1/similarity/:id2
 * Calculate similarity between two politicians
 */
router.get('/:id1/similarity/:id2', async (req, res) => {
  try {
    const { id1, id2 } = req.params;
    
    // Check cache
    let similarity = await cache.getCachedSimilarity(id1, id2);
    
    if (similarity === null) {
      // Calculate similarity
      const votes1 = await getVotingRecords(id1);
      const votes2 = await getVotingRecords(id2);
      
      similarity = votingAnalysis.calculateVotingSimilarity(votes1, votes2);
      
      // Cache result
      await cache.cacheSimilarity(id1, id2, similarity);
    }
    
    res.json({ similarity });
  } catch (error) {
    res.status(500).json({ error: 'Failed to calculate similarity' });
  }
});

export default router;

// Placeholder
async function getVotingRecords(id: string) {
  return [];
}
```

## Testing Strategy

### Unit Tests

```typescript
// tests/unit/voting-analysis.test.ts

import { VotingAnalysisService } from '../../src/services/analysis/voting-analysis.service';

describe('VotingAnalysisService', () => {
  let service: VotingAnalysisService;

  beforeEach(() => {
    service = new VotingAnalysisService();
  });

  describe('calculateVotingSimilarity', () => {
    it('should return 1 for identical voting patterns', () => {
      const votes1 = [
        { position: 'Yea' },
        { position: 'Nay' },
        { position: 'Yea' }
      ];
      const votes2 = [
        { position: 'Yea' },
        { position: 'Nay' },
        { position: 'Yea' }
      ];

      const similarity = service.calculateVotingSimilarity(votes1, votes2);
      
      expect(similarity).toBeCloseTo(1.0);
    });

    it('should return -1 for opposite voting patterns', () => {
      const votes1 = [
        { position: 'Yea' },
        { position: 'Yea' },
        { position: 'Yea' }
      ];
      const votes2 = [
        { position: 'Nay' },
        { position: 'Nay' },
        { position: 'Nay' }
      ];

      const similarity = service.calculateVotingSimilarity(votes1, votes2);
      
      expect(similarity).toBeCloseTo(-1.0);
    });
  });
});
```

### Integration Tests

```typescript
// tests/integration/api.test.ts

import request from 'supertest';
import app from '../../src/index';

describe('Politicians API', () => {
  describe('GET /api/politicians/:id/metrics', () => {
    it('should return metrics for valid politician', async () => {
      const response = await request(app)
        .get('/api/politicians/B000944/metrics')
        .expect(200);

      expect(response.body).toHaveProperty('bipartisanScore');
      expect(response.body).toHaveProperty('participationRate');
      expect(response.body.bipartisanScore).toBeGreaterThanOrEqual(0);
      expect(response.body.bipartisanScore).toBeLessThanOrEqual(100);
    });

    it('should return 404 for invalid politician', async () => {
      await request(app)
        .get('/api/politicians/INVALID/metrics')
        .expect(404);
    });
  });
});
```

## Performance Optimization

### Database Indexing

```sql
-- Optimize queries for voting records
CREATE INDEX idx_voting_records_politician 
ON voting_records(politician_id, vote_date);

CREATE INDEX idx_voting_records_bill 
ON voting_records(bill_id);

-- Optimize bill searches
CREATE INDEX idx_bills_policy_area 
ON bills(policy_area);

CREATE INDEX idx_bills_congress 
ON bills(congress, bill_type, bill_number);

-- Full text search on bill titles
CREATE INDEX idx_bills_title_search 
ON bills USING gin(to_tsvector('english', title));
```

### Query Optimization

```typescript
// Use connection pooling
import { Pool } from 'pg';

const pool = new Pool({
  max: 20,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
});

// Batch queries
async function getBillsInBatch(billIds: string[]): Promise<Bill[]> {
  const query = `
    SELECT * FROM bills 
    WHERE id = ANY($1)
  `;
  
  const result = await pool.query(query, [billIds]);
  return result.rows;
}
```

## Deployment

### Environment Variables

```bash
# .env.example
NODE_ENV=production
PORT=8000

# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/political_db
DATABASE_POOL_SIZE=20

# Redis
REDIS_URL=redis://localhost:6379

# API Keys
CONGRESS_API_KEY=your_key_here
GOVINFO_API_KEY=your_key_here
OPENSTATES_API_KEY=your_key_here

# Security
JWT_SECRET=your_secret_here
CORS_ORIGIN=https://yourdomain.com

# Logging
LOG_LEVEL=info
```

### Docker Configuration

```dockerfile
# Dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

EXPOSE 8000

CMD ["node", "dist/index.js"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: ./backend
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://postgres:password@db:5432/political_db
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis

  db:
    image: postgres:14
    environment:
      - POSTGRES_DB=political_db
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

## Monitoring and Logging

### Logging Setup

```typescript
// src/utils/logger.ts

import winston from 'winston';

export const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
  transports: [
    new winston.transports.File({ filename: 'error.log', level: 'error' }),
    new winston.transports.File({ filename: 'combined.log' }),
  ],
});

if (process.env.NODE_ENV !== 'production') {
  logger.add(new winston.transports.Console({
    format: winston.format.simple(),
  }));
}
```

## Best Practices

1. **Always validate input**: Use validation middleware
2. **Cache expensive calculations**: Use Redis for caching
3. **Use background jobs**: For long-running tasks
4. **Implement rate limiting**: Protect APIs
5. **Monitor performance**: Track slow queries
6. **Version your APIs**: Support backward compatibility
7. **Document your code**: Use JSDoc/TSDoc
8. **Write tests**: Aim for >80% coverage
9. **Use TypeScript**: For type safety
10. **Follow security best practices**: Sanitize inputs, use HTTPS

## Revision History

- v1.0 (Initial) - Comprehensive implementation guide
