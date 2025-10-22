-- Database Views for Political Analysis Platform
-- These views provide optimized queries for common analytical operations

-- View: Politician Voting Summary
-- Aggregates voting statistics for each politician
CREATE OR REPLACE VIEW politician_voting_summary AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  p.chamber,
  COUNT(DISTINCT v.id) AS total_votes,
  COUNT(DISTINCT CASE WHEN v.vote_position = 'yes' THEN v.id END) AS yes_votes,
  COUNT(DISTINCT CASE WHEN v.vote_position = 'no' THEN v.id END) AS no_votes,
  COUNT(DISTINCT CASE WHEN v.vote_position = 'present' THEN v.id END) AS present_votes,
  COUNT(DISTINCT CASE WHEN v.vote_position = 'not_voting' THEN v.id END) AS not_voting,
  ROUND(
    (COUNT(DISTINCT CASE WHEN v.vote_position IN ('yes', 'no') THEN v.id END)::NUMERIC / 
    NULLIF(COUNT(DISTINCT v.id), 0)) * 100, 2
  ) AS participation_rate,
  COUNT(DISTINCT CASE WHEN v.party_line_vote = true THEN v.id END) AS party_line_votes,
  ROUND(
    (COUNT(DISTINCT CASE WHEN v.party_line_vote = true THEN v.id END)::NUMERIC / 
    NULLIF(COUNT(DISTINCT v.id), 0)) * 100, 2
  ) AS party_loyalty_percentage,
  MIN(v.vote_date) AS first_vote_date,
  MAX(v.vote_date) AS last_vote_date
FROM politicians p
LEFT JOIN votes v ON p.id = v.politician_id
WHERE p.in_office = true
GROUP BY p.id, p.full_name, p.party, p.state, p.chamber;

-- View: Bill Sponsorship Analysis
-- Shows bill sponsorship patterns and success rates
CREATE OR REPLACE VIEW bill_sponsorship_analysis AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  COUNT(DISTINCT b.id) AS bills_sponsored,
  COUNT(DISTINCT CASE WHEN b.status IN ('signed', 'enacted') THEN b.id END) AS bills_enacted,
  COUNT(DISTINCT CASE WHEN b.status = 'passed_house' THEN b.id END) AS passed_house,
  COUNT(DISTINCT CASE WHEN b.status = 'passed_senate' THEN b.id END) AS passed_senate,
  COUNT(DISTINCT CASE WHEN b.status = 'failed' THEN b.id END) AS bills_failed,
  ROUND(
    (COUNT(DISTINCT CASE WHEN b.status IN ('signed', 'enacted') THEN b.id END)::NUMERIC /
    NULLIF(COUNT(DISTINCT b.id), 0)) * 100, 2
  ) AS success_rate,
  AVG(b.bipartisan_support) AS avg_bipartisan_support,
  AVG(b.controversy_score) AS avg_controversy_score
FROM politicians p
LEFT JOIN bills b ON p.id = b.sponsor_id
GROUP BY p.id, p.full_name, p.party, p.state;

-- View: Social Media Activity Summary
-- Aggregates social media posting and engagement
CREATE OR REPLACE VIEW social_media_summary AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  COUNT(DISTINCT smp.id) AS total_posts,
  COUNT(DISTINCT CASE WHEN smp.platform = 'twitter' THEN smp.id END) AS twitter_posts,
  COUNT(DISTINCT CASE WHEN smp.platform = 'facebook' THEN smp.id END) AS facebook_posts,
  COUNT(DISTINCT CASE WHEN smp.platform = 'instagram' THEN smp.id END) AS instagram_posts,
  SUM(smp.likes_count) AS total_likes,
  SUM(smp.shares_count) AS total_shares,
  SUM(smp.comments_count) AS total_comments,
  ROUND(AVG(smp.likes_count), 0) AS avg_likes_per_post,
  ROUND(AVG(smp.shares_count), 0) AS avg_shares_per_post,
  COUNT(CASE WHEN smp.is_deleted = true THEN 1 END) AS deleted_posts,
  MIN(smp.posted_at) AS first_post_date,
  MAX(smp.posted_at) AS last_post_date
FROM politicians p
LEFT JOIN social_media_posts smp ON p.id = smp.politician_id
GROUP BY p.id, p.full_name, p.party, p.state;

-- View: NLP Analysis Summary
-- Aggregates NLP analysis results
CREATE OR REPLACE VIEW nlp_analysis_summary AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  COUNT(DISTINCT nlp.id) AS total_analyses,
  AVG(nlp.sentiment_score) AS avg_sentiment,
  AVG(nlp.toxicity_score) AS avg_toxicity,
  COUNT(CASE WHEN nlp.contains_hate_speech = true THEN 1 END) AS hate_speech_count,
  COUNT(CASE WHEN nlp.sentiment_label = 'positive' THEN 1 END) AS positive_content_count,
  COUNT(CASE WHEN nlp.sentiment_label = 'negative' THEN 1 END) AS negative_content_count,
  COUNT(CASE WHEN nlp.sentiment_label = 'neutral' THEN 1 END) AS neutral_content_count,
  AVG(nlp.profanity_score) AS avg_profanity,
  AVG(nlp.identity_attack_score) AS avg_identity_attack,
  AVG(nlp.insult_score) AS avg_insult,
  AVG(nlp.threat_score) AS avg_threat
FROM politicians p
LEFT JOIN nlp_analysis nlp ON p.id = nlp.politician_id
GROUP BY p.id, p.full_name, p.party, p.state;

-- View: Consistency Analysis Summary
-- Shows consistency scores across policy areas
CREATE OR REPLACE VIEW consistency_summary AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  AVG(ca.consistency_score) AS overall_consistency,
  AVG(ca.statement_vote_alignment) AS avg_statement_alignment,
  COUNT(DISTINCT ca.policy_topic) AS topics_analyzed,
  SUM(ca.total_votes_analyzed) AS total_votes_in_analysis,
  SUM(ca.total_statements_analyzed) AS total_statements_in_analysis,
  COUNT(CASE WHEN ca.consistency_score < 50 THEN 1 END) AS inconsistent_topics,
  COUNT(CASE WHEN ca.consistency_score >= 80 THEN 1 END) AS highly_consistent_topics
FROM politicians p
LEFT JOIN consistency_analysis ca ON p.id = ca.politician_id
GROUP BY p.id, p.full_name, p.party, p.state;

-- View: Campaign Promise Tracking
-- Summarizes campaign promise fulfillment
CREATE OR REPLACE VIEW campaign_promise_summary AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  COUNT(DISTINCT cp.id) AS total_promises,
  COUNT(CASE WHEN cp.status = 'kept' THEN 1 END) AS promises_kept,
  COUNT(CASE WHEN cp.status = 'broken' THEN 1 END) AS promises_broken,
  COUNT(CASE WHEN cp.status = 'in_progress' THEN 1 END) AS promises_in_progress,
  COUNT(CASE WHEN cp.status = 'compromised' THEN 1 END) AS promises_compromised,
  COUNT(CASE WHEN cp.status = 'stalled' THEN 1 END) AS promises_stalled,
  ROUND(
    (COUNT(CASE WHEN cp.status = 'kept' THEN 1 END)::NUMERIC /
    NULLIF(COUNT(DISTINCT cp.id), 0)) * 100, 2
  ) AS promise_keep_rate,
  AVG(cp.fulfillment_percentage) AS avg_fulfillment_percentage
FROM politicians p
LEFT JOIN campaign_promises cp ON p.id = cp.politician_id
GROUP BY p.id, p.full_name, p.party, p.state;

-- View: Comprehensive Politician Dashboard
-- Combines multiple metrics for a complete overview
CREATE OR REPLACE VIEW politician_dashboard AS
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  p.chamber,
  p.in_office,
  p.term_start,
  p.term_end,
  p.image_url,
  p.website,
  p.twitter_account,
  p.facebook_account,
  
  -- Voting metrics
  pvs.total_votes,
  pvs.participation_rate,
  pvs.party_loyalty_percentage,
  
  -- Bill metrics
  bsa.bills_sponsored,
  bsa.bills_enacted,
  bsa.success_rate AS bill_success_rate,
  bsa.avg_bipartisan_support,
  
  -- Social media metrics
  sms.total_posts AS social_media_posts,
  sms.total_likes,
  sms.avg_likes_per_post,
  sms.deleted_posts,
  
  -- NLP metrics
  nls.avg_sentiment,
  nls.avg_toxicity,
  nls.hate_speech_count,
  
  -- Consistency metrics
  cs.overall_consistency,
  cs.avg_statement_alignment,
  
  -- Promise metrics
  cps.total_promises,
  cps.promise_keep_rate,
  
  -- KPI scores (latest)
  (SELECT integrity_score FROM politician_kpis WHERE politician_id = p.id ORDER BY period_end DESC LIMIT 1) AS integrity_score,
  (SELECT honesty_score FROM politician_kpis WHERE politician_id = p.id ORDER BY period_end DESC LIMIT 1) AS honesty_score,
  (SELECT transparency_score FROM politician_kpis WHERE politician_id = p.id ORDER BY period_end DESC LIMIT 1) AS transparency_score,
  (SELECT effectiveness_score FROM politician_kpis WHERE politician_id = p.id ORDER BY period_end DESC LIMIT 1) AS effectiveness_score,
  (SELECT bipartisan_score FROM politician_kpis WHERE politician_id = p.id ORDER BY period_end DESC LIMIT 1) AS bipartisan_score
  
FROM politicians p
LEFT JOIN politician_voting_summary pvs ON p.id = pvs.politician_id
LEFT JOIN bill_sponsorship_analysis bsa ON p.id = bsa.politician_id
LEFT JOIN social_media_summary sms ON p.id = sms.politician_id
LEFT JOIN nlp_analysis_summary nls ON p.id = nls.politician_id
LEFT JOIN consistency_summary cs ON p.id = cs.politician_id
LEFT JOIN campaign_promise_summary cps ON p.id = cps.politician_id
WHERE p.in_office = true;

-- View: Bipartisan Collaboration Network
-- Shows cross-party collaboration patterns
CREATE OR REPLACE VIEW bipartisan_collaboration AS
SELECT 
  p1.id AS politician_1_id,
  p1.full_name AS politician_1_name,
  p1.party AS politician_1_party,
  p2.id AS politician_2_id,
  p2.full_name AS politician_2_name,
  p2.party AS politician_2_party,
  ps.voting_similarity,
  ps.policy_similarity,
  ps.overall_similarity,
  ps.common_votes,
  ps.agreement_count,
  ROUND((ps.agreement_count::NUMERIC / NULLIF(ps.common_votes, 0)) * 100, 2) AS agreement_percentage
FROM politician_similarity ps
JOIN politicians p1 ON ps.politician_1_id = p1.id
JOIN politicians p2 ON ps.politician_2_id = p2.id
WHERE p1.party != p2.party
  AND ps.overall_similarity > 0.5
ORDER BY ps.overall_similarity DESC;

-- View: Policy Topic Voting Patterns
-- Analyzes voting patterns by policy topic
CREATE OR REPLACE VIEW policy_topic_analysis AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  jsonb_array_elements_text(b.policy_areas::jsonb) AS policy_area,
  COUNT(DISTINCT v.id) AS votes_on_topic,
  COUNT(CASE WHEN v.vote_position = 'yes' THEN 1 END) AS yes_votes,
  COUNT(CASE WHEN v.vote_position = 'no' THEN 1 END) AS no_votes,
  ROUND(
    (COUNT(CASE WHEN v.vote_position = 'yes' THEN 1 END)::NUMERIC /
    NULLIF(COUNT(DISTINCT v.id), 0)) * 100, 2
  ) AS support_percentage
FROM politicians p
JOIN votes v ON p.id = v.politician_id
JOIN bills b ON v.bill_id = b.id
WHERE b.policy_areas IS NOT NULL
  AND jsonb_array_length(b.policy_areas::jsonb) > 0
GROUP BY p.id, p.full_name, p.party, p.state, policy_area;

-- View: Recent Alerts and Issues
-- Shows recent alerts for politicians
CREATE OR REPLACE VIEW recent_politician_alerts AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  pa.alert_type,
  pa.severity,
  pa.category,
  pa.title,
  pa.description,
  pa.alert_date,
  pa.is_verified
FROM politicians p
JOIN politician_alerts pa ON p.id = pa.politician_id
WHERE pa.is_dismissed = false
  AND pa.alert_date >= CURRENT_DATE - INTERVAL '30 days'
ORDER BY pa.alert_date DESC, pa.severity DESC;

-- View: Monthly KPI Trends
-- Tracks KPI changes over time
CREATE OR REPLACE VIEW monthly_kpi_trends AS
SELECT 
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.state,
  pk.period_start,
  pk.period_end,
  pk.integrity_score,
  pk.honesty_score,
  pk.consistency_score,
  pk.transparency_score,
  pk.effectiveness_score,
  pk.bipartisan_score,
  pk.attendance_rate,
  pk.bills_sponsored,
  pk.bills_passed,
  LAG(pk.integrity_score) OVER (PARTITION BY p.id ORDER BY pk.period_start) AS prev_integrity_score,
  LAG(pk.effectiveness_score) OVER (PARTITION BY p.id ORDER BY pk.period_start) AS prev_effectiveness_score,
  pk.integrity_score - LAG(pk.integrity_score) OVER (PARTITION BY p.id ORDER BY pk.period_start) AS integrity_change,
  pk.effectiveness_score - LAG(pk.effectiveness_score) OVER (PARTITION BY p.id ORDER BY pk.period_start) AS effectiveness_change
FROM politicians p
JOIN politician_kpis pk ON p.id = pk.politician_id
WHERE pk.period_type = 'monthly'
ORDER BY p.full_name, pk.period_start DESC;

-- View: State Comparison
-- Compares politicians within the same state
CREATE OR REPLACE VIEW state_politician_comparison AS
SELECT 
  p.state,
  p.id AS politician_id,
  p.full_name,
  p.party,
  p.chamber,
  pvs.participation_rate,
  pvs.party_loyalty_percentage,
  bsa.bills_sponsored,
  bsa.success_rate AS bill_success_rate,
  cs.overall_consistency,
  cps.promise_keep_rate,
  AVG(pvs.participation_rate) OVER (PARTITION BY p.state) AS state_avg_participation,
  AVG(cs.overall_consistency) OVER (PARTITION BY p.state) AS state_avg_consistency,
  RANK() OVER (PARTITION BY p.state ORDER BY pvs.participation_rate DESC) AS participation_rank,
  RANK() OVER (PARTITION BY p.state ORDER BY cs.overall_consistency DESC) AS consistency_rank
FROM politicians p
LEFT JOIN politician_voting_summary pvs ON p.id = pvs.politician_id
LEFT JOIN bill_sponsorship_analysis bsa ON p.id = bsa.politician_id
LEFT JOIN consistency_summary cs ON p.id = cs.politician_id
LEFT JOIN campaign_promise_summary cps ON p.id = cps.politician_id
WHERE p.in_office = true
ORDER BY p.state, participation_rank;
