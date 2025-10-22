import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('bills', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('congress_bill_id').unique().notNullable(); // e.g., "hr1-117"
    table.string('bill_number').notNullable(); // e.g., "H.R.1"
    table.integer('congress').notNullable(); // Congress session number
    table.enum('bill_type', ['hr', 's', 'hjres', 'sjres', 'hconres', 'sconres', 'hres', 'sres']);
    table.string('title').notNullable();
    table.text('short_title');
    table.text('official_title');
    table.text('summary');
    table.text('full_text'); // Full bill text
    table.uuid('sponsor_id').references('id').inTable('politicians');
    table.json('cosponsors').defaultTo('[]'); // Array of politician IDs
    table.enum('status', [
      'introduced',
      'referred',
      'reported',
      'passed_house',
      'passed_senate',
      'to_president',
      'signed',
      'vetoed',
      'failed'
    ]).defaultTo('introduced');
    table.json('actions').defaultTo('[]'); // Array of legislative actions
    table.json('committees').defaultTo('[]'); // Committees that reviewed the bill
    table.json('subjects').defaultTo('[]'); // Policy areas/subjects
    table.json('policy_areas').defaultTo('[]'); // NLP-derived policy classifications
    table.decimal('controversy_score', 5, 2).defaultTo(0); // Calculated controversy metric
    table.decimal('bipartisan_support', 5, 2).defaultTo(0); // Cross-party support percentage
    table.integer('total_votes').defaultTo(0);
    table.integer('yes_votes').defaultTo(0);
    table.integer('no_votes').defaultTo(0);
    table.integer('abstain_votes').defaultTo(0);
    table.date('introduced_date');
    table.date('last_action_date');
    table.string('congress_url'); // Link to Congress.gov
    table.string('govtrack_url'); // Link to GovTrack
    table.json('related_bills').defaultTo('[]'); // Related/companion bills
    table.text('nlp_summary'); // AI-generated summary
    table.json('key_topics').defaultTo('[]'); // NLP-extracted topics
    table.decimal('complexity_score', 5, 2).defaultTo(0); // Text complexity metric
    table.timestamps(true, true);
    
    // Indexes
    table.index(['congress_bill_id']);
    table.index(['bill_number']);
    table.index(['congress']);
    table.index(['bill_type']);
    table.index(['status']);
    table.index(['sponsor_id']);
    table.index(['introduced_date']);
    table.index(['controversy_score']);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('bills');
}

