# Data Ingestion and Analytics System

## Overview

This comprehensive system ingests political data from multiple authoritative sources, analyzes it using advanced NLP and statistical methods, and provides actionable insights into political behavior, consistency, and integrity.

## Data Sources

### 1. Congress.gov API
- **Purpose**: Federal legislative data
- **Coverage**: Bills, votes, members, committees, amendments
- **Update Frequency**: Daily
- **Service**: `backend/src/services/ingestion/congressGovIngestion.ts`

#### Features
- Full member profiles with social media links
- Complete bill text and metadata
- Roll call vote records with individual positions
- Committee assignments and leadership roles
- Historical data back to 118th Congress

### 2. GovInfo.gov API
- **Purpose**: Government publications and records
- **Coverage**: Congressional Record, Federal Register, compiled bills, hearings, reports
- **Update Frequency**: Daily
- **Service**: `backend/src/services/ingestion/govInfoIngestion.ts`

#### Features
- Full-text Congressional Record with speeches
- Federal Register documents
- Committee reports and hearing transcripts
- Compiled bill versions (as introduced, engrossed, enrolled)
- Bulk data downloads for historical records

### 3. OpenStates API
- **Purpose**: State-level political data
- **Coverage**: All 50 states + DC and territories
- **Update Frequency**: Daily for active sessions
- **Service**: `backend/src/services/ingestion/openStatesIngestion.ts`

#### Features
- State legislator profiles
- State bills and resolutions
- State voting records
- Committee information
- Session data for current and historical sessions

## Database Schema

### Core Tables

#### Politicians
- Federal and state elected officials
- Biographical information
- Social media accounts
- Committee memberships
- Calculated KPI scores

#### Bills
- Federal and state legislation
- Sponsorship and co-sponsorship
- Legislative actions and status
- Policy area classifications
- NLP-derived summaries and topics

#### Votes
- Individual voting records
- Bill-level vote tallies
- Party-line vote indicators
- Vote context and descriptions

#### Social Media Posts
- Twitter/X, Facebook, Instagram, YouTube
- Post content and metadata
- Engagement metrics (likes, shares, comments)
- Deletion tracking for transparency

### Analytics Tables

#### NLP Analysis
- Sentiment scores
- Toxicity detection
- Topic extraction
- Entity recognition
- Policy position identification
- Hate speech detection

#### Consistency Analysis
- Statement-vote alignment
- Position changes over time
- Inconsistency detection and categorization
- Evidence gathering

#### Campaign Promises
- Promise tracking
- Fulfillment status
- Related legislative actions
- Timeline analysis

#### Bias Analysis
- Demographic voting patterns
- Industry influence
- Geographic bias
- Statistical significance testing

#### Text Embeddings
- Semantic similarity vectors
- Politician clustering
- Bill similarity matching
- Policy position clustering

#### KPI Metrics
- Integrity scores
- Honesty and consistency ratings
- Transparency measures
- Effectiveness metrics
- Bipartisan collaboration scores

## Analytics Services

### 1. NLP Analysis Service
**File**: `backend/src/services/analytics/nlpAnalysis.ts`

#### Capabilities
- **Sentiment Analysis**: Multi-method sentiment detection with magnitude
- **Topic Extraction**: TF-IDF and noun phrase analysis
- **Entity Recognition**: People, organizations, places
- **Toxicity Detection**: Profanity, insults, threats, identity attacks
- **Hate Speech Detection**: Target identification and severity classification
- **Policy Position Extraction**: Stance detection (support/oppose/neutral)
- **Complexity Metrics**: Reading level, Flesch score, vocabulary analysis

#### Usage
```typescript
import nlpAnalysis from './services/analytics/nlpAnalysis';

const result = await nlpAnalysis.analyzeText(
  text,
  'social_post',
  postId,
  politicianId
);
```

### 2. Consistency Analysis Service
**File**: `backend/src/services/analytics/consistencyAnalysis.ts`

#### Capabilities
- Statement-vote alignment scoring
- Inconsistency detection with severity levels
- Position change tracking
- Evidence aggregation
- Topic-specific consistency analysis

#### Metrics
- **Consistency Score**: 0-100 based on alignment rate
- **Statement-Vote Alignment**: Percentage of matching positions
- **Inconsistency Severity**: Low/Medium/High based on contradiction strength

#### Usage
```typescript
import consistencyAnalysis from './services/analytics/consistencyAnalysis';

const results = await consistencyAnalysis.analyzeConsistency(
  politicianId,
  startDate,
  endDate
);
```

### 3. KPI Aggregation Service
**File**: `backend/src/services/analytics/kpiAggregation.ts`

#### Core KPIs

1. **Integrity Score** (0-100)
   - Weighted combination of consistency, honesty, transparency, bias
   - Formula: `consistency × 0.3 + honesty × 0.3 + transparency × 0.25 + (100-bias) × 0.15`

2. **Honesty Score** (0-100)
   - Campaign promise keep rate
   - Statement-action alignment

3. **Consistency Score** (0-100)
   - Topic-level consistency across all policy areas
   - Position stability over time

4. **Transparency Score** (0-100)
   - Voting attendance rate (70%)
   - Communication frequency (30%)

5. **Effectiveness Score** (0-100)
   - Bill passage success rate
   - Amendment success rate
   - Leadership positions

6. **Bipartisan Score** (0-100)
   - Cross-party voting frequency
   - Bipartisan bill sponsorship
   - Cross-party collaboration

#### Usage
```typescript
import kpiAggregation from './services/analytics/kpiAggregation';

const kpis = await kpiAggregation.calculateKPIs(
  politicianId,
  periodStart,
  periodEnd,
  'monthly'
);
```

### 4. Feed Integration Service
**File**: `backend/src/services/analytics/feedIntegration.ts`

#### Features
- Unified timeline combining votes, bills, social posts, alerts
- Personalized feeds based on user preferences
- Trending items detection
- Comparative politician feeds
- Importance and engagement scoring

#### Feed Item Types
- **Votes**: With alignment analysis
- **Bills**: With sponsorship and status
- **Social Posts**: With NLP analysis results
- **Alerts**: Scandals, achievements, position changes

#### Usage
```typescript
import feedIntegration from './services/analytics/feedIntegration';

const feed = await feedIntegration.getFeed(
  {
    politicianIds: ['...'],
    startDate: new Date('2024-01-01'),
    minImportance: 0.5
  },
  page,
  pageSize
);
```

## SQL Views and Queries

### Analytical Views
**File**: `backend/src/database/views/analytical_views.sql`

#### Available Views

1. **politician_voting_summary**: Voting participation and patterns
2. **bill_sponsorship_analysis**: Bill success rates and patterns
3. **social_media_summary**: Social media activity aggregation
4. **nlp_analysis_summary**: NLP metrics aggregation
5. **consistency_summary**: Consistency across policy topics
6. **campaign_promise_summary**: Promise fulfillment tracking
7. **politician_dashboard**: Comprehensive KPI overview
8. **bipartisan_collaboration**: Cross-party collaboration network
9. **policy_topic_analysis**: Topic-specific voting patterns
10. **recent_politician_alerts**: Recent alerts and issues
11. **monthly_kpi_trends**: KPI changes over time
12. **state_politician_comparison**: Within-state comparisons

### Analytical Queries
**File**: `backend/src/database/analytics/analytical_queries.sql`

#### Query Categories

**Voting Record Analysis**
- Party loyalty identification
- Controversial bill detection
- Voting pattern changes after events

**Consistency Analysis**
- Statement-vote gaps
- Promise tracking
- Position flip-flops

**Bias Detection**
- Demographic voting patterns
- Industry funding correlations
- Geographic bias

**Social Media Analysis**
- Toxicity rankings
- Deleted post analysis
- Sentiment-vote correlation

**Effectiveness Metrics**
- Legislative success rates
- Bipartisan leadership
- Overall integrity scorecard

## Running the System

### Prerequisites
```bash
# Environment variables required
CONGRESS_API_KEY=your_key_here
GOVINFO_API_KEY=your_key_here
OPENSTATES_API_KEY=your_key_here
DATABASE_URL=postgresql://user:pass@host:5432/dbname
```

### Database Setup
```bash
# Run migrations
cd backend
npm run migrate

# Create views
psql -d political_social_network -f src/database/views/analytical_views.sql
```

### Data Ingestion

#### Full Initial Ingestion
```typescript
// Congress data
import congressGov from './services/ingestion/congressGovIngestion';
await congressGov.runFullIngestion(118); // Current Congress

// GovInfo data
import govInfo from './services/ingestion/govInfoIngestion';
await govInfo.runFullIngestion(118);

// OpenStates data (all states)
import openStates from './services/ingestion/openStatesIngestion';
await openStates.runFullIngestion();
```

#### Incremental Updates
```typescript
// Daily updates
await congressGov.incrementalUpdate(118);
await govInfo.ingestCongressionalRecord(yesterday, today);
await openStates.incrementalUpdate(); // Updates all states
```

### Analytics Execution

#### Run NLP Analysis on Social Media
```typescript
import nlpAnalysis from './services/analytics/nlpAnalysis';

const posts = await db('social_media_posts')
  .whereNull('analyzed_at')
  .limit(100);

for (const post of posts) {
  await nlpAnalysis.analyzeText(
    post.content,
    'social_post',
    post.id,
    post.politician_id
  );
}
```

#### Calculate Consistency Scores
```typescript
import consistencyAnalysis from './services/analytics/consistencyAnalysis';

const politicians = await db('politicians')
  .where('in_office', true);

for (const politician of politicians) {
  await consistencyAnalysis.analyzeConsistency(
    politician.id,
    oneYearAgo,
    today
  );
}
```

#### Update KPIs
```typescript
import kpiAggregation from './services/analytics/kpiAggregation';

// Monthly KPI calculation for all politicians
await kpiAggregation.calculateAllPoliticiansKPIs(
  monthStart,
  monthEnd,
  'monthly'
);
```

## API Endpoints (Recommended)

### Politicians
```
GET /api/politicians/:id
GET /api/politicians/:id/feed
GET /api/politicians/:id/kpis
GET /api/politicians/:id/consistency
GET /api/politicians/compare?ids=1,2,3
GET /api/politicians/rankings/:metric
```

### Feed
```
GET /api/feed/personalized
GET /api/feed/trending
GET /api/feed/by-politician/:id
POST /api/feed/filter
```

### Analytics
```
GET /api/analytics/voting-patterns/:politicianId
GET /api/analytics/consistency/:politicianId
GET /api/analytics/social-sentiment/:politicianId
GET /api/analytics/campaign-promises/:politicianId
```

## Scheduled Jobs

### Recommended Cron Schedule

```javascript
// Daily at 2 AM - Ingest new data
'0 2 * * *': async () => {
  await congressGov.incrementalUpdate();
  await govInfo.ingestCongressionalRecord();
  await openStates.incrementalUpdate();
}

// Daily at 3 AM - Run NLP analysis
'0 3 * * *': async () => {
  await nlpAnalysis.batchAnalyze(unanalyzedContent);
}

// Daily at 4 AM - Calculate consistency
'0 4 * * *': async () => {
  await consistencyAnalysis.analyzeConsistency();
}

// Monthly on 1st at 1 AM - Calculate KPIs
'0 1 1 * *': async () => {
  await kpiAggregation.calculateAllPoliticiansKPIs();
}
```

## Performance Optimization

### Indexing Strategy
- All foreign keys indexed
- Timestamp columns indexed for temporal queries
- JSON fields indexed with GIN for topic/policy searches
- Composite indexes for common query patterns

### Caching Recommendations
- Cache politician profiles (Redis, 1 hour TTL)
- Cache KPI scores (Redis, 24 hour TTL)
- Cache feed items (Redis, 15 minute TTL)
- Cache view results for popular queries

### Query Optimization
- Use materialized views for complex aggregations
- Implement pagination for all list endpoints
- Use database connection pooling
- Batch database operations when possible

## Security Considerations

### Data Privacy
- Only public data is ingested
- No personal contact information exposed
- Deleted social media posts tracked but content may be removed

### API Rate Limiting
- Congress.gov: 1000 requests/hour
- GovInfo.gov: 1000 requests/hour
- OpenStates: varies by plan
- Implement exponential backoff for all services

### Input Validation
- All user inputs sanitized
- SQL injection prevention via parameterized queries
- XSS prevention in feed content display

## Monitoring and Logging

### Key Metrics to Monitor
- Ingestion success rates
- API error rates
- NLP processing time
- Database query performance
- Feed generation latency

### Logging Strategy
- Info: Successful ingestion batches
- Warning: API rate limits approached
- Error: Failed ingestion attempts, analysis errors
- Debug: Detailed processing information

## Future Enhancements

### Planned Features
1. **Real-time streaming** for social media monitoring
2. **Machine learning models** for vote prediction
3. **Graph analytics** for influence networks
4. **Advanced embeddings** with BERT/GPT models
5. **Automated fact-checking** against statements
6. **Donor influence analysis** with campaign finance data
7. **Constituent survey integration** for approval ratings
8. **Multi-language support** for non-English content

### Research Applications
- Political science research datasets
- Voting pattern studies
- Social media impact analysis
- Predictive modeling for legislation outcomes
- Partisanship trend analysis

## Contributing

See main repository CONTRIBUTING.md for guidelines.

## License

MIT License - see LICENSE file

---

**Last Updated**: 2025-10-22
**Version**: 1.0
**Maintained by**: Political Social Network Team
