import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  // Social media posts table
  await knex.schema.createTable('social_media_posts', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.string('platform').notNullable(); // twitter, facebook, instagram, youtube
    table.string('platform_post_id').notNullable();
    table.text('content');
    table.string('media_type'); // text, image, video, link
    table.json('media_urls');
    table.string('post_url');
    table.timestamp('posted_at').notNullable();
    table.integer('likes_count').defaultTo(0);
    table.integer('shares_count').defaultTo(0);
    table.integer('comments_count').defaultTo(0);
    table.integer('retweets_count').defaultTo(0);
    table.boolean('is_deleted').defaultTo(false);
    table.timestamp('deleted_at');
    table.json('hashtags');
    table.json('mentions');
    table.json('raw_data'); // Store full API response
    table.timestamps(true, true);
    
    table.index(['politician_id']);
    table.index(['platform']);
    table.index(['posted_at']);
    table.index(['platform', 'platform_post_id']);
    table.unique(['platform', 'platform_post_id']);
  });

  // NLP analysis results table
  await knex.schema.createTable('nlp_analysis', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.string('content_type').notNullable(); // social_post, speech, statement, bill_text
    table.uuid('content_id'); // ID of the analyzed content
    table.text('analyzed_text');
    table.decimal('sentiment_score', 5, 4); // -1 to 1
    table.string('sentiment_label'); // positive, negative, neutral
    table.decimal('sentiment_magnitude', 5, 4); // 0 to 1
    table.json('topics'); // Extracted topics
    table.json('entities'); // Named entities (people, orgs, places)
    table.json('key_phrases');
    table.decimal('toxicity_score', 5, 4); // 0 to 1
    table.decimal('profanity_score', 5, 4);
    table.decimal('identity_attack_score', 5, 4);
    table.decimal('insult_score', 5, 4);
    table.decimal('threat_score', 5, 4);
    table.boolean('contains_hate_speech').defaultTo(false);
    table.string('hate_speech_target'); // demographic group if detected
    table.json('policy_positions'); // Extracted policy stances
    table.string('language').defaultTo('en');
    table.json('complexity_metrics'); // Reading level, etc.
    table.timestamp('analyzed_at').notNullable();
    table.timestamps(true, true);
    
    table.index(['politician_id']);
    table.index(['content_type']);
    table.index(['sentiment_label']);
    table.index(['contains_hate_speech']);
    table.index(['analyzed_at']);
  });

  // Text embeddings table (for semantic similarity)
  await knex.schema.createTable('text_embeddings', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('content_type').notNullable(); // bill, statement, post, speech
    table.uuid('content_id').notNullable();
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.string('embedding_model').notNullable(); // bert-base, sentence-transformers, etc.
    table.specificType('embedding_vector', 'vector(768)'); // PostgreSQL pgvector extension
    table.text('text_sample'); // First 500 chars for reference
    table.json('metadata');
    table.timestamps(true, true);
    
    table.index(['content_type', 'content_id']);
    table.index(['politician_id']);
    table.index(['embedding_model']);
  });

  // Politician consistency tracking
  await knex.schema.createTable('consistency_analysis', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.string('policy_topic').notNullable(); // healthcare, immigration, etc.
    table.date('analysis_period_start').notNullable();
    table.date('analysis_period_end').notNullable();
    table.decimal('consistency_score', 5, 2); // 0-100
    table.decimal('statement_vote_alignment', 5, 2); // 0-100
    table.integer('total_votes_analyzed');
    table.integer('total_statements_analyzed');
    table.json('inconsistencies'); // Specific contradictions found
    table.json('position_changes'); // Timeline of position shifts
    table.json('supporting_evidence');
    table.text('summary');
    table.timestamps(true, true);
    
    table.index(['politician_id']);
    table.index(['policy_topic']);
    table.index(['analysis_period_start']);
    table.index(['consistency_score']);
  });

  // Campaign promises tracking
  await knex.schema.createTable('campaign_promises', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.text('promise_text').notNullable();
    table.string('promise_category'); // healthcare, economy, etc.
    table.string('source_url'); // Where the promise was made
    table.date('promise_date');
    table.string('campaign_year');
    table.enum('status', [
      'not_started',
      'in_progress',
      'kept',
      'broken',
      'compromised',
      'stalled'
    ]).defaultTo('not_started');
    table.decimal('fulfillment_percentage', 5, 2).defaultTo(0);
    table.json('related_bills'); // Bills that address this promise
    table.json('related_votes'); // Votes related to this promise
    table.json('related_actions'); // Other legislative actions
    table.text('progress_notes');
    table.date('last_action_date');
    table.timestamps(true, true);
    
    table.index(['politician_id']);
    table.index(['status']);
    table.index(['promise_category']);
    table.index(['campaign_year']);
  });

  // Bias detection results
  await knex.schema.createTable('bias_analysis', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.string('bias_type').notNullable(); // demographic, geographic, economic, industry
    table.string('bias_category'); // specific category within type
    table.date('analysis_period_start').notNullable();
    table.date('analysis_period_end').notNullable();
    table.decimal('bias_score', 5, 2); // 0-100, higher = more biased
    table.string('bias_direction'); // favors/against specific group
    table.json('voting_patterns'); // Supporting vote data
    table.json('funding_correlations'); // Campaign finance correlations
    table.json('evidence'); // Specific examples
    table.integer('votes_analyzed');
    table.decimal('statistical_significance', 5, 4); // p-value
    table.text('summary');
    table.timestamps(true, true);
    
    table.index(['politician_id']);
    table.index(['bias_type']);
    table.index(['bias_score']);
    table.index(['analysis_period_start']);
  });

  // Politician similarity matrix
  await knex.schema.createTable('politician_similarity', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_1_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.uuid('politician_2_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.decimal('voting_similarity', 5, 4); // 0-1, cosine similarity
    table.decimal('policy_similarity', 5, 4); // Based on text embeddings
    table.decimal('overall_similarity', 5, 4); // Combined score
    table.integer('common_votes'); // Number of votes both participated in
    table.integer('agreement_count'); // How many times they voted the same
    table.json('divergent_topics'); // Topics where they differ most
    table.json('aligned_topics'); // Topics where they align
    table.date('calculated_at').notNullable();
    table.timestamps(true, true);
    
    table.index(['politician_1_id']);
    table.index(['politician_2_id']);
    table.index(['overall_similarity']);
    table.unique(['politician_1_id', 'politician_2_id']);
  });

  // Comprehensive KPI tracking
  await knex.schema.createTable('politician_kpis', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.date('period_start').notNullable();
    table.date('period_end').notNullable();
    table.string('period_type'); // weekly, monthly, quarterly, annual
    
    // Core metrics
    table.decimal('integrity_score', 5, 2); // 0-100
    table.decimal('honesty_score', 5, 2);
    table.decimal('consistency_score', 5, 2);
    table.decimal('transparency_score', 5, 2);
    table.decimal('effectiveness_score', 5, 2);
    table.decimal('bipartisan_score', 5, 2);
    
    // Engagement metrics
    table.integer('votes_cast');
    table.integer('votes_missed');
    table.decimal('attendance_rate', 5, 2);
    table.integer('bills_sponsored');
    table.integer('bills_cosponsored');
    table.integer('amendments_proposed');
    table.integer('speeches_given');
    
    // Social media metrics
    table.integer('social_posts_count');
    table.integer('total_engagement');
    table.decimal('average_toxicity', 5, 4);
    table.integer('controversial_posts');
    
    // Influence metrics
    table.decimal('leadership_score', 5, 2);
    table.decimal('influence_score', 5, 2);
    table.integer('bills_passed');
    table.integer('successful_amendments');
    
    // Constituent alignment
    table.decimal('constituent_alignment', 5, 2);
    table.decimal('district_approval_rating', 5, 2);
    
    table.timestamps(true, true);
    
    table.index(['politician_id']);
    table.index(['period_start', 'period_end']);
    table.index(['integrity_score']);
    table.index(['effectiveness_score']);
    table.unique(['politician_id', 'period_start', 'period_type']);
  });

  // Alert/notification triggers for significant events
  await knex.schema.createTable('politician_alerts', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.string('alert_type').notNullable(); // position_change, scandal, achievement, inconsistency
    table.string('severity'); // low, medium, high, critical
    table.string('category'); // voting, social_media, ethics, etc.
    table.text('title').notNullable();
    table.text('description');
    table.json('evidence'); // Supporting data
    table.string('source_url');
    table.boolean('is_verified').defaultTo(false);
    table.boolean('is_dismissed').defaultTo(false);
    table.timestamp('alert_date').notNullable();
    table.timestamps(true, true);
    
    table.index(['politician_id']);
    table.index(['alert_type']);
    table.index(['severity']);
    table.index(['alert_date']);
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTableIfExists('politician_alerts');
  await knex.schema.dropTableIfExists('politician_kpis');
  await knex.schema.dropTableIfExists('politician_similarity');
  await knex.schema.dropTableIfExists('bias_analysis');
  await knex.schema.dropTableIfExists('campaign_promises');
  await knex.schema.dropTableIfExists('consistency_analysis');
  await knex.schema.dropTableIfExists('text_embeddings');
  await knex.schema.dropTableIfExists('nlp_analysis');
  await knex.schema.dropTableIfExists('social_media_posts');
}
