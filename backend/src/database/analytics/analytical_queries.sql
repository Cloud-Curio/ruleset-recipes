-- Analytical SQL Scripts for Political Analysis
-- These scripts perform complex analytical queries for insights

-- ===================================
-- VOTING RECORD ANALYSIS
-- ===================================

-- 1. Identify Politicians Who Vote Against Their Party Most Often
-- This finds politicians with low party loyalty
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  p.chamber,
  COUNT(DISTINCT v.id) AS total_votes,
  COUNT(CASE WHEN v.party_line_vote = false THEN 1 END) AS cross_party_votes,
  ROUND(
    (COUNT(CASE WHEN v.party_line_vote = false THEN 1 END)::NUMERIC / 
    NULLIF(COUNT(DISTINCT v.id), 0)) * 100, 2
  ) AS cross_party_percentage,
  RANK() OVER (PARTITION BY p.party ORDER BY 
    (COUNT(CASE WHEN v.party_line_vote = false THEN 1 END)::NUMERIC / 
    NULLIF(COUNT(DISTINCT v.id), 0)) DESC
  ) AS party_rank
FROM politicians p
JOIN votes v ON p.id = v.politician_id
WHERE p.in_office = true
  AND v.vote_date >= CURRENT_DATE - INTERVAL '2 years'
GROUP BY p.id, p.full_name, p.party, p.state, p.chamber
HAVING COUNT(DISTINCT v.id) >= 50
ORDER BY cross_party_percentage DESC
LIMIT 50;

-- 2. Find Bills with Most Controversial Voting (Close Splits)
SELECT 
  b.id,
  b.congress_bill_id,
  b.title,
  b.bill_type,
  b.sponsor_id,
  b.policy_areas,
  v_summary.total_votes,
  v_summary.yes_votes,
  v_summary.no_votes,
  ABS(v_summary.yes_votes - v_summary.no_votes) AS vote_margin,
  ROUND(
    ABS(v_summary.yes_votes - v_summary.no_votes)::NUMERIC / 
    NULLIF(v_summary.total_votes, 0) * 100, 2
  ) AS margin_percentage,
  v_summary.republican_yes,
  v_summary.democrat_yes,
  b.status
FROM bills b
JOIN (
  SELECT 
    v.bill_id,
    COUNT(DISTINCT v.id) AS total_votes,
    COUNT(CASE WHEN v.vote_position = 'yes' THEN 1 END) AS yes_votes,
    COUNT(CASE WHEN v.vote_position = 'no' THEN 1 END) AS no_votes,
    COUNT(CASE WHEN v.vote_position = 'yes' AND p.party = 'Republican' THEN 1 END) AS republican_yes,
    COUNT(CASE WHEN v.vote_position = 'yes' AND p.party = 'Democratic' THEN 1 END) AS democrat_yes
  FROM votes v
  JOIN politicians p ON v.politician_id = p.id
  GROUP BY v.bill_id
) v_summary ON b.id = v_summary.bill_id
WHERE v_summary.total_votes >= 100
ORDER BY margin_percentage ASC
LIMIT 100;

-- 3. Voting Pattern Changes After Major Events
-- Compares voting before and after a specific date (e.g., election, scandal)
WITH politician_votes_before AS (
  SELECT 
    v.politician_id,
    COUNT(DISTINCT v.id) AS votes_before,
    COUNT(CASE WHEN v.vote_position = 'yes' THEN 1 END) AS yes_before,
    COUNT(CASE WHEN v.party_line_vote = true THEN 1 END) AS party_line_before
  FROM votes v
  WHERE v.vote_date < '2024-01-01' -- Change this date
    AND v.vote_date >= '2023-01-01'
  GROUP BY v.politician_id
),
politician_votes_after AS (
  SELECT 
    v.politician_id,
    COUNT(DISTINCT v.id) AS votes_after,
    COUNT(CASE WHEN v.vote_position = 'yes' THEN 1 END) AS yes_after,
    COUNT(CASE WHEN v.party_line_vote = true THEN 1 END) AS party_line_after
  FROM votes v
  WHERE v.vote_date >= '2024-01-01' -- Change this date
    AND v.vote_date < '2025-01-01'
  GROUP BY v.politician_id
)
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  ROUND((before.yes_before::NUMERIC / NULLIF(before.votes_before, 0)) * 100, 2) AS yes_rate_before,
  ROUND((after.yes_after::NUMERIC / NULLIF(after.votes_after, 0)) * 100, 2) AS yes_rate_after,
  ROUND((after.yes_after::NUMERIC / NULLIF(after.votes_after, 0)) * 100, 2) - 
    ROUND((before.yes_before::NUMERIC / NULLIF(before.votes_before, 0)) * 100, 2) AS yes_rate_change,
  ROUND((before.party_line_before::NUMERIC / NULLIF(before.votes_before, 0)) * 100, 2) AS party_loyalty_before,
  ROUND((after.party_line_after::NUMERIC / NULLIF(after.votes_after, 0)) * 100, 2) AS party_loyalty_after,
  ROUND((after.party_line_after::NUMERIC / NULLIF(after.votes_after, 0)) * 100, 2) - 
    ROUND((before.party_line_before::NUMERIC / NULLIF(before.votes_before, 0)) * 100, 2) AS loyalty_change
FROM politicians p
JOIN politician_votes_before before ON p.id = before.politician_id
JOIN politician_votes_after after ON p.id = after.politician_id
WHERE before.votes_before >= 20 AND after.votes_after >= 20
ORDER BY ABS(loyalty_change) DESC
LIMIT 50;

-- ===================================
-- CONSISTENCY ANALYSIS
-- ===================================

-- 4. Politicians with Largest Gaps Between Statements and Votes
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  ca.policy_topic,
  ca.consistency_score,
  ca.statement_vote_alignment,
  ca.total_votes_analyzed,
  ca.total_statements_analyzed,
  ca.inconsistencies,
  100 - ca.statement_vote_alignment AS alignment_gap
FROM politicians p
JOIN consistency_analysis ca ON p.id = ca.politician_id
WHERE ca.analysis_period_end >= CURRENT_DATE - INTERVAL '1 year'
  AND ca.total_votes_analyzed >= 5
  AND ca.total_statements_analyzed >= 5
ORDER BY alignment_gap DESC
LIMIT 100;

-- 5. Campaign Promises: Kept vs Broken Analysis
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  cp.promise_category,
  COUNT(*) AS total_promises,
  COUNT(CASE WHEN cp.status = 'kept' THEN 1 END) AS kept,
  COUNT(CASE WHEN cp.status = 'broken' THEN 1 END) AS broken,
  COUNT(CASE WHEN cp.status = 'in_progress' THEN 1 END) AS in_progress,
  ROUND(
    (COUNT(CASE WHEN cp.status = 'kept' THEN 1 END)::NUMERIC / 
    NULLIF(COUNT(*), 0)) * 100, 2
  ) AS keep_rate,
  ROUND(
    (COUNT(CASE WHEN cp.status = 'broken' THEN 1 END)::NUMERIC / 
    NULLIF(COUNT(*), 0)) * 100, 2
  ) AS break_rate
FROM politicians p
JOIN campaign_promises cp ON p.id = cp.politician_id
WHERE cp.campaign_year >= '2020'
GROUP BY p.id, p.full_name, p.party, p.state, cp.promise_category
ORDER BY total_promises DESC, keep_rate DESC;

-- ===================================
-- BIAS DETECTION
-- ===================================

-- 6. Detect Voting Bias by Demographic Impact
-- This identifies politicians who consistently vote against bills helping certain demographics
WITH demographic_votes AS (
  SELECT 
    v.politician_id,
    CASE 
      WHEN b.subjects::text ILIKE '%women%' OR b.subjects::text ILIKE '%gender%' THEN 'women'
      WHEN b.subjects::text ILIKE '%minority%' OR b.subjects::text ILIKE '%civil rights%' THEN 'minorities'
      WHEN b.subjects::text ILIKE '%labor%' OR b.subjects::text ILIKE '%worker%' THEN 'workers'
      WHEN b.subjects::text ILIKE '%veteran%' THEN 'veterans'
      WHEN b.subjects::text ILIKE '%senior%' OR b.subjects::text ILIKE '%elderly%' THEN 'seniors'
      ELSE 'other'
    END AS demographic_group,
    v.vote_position,
    b.title
  FROM votes v
  JOIN bills b ON v.bill_id = b.id
  WHERE v.vote_date >= CURRENT_DATE - INTERVAL '2 years'
)
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  dv.demographic_group,
  COUNT(*) AS total_votes,
  COUNT(CASE WHEN dv.vote_position = 'yes' THEN 1 END) AS yes_votes,
  COUNT(CASE WHEN dv.vote_position = 'no' THEN 1 END) AS no_votes,
  ROUND(
    (COUNT(CASE WHEN dv.vote_position = 'no' THEN 1 END)::NUMERIC / 
    NULLIF(COUNT(*), 0)) * 100, 2
  ) AS opposition_rate
FROM politicians p
JOIN demographic_votes dv ON p.id = dv.politician_id
WHERE dv.demographic_group != 'other'
GROUP BY p.id, p.full_name, p.party, p.state, dv.demographic_group
HAVING COUNT(*) >= 5
ORDER BY opposition_rate DESC;

-- 7. Analyze Funding Source Correlation with Voting
-- Requires campaign finance data integration
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  ba.bias_type,
  ba.bias_category,
  ba.bias_score,
  ba.bias_direction,
  ba.votes_analyzed,
  ba.statistical_significance,
  ba.funding_correlations
FROM politicians p
JOIN bias_analysis ba ON p.id = ba.politician_id
WHERE ba.bias_type = 'industry'
  AND ba.bias_score >= 60
  AND ba.statistical_significance < 0.05
ORDER BY ba.bias_score DESC;

-- ===================================
-- SOCIAL MEDIA ANALYSIS
-- ===================================

-- 8. Politicians with Highest Toxicity in Social Media
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  nls.avg_toxicity,
  nls.hate_speech_count,
  nls.avg_profanity,
  nls.avg_identity_attack,
  nls.avg_insult,
  nls.avg_threat,
  sms.total_posts,
  ROUND(
    (nls.hate_speech_count::NUMERIC / NULLIF(sms.total_posts, 0)) * 100, 2
  ) AS hate_speech_percentage
FROM politicians p
JOIN nlp_analysis_summary nls ON p.id = nls.politician_id
JOIN social_media_summary sms ON p.id = sms.politician_id
WHERE sms.total_posts >= 50
ORDER BY nls.avg_toxicity DESC
LIMIT 50;

-- 9. Deleted Posts Analysis (Transparency Issue)
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  COUNT(*) AS total_deleted,
  COUNT(CASE WHEN smp.platform = 'twitter' THEN 1 END) AS deleted_tweets,
  COUNT(CASE WHEN smp.platform = 'facebook' THEN 1 END) AS deleted_facebook,
  MIN(smp.posted_at) AS earliest_deletion,
  MAX(smp.posted_at) AS latest_deletion,
  AVG(EXTRACT(EPOCH FROM (smp.deleted_at - smp.posted_at)) / 3600) AS avg_hours_before_deletion
FROM politicians p
JOIN social_media_posts smp ON p.id = smp.politician_id
WHERE smp.is_deleted = true
  AND smp.deleted_at >= CURRENT_DATE - INTERVAL '1 year'
GROUP BY p.id, p.full_name, p.party, p.state
ORDER BY total_deleted DESC;

-- ===================================
-- COMPREHENSIVE SCORING
-- ===================================

-- 10. Overall Integrity Scorecard
SELECT 
  p.id,
  p.full_name,
  p.party,
  p.state,
  p.chamber,
  COALESCE(cs.overall_consistency, 0) AS consistency_score,
  COALESCE(cps.promise_keep_rate, 0) AS honesty_score,
  COALESCE(pvs.participation_rate, 0) AS transparency_score,
  100 - COALESCE(nls.avg_toxicity * 100, 0) AS civility_score,
  COALESCE(bsa.avg_bipartisan_support, 0) AS bipartisan_score,
  COALESCE(bsa.success_rate, 0) AS effectiveness_score,
  -- Overall integrity score (weighted average)
  ROUND(
    (COALESCE(cs.overall_consistency, 0) * 0.25 +
     COALESCE(cps.promise_keep_rate, 0) * 0.25 +
     COALESCE(pvs.participation_rate, 0) * 0.20 +
     (100 - COALESCE(nls.avg_toxicity * 100, 0)) * 0.15 +
     COALESCE(bsa.avg_bipartisan_support, 0) * 0.10 +
     COALESCE(bsa.success_rate, 0) * 0.05), 2
  ) AS overall_integrity_score
FROM politicians p
LEFT JOIN consistency_summary cs ON p.id = cs.politician_id
LEFT JOIN campaign_promise_summary cps ON p.id = cps.politician_id
LEFT JOIN politician_voting_summary pvs ON p.id = pvs.politician_id
LEFT JOIN nlp_analysis_summary nls ON p.id = nls.politician_id
LEFT JOIN bill_sponsorship_analysis bsa ON p.id = bsa.politician_id
WHERE p.in_office = true
ORDER BY overall_integrity_score DESC;
