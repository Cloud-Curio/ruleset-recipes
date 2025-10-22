import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('votes', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('vote_id').unique().notNullable(); // Congress.gov vote ID
    table.uuid('bill_id').references('id').inTable('bills').onDelete('CASCADE');
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.enum('chamber', ['house', 'senate']).notNullable();
    table.integer('congress').notNullable();
    table.integer('session').notNullable();
    table.integer('roll_call').notNullable();
    table.enum('vote_type', ['passage', 'amendment', 'procedural', 'nomination', 'other']);
    table.enum('vote_position', ['yes', 'no', 'present', 'not_voting']).notNullable();
    table.string('question'); // What was being voted on
    table.text('description');
    table.enum('result', ['passed', 'failed', 'agreed_to', 'rejected']);
    table.integer('total_yes').defaultTo(0);
    table.integer('total_no').defaultTo(0);
    table.integer('total_present').defaultTo(0);
    table.integer('total_not_voting').defaultTo(0);
    table.date('vote_date').notNullable();
    table.time('vote_time');
    table.string('congress_url'); // Link to vote record
    table.boolean('party_line_vote').defaultTo(false); // Calculated field
    table.decimal('bipartisan_index', 5, 2).defaultTo(0); // Cross-party voting measure
    table.timestamps(true, true);
    
    // Indexes
    table.index(['vote_id']);
    table.index(['bill_id']);
    table.index(['politician_id']);
    table.index(['chamber']);
    table.index(['vote_date']);
    table.index(['vote_position']);
    table.index(['party_line_vote']);
    
    // Composite indexes
    table.index(['politician_id', 'vote_date']);
    table.index(['bill_id', 'vote_position']);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('votes');
}

