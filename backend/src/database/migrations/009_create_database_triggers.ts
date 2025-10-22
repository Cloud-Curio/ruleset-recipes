import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  // Trigger to automatically update politician stats when votes are inserted/updated
  await knex.raw(`
    CREATE OR REPLACE FUNCTION update_politician_vote_stats()
    RETURNS TRIGGER AS $$
    BEGIN
      -- Update voting statistics on politicians table
      UPDATE politicians p
      SET 
        attendance_rate = (
          SELECT ROUND(
            (COUNT(CASE WHEN vote_position IN ('yes', 'no') THEN 1 END)::numeric / 
            NULLIF(COUNT(*), 0)) * 100, 2
          )
          FROM votes
          WHERE politician_id = NEW.politician_id
        )
      WHERE p.id = NEW.politician_id;
      
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER politician_vote_stats_trigger
    AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW
    EXECUTE FUNCTION update_politician_vote_stats();
  `);

  // Trigger to update bill statistics when votes are recorded
  await knex.raw(`
    CREATE OR REPLACE FUNCTION update_bill_vote_stats()
    RETURNS TRIGGER AS $$
    BEGIN
      IF NEW.bill_id IS NOT NULL THEN
        UPDATE bills b
        SET 
          total_votes = (SELECT COUNT(*) FROM votes WHERE bill_id = NEW.bill_id),
          yes_votes = (SELECT COUNT(*) FROM votes WHERE bill_id = NEW.bill_id AND vote_position = 'yes'),
          no_votes = (SELECT COUNT(*) FROM votes WHERE bill_id = NEW.bill_id AND vote_position = 'no'),
          abstain_votes = (SELECT COUNT(*) FROM votes WHERE bill_id = NEW.bill_id AND vote_position = 'present')
        WHERE b.id = NEW.bill_id;
      END IF;
      
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER bill_vote_stats_trigger
    AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW
    EXECUTE FUNCTION update_bill_vote_stats();
  `);

  // Trigger to calculate bipartisan support for bills
  await knex.raw(`
    CREATE OR REPLACE FUNCTION calculate_bipartisan_support()
    RETURNS TRIGGER AS $$
    DECLARE
      total_yes INTEGER;
      republican_yes INTEGER;
      democrat_yes INTEGER;
      bipartisan_score NUMERIC;
    BEGIN
      IF NEW.bill_id IS NOT NULL THEN
        SELECT 
          COUNT(CASE WHEN v.vote_position = 'yes' THEN 1 END),
          COUNT(CASE WHEN v.vote_position = 'yes' AND p.party = 'Republican' THEN 1 END),
          COUNT(CASE WHEN v.vote_position = 'yes' AND p.party = 'Democratic' THEN 1 END)
        INTO total_yes, republican_yes, democrat_yes
        FROM votes v
        JOIN politicians p ON v.politician_id = p.id
        WHERE v.bill_id = NEW.bill_id;
        
        -- Calculate bipartisan score (both parties represented in yes votes)
        IF total_yes > 0 AND republican_yes > 0 AND democrat_yes > 0 THEN
          bipartisan_score = LEAST(republican_yes, democrat_yes)::numeric / total_yes * 100;
        ELSE
          bipartisan_score = 0;
        END IF;
        
        UPDATE bills
        SET bipartisan_support = bipartisan_score
        WHERE id = NEW.bill_id;
      END IF;
      
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER bipartisan_support_trigger
    AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW
    EXECUTE FUNCTION calculate_bipartisan_support();
  `);

  // Trigger to detect party-line votes
  await knex.raw(`
    CREATE OR REPLACE FUNCTION detect_party_line_vote()
    RETURNS TRIGGER AS $$
    DECLARE
      politician_party TEXT;
      majority_party_position TEXT;
      is_party_line BOOLEAN;
    BEGIN
      -- Get politician's party
      SELECT party INTO politician_party
      FROM politicians
      WHERE id = NEW.politician_id;
      
      -- Determine if this is a party-line vote
      -- A vote is party-line if >80% of party members voted the same way
      WITH party_votes AS (
        SELECT 
          p.party,
          v.vote_position,
          COUNT(*) as vote_count,
          SUM(COUNT(*)) OVER (PARTITION BY p.party) as party_total
        FROM votes v
        JOIN politicians p ON v.politician_id = p.id
        WHERE v.vote_id = NEW.vote_id
          AND p.party = politician_party
        GROUP BY p.party, v.vote_position
      )
      SELECT 
        CASE 
          WHEN MAX(vote_count::numeric / party_total) > 0.8 THEN TRUE
          ELSE FALSE
        END
      INTO is_party_line
      FROM party_votes
      WHERE party = politician_party
        AND vote_position = NEW.vote_position;
      
      NEW.party_line_vote = COALESCE(is_party_line, FALSE);
      
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER party_line_vote_trigger
    BEFORE INSERT OR UPDATE ON votes
    FOR EACH ROW
    EXECUTE FUNCTION detect_party_line_vote();
  `);

  // Trigger to update politician bill counts
  await knex.raw(`
    CREATE OR REPLACE FUNCTION update_politician_bill_stats()
    RETURNS TRIGGER AS $$
    BEGIN
      UPDATE politicians
      SET 
        bills_sponsored = (
          SELECT COUNT(*)
          FROM bills
          WHERE sponsor_id = NEW.sponsor_id
        )
      WHERE id = NEW.sponsor_id;
      
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER politician_bill_stats_trigger
    AFTER INSERT OR UPDATE ON bills
    FOR EACH ROW
    EXECUTE FUNCTION update_politician_bill_stats();
  `);

  // Trigger to create alerts for significant events
  await knex.raw(`
    CREATE OR REPLACE FUNCTION create_consistency_alert()
    RETURNS TRIGGER AS $$
    BEGIN
      -- Create alert if consistency score drops significantly
      IF NEW.consistency_score < 50 AND 
         NEW.total_votes_analyzed >= 10 AND
         NEW.total_statements_analyzed >= 10 THEN
        
        INSERT INTO politician_alerts (
          politician_id,
          alert_type,
          severity,
          category,
          title,
          description,
          evidence,
          alert_date
        ) VALUES (
          NEW.politician_id,
          'inconsistency',
          CASE 
            WHEN NEW.consistency_score < 30 THEN 'high'
            WHEN NEW.consistency_score < 40 THEN 'medium'
            ELSE 'low'
          END,
          'voting',
          'Low Consistency Detected',
          'Politician shows low consistency between statements and voting record on ' || NEW.policy_topic,
          jsonb_build_object(
            'consistency_score', NEW.consistency_score,
            'policy_topic', NEW.policy_topic,
            'inconsistencies', NEW.inconsistencies
          ),
          CURRENT_TIMESTAMP
        )
        ON CONFLICT DO NOTHING;
      END IF;
      
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER consistency_alert_trigger
    AFTER INSERT OR UPDATE ON consistency_analysis
    FOR EACH ROW
    EXECUTE FUNCTION create_consistency_alert();
  `);

  // Trigger to create alerts for high toxicity
  await knex.raw(`
    CREATE OR REPLACE FUNCTION create_toxicity_alert()
    RETURNS TRIGGER AS $$
    BEGIN
      -- Create alert if toxicity is high or hate speech detected
      IF NEW.toxicity_score > 0.7 OR NEW.contains_hate_speech = true THEN
        
        INSERT INTO politician_alerts (
          politician_id,
          alert_type,
          severity,
          category,
          title,
          description,
          evidence,
          alert_date
        ) VALUES (
          NEW.politician_id,
          CASE 
            WHEN NEW.contains_hate_speech THEN 'hate_speech'
            ELSE 'toxic_content'
          END,
          CASE 
            WHEN NEW.contains_hate_speech THEN 'critical'
            WHEN NEW.toxicity_score > 0.9 THEN 'high'
            ELSE 'medium'
          END,
          'social_media',
          CASE 
            WHEN NEW.contains_hate_speech THEN 'Hate Speech Detected'
            ELSE 'High Toxicity Content'
          END,
          'Content with high toxicity score detected in ' || NEW.content_type,
          jsonb_build_object(
            'toxicity_score', NEW.toxicity_score,
            'content_type', NEW.content_type,
            'content_id', NEW.content_id,
            'hate_speech', NEW.contains_hate_speech,
            'target', NEW.hate_speech_target
          ),
          CURRENT_TIMESTAMP
        )
        ON CONFLICT DO NOTHING;
      END IF;
      
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER toxicity_alert_trigger
    AFTER INSERT ON nlp_analysis
    FOR EACH ROW
    EXECUTE FUNCTION create_toxicity_alert();
  `);

  // Trigger to update last_updated timestamp
  await knex.raw(`
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
      NEW.updated_at = CURRENT_TIMESTAMP;
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    -- Apply to all tables with updated_at column
    DO $$ 
    DECLARE
      table_name TEXT;
    BEGIN
      FOR table_name IN 
        SELECT tablename 
        FROM pg_tables 
        WHERE schemaname = 'public'
        AND tablename IN (
          'politicians', 'bills', 'votes', 'social_media_posts',
          'nlp_analysis', 'consistency_analysis', 'campaign_promises',
          'bias_analysis', 'politician_kpis'
        )
      LOOP
        EXECUTE format('
          CREATE TRIGGER update_%I_updated_at
          BEFORE UPDATE ON %I
          FOR EACH ROW
          EXECUTE FUNCTION update_updated_at_column();
        ', table_name, table_name);
      END LOOP;
    END $$;
  `);

  // Trigger to calculate controversy score for bills
  await knex.raw(`
    CREATE OR REPLACE FUNCTION calculate_bill_controversy()
    RETURNS TRIGGER AS $$
    DECLARE
      vote_margin NUMERIC;
      party_split NUMERIC;
      controversy NUMERIC;
    BEGIN
      IF NEW.bill_id IS NOT NULL THEN
        -- Calculate vote margin and party split
        SELECT 
          ABS(
            (SELECT COUNT(*) FROM votes WHERE bill_id = NEW.bill_id AND vote_position = 'yes') -
            (SELECT COUNT(*) FROM votes WHERE bill_id = NEW.bill_id AND vote_position = 'no')
          )::numeric / NULLIF(
            (SELECT COUNT(*) FROM votes WHERE bill_id = NEW.bill_id), 0
          ) as margin,
          ABS(
            (SELECT COUNT(*) FROM votes v 
             JOIN politicians p ON v.politician_id = p.id 
             WHERE v.bill_id = NEW.bill_id AND v.vote_position = 'yes' AND p.party = 'Republican') -
            (SELECT COUNT(*) FROM votes v 
             JOIN politicians p ON v.politician_id = p.id 
             WHERE v.bill_id = NEW.bill_id AND v.vote_position = 'yes' AND p.party = 'Democratic')
          )::numeric / NULLIF(
            (SELECT COUNT(*) FROM votes WHERE bill_id = NEW.bill_id AND vote_position = 'yes'), 0
          ) as split
        INTO vote_margin, party_split;
        
        -- Controversy is higher when margin is smaller and parties are split
        controversy = (1 - COALESCE(vote_margin, 0)) * 50 + COALESCE(party_split, 0) * 50;
        
        UPDATE bills
        SET controversy_score = ROUND(controversy, 2)
        WHERE id = NEW.bill_id;
      END IF;
      
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER bill_controversy_trigger
    AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW
    EXECUTE FUNCTION calculate_bill_controversy();
  `);
}

export async function down(knex: Knex): Promise<void> {
  // Drop all triggers
  await knex.raw(`
    DROP TRIGGER IF EXISTS politician_vote_stats_trigger ON votes;
    DROP TRIGGER IF EXISTS bill_vote_stats_trigger ON votes;
    DROP TRIGGER IF EXISTS bipartisan_support_trigger ON votes;
    DROP TRIGGER IF EXISTS party_line_vote_trigger ON votes;
    DROP TRIGGER IF EXISTS politician_bill_stats_trigger ON bills;
    DROP TRIGGER IF EXISTS consistency_alert_trigger ON consistency_analysis;
    DROP TRIGGER IF EXISTS toxicity_alert_trigger ON nlp_analysis;
    DROP TRIGGER IF EXISTS bill_controversy_trigger ON votes;
  `);

  // Drop update triggers for each table
  await knex.raw(`
    DROP TRIGGER IF EXISTS update_politicians_updated_at ON politicians;
    DROP TRIGGER IF EXISTS update_bills_updated_at ON bills;
    DROP TRIGGER IF EXISTS update_votes_updated_at ON votes;
    DROP TRIGGER IF EXISTS update_social_media_posts_updated_at ON social_media_posts;
    DROP TRIGGER IF EXISTS update_nlp_analysis_updated_at ON nlp_analysis;
    DROP TRIGGER IF EXISTS update_consistency_analysis_updated_at ON consistency_analysis;
    DROP TRIGGER IF EXISTS update_campaign_promises_updated_at ON campaign_promises;
    DROP TRIGGER IF EXISTS update_bias_analysis_updated_at ON bias_analysis;
    DROP TRIGGER IF EXISTS update_politician_kpis_updated_at ON politician_kpis;
  `);

  // Drop all functions
  await knex.raw(`
    DROP FUNCTION IF EXISTS update_politician_vote_stats();
    DROP FUNCTION IF EXISTS update_bill_vote_stats();
    DROP FUNCTION IF EXISTS calculate_bipartisan_support();
    DROP FUNCTION IF EXISTS detect_party_line_vote();
    DROP FUNCTION IF EXISTS update_politician_bill_stats();
    DROP FUNCTION IF EXISTS create_consistency_alert();
    DROP FUNCTION IF EXISTS create_toxicity_alert();
    DROP FUNCTION IF EXISTS update_updated_at_column();
    DROP FUNCTION IF EXISTS calculate_bill_controversy();
  `);
}
