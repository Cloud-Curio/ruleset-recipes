import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  // Politician similarity scores
  await knex.schema.createTable('politician_similarities', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_a_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.uuid('politician_b_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.decimal('voting_similarity', 5, 4).notNullable(); // Cosine similarity score
    table.decimal('policy_similarity', 5, 4).notNullable(); // Policy position similarity
    table.decimal('overall_similarity', 5, 4).notNullable(); // Combined similarity score
    table.integer('common_votes').defaultTo(0); // Number of votes both participated in
    table.integer('agreement_count').defaultTo(0); // Number of times they voted the same way
    table.json('similarity_breakdown').defaultTo('{}'); // Detailed similarity metrics
    table.timestamp('calculated_at').defaultTo(knex.fn.now());
    table.timestamps(true, true);
    
    // Ensure no duplicate pairs
    table.unique(['politician_a_id', 'politician_b_id']);
    
    // Indexes
    table.index(['politician_a_id']);
    table.index(['politician_b_id']);
    table.index(['overall_similarity']);
    table.index(['calculated_at']);
  });

  // Bill topic classifications
  await knex.schema.createTable('bill_topics', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('bill_id').references('id').inTable('bills').onDelete('CASCADE');
    table.string('topic').notNullable(); // e.g., "healthcare", "environment"
    table.decimal('confidence', 5, 4).notNullable(); // NLP confidence score
    table.string('extraction_method'); // Which NLP method was used
    table.json('keywords').defaultTo('[]'); // Keywords that led to this classification
    table.timestamps(true, true);
    
    // Indexes
    table.index(['bill_id']);
    table.index(['topic']);
    table.index(['confidence']);
  });

  // Politician voting patterns
  await knex.schema.createTable('voting_patterns', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.string('topic').notNullable(); // Policy area
    table.integer('total_votes').defaultTo(0);
    table.integer('yes_votes').defaultTo(0);
    table.integer('no_votes').defaultTo(0);
    table.integer('abstain_votes').defaultTo(0);
    table.decimal('yes_percentage', 5, 2).defaultTo(0);
    table.decimal('party_alignment', 5, 2).defaultTo(0); // How often votes with party
    table.decimal('bipartisan_score', 5, 2).defaultTo(0); // Cross-party voting frequency
    table.json('recent_votes').defaultTo('[]'); // Last 10 votes in this topic
    table.timestamp('last_calculated').defaultTo(knex.fn.now());
    table.timestamps(true, true);
    
    // Ensure unique politician-topic combinations
    table.unique(['politician_id', 'topic']);
    
    // Indexes
    table.index(['politician_id']);
    table.index(['topic']);
    table.index(['party_alignment']);
    table.index(['bipartisan_score']);
  });

  // Daily analytics snapshots
  await knex.schema.createTable('daily_analytics', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.date('date').notNullable();
    table.string('metric_type').notNullable(); // e.g., "user_activity", "bill_activity"
    table.json('data').notNullable(); // Metric data
    table.timestamps(true, true);
    
    // Ensure unique date-metric combinations
    table.unique(['date', 'metric_type']);
    
    // Indexes
    table.index(['date']);
    table.index(['metric_type']);
  });

  // User engagement metrics
  await knex.schema.createTable('user_engagement', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('user_id').references('id').inTable('users').onDelete('CASCADE');
    table.date('date').notNullable();
    table.integer('posts_created').defaultTo(0);
    table.integer('comments_made').defaultTo(0);
    table.integer('likes_given').defaultTo(0);
    table.integer('politicians_followed').defaultTo(0);
    table.integer('bills_viewed').defaultTo(0);
    table.integer('votes_analyzed').defaultTo(0);
    table.integer('session_duration').defaultTo(0); // in minutes
    table.timestamps(true, true);
    
    // Ensure unique user-date combinations
    table.unique(['user_id', 'date']);
    
    // Indexes
    table.index(['user_id']);
    table.index(['date']);
  });

  // Trending topics
  await knex.schema.createTable('trending_topics', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('topic').notNullable();
    table.integer('mention_count').defaultTo(0);
    table.integer('engagement_score').defaultTo(0); // Weighted engagement metric
    table.decimal('growth_rate', 5, 2).defaultTo(0); // Percentage growth
    table.json('related_bills').defaultTo('[]'); // Bills related to this topic
    table.json('related_politicians').defaultTo('[]'); // Politicians associated with topic
    table.date('trending_date').notNullable();
    table.timestamps(true, true);
    
    // Ensure unique topic-date combinations
    table.unique(['topic', 'trending_date']);
    
    // Indexes
    table.index(['topic']);
    table.index(['trending_date']);
    table.index(['engagement_score']);
    table.index(['growth_rate']);
  });

  return Promise.resolve();
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTable('trending_topics');
  await knex.schema.dropTable('user_engagement');
  await knex.schema.dropTable('daily_analytics');
  await knex.schema.dropTable('voting_patterns');
  await knex.schema.dropTable('bill_topics');
  await knex.schema.dropTable('politician_similarities');
  return Promise.resolve();
}

