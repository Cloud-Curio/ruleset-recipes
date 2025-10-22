import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  // Congressional Record table
  await knex.schema.createTable('congressional_records', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('package_id').unique().notNullable();
    table.text('title').notNullable();
    table.date('date_issued').notNullable();
    table.string('doc_class');
    table.string('collection').defaultTo('CREC');
    table.integer('congress');
    table.integer('session');
    table.string('pages');
    table.json('government_author');
    table.text('content_text');
    table.string('download_url');
    table.string('pdf_url');
    table.timestamp('last_modified');
    table.timestamps(true, true);
    
    table.index(['package_id']);
    table.index(['date_issued']);
    table.index(['congress']);
  });

  // Federal Register documents table
  await knex.schema.createTable('federal_register_documents', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('package_id').unique().notNullable();
    table.text('title').notNullable();
    table.date('publication_date').notNullable();
    table.string('volume');
    table.string('issue_number');
    table.string('collection').defaultTo('FR');
    table.string('doc_class');
    table.string('pages');
    table.json('agencies');
    table.json('subjects');
    table.string('download_url');
    table.string('pdf_url');
    table.timestamp('last_modified');
    table.timestamps(true, true);
    
    table.index(['package_id']);
    table.index(['publication_date']);
  });

  // Compiled bills table (from GovInfo)
  await knex.schema.createTable('compiled_bills', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('package_id').unique().notNullable();
    table.integer('congress').notNullable();
    table.string('bill_type').notNullable();
    table.integer('bill_number').notNullable();
    table.string('version'); // ih (introduced house), es (engrossed senate), etc.
    table.text('title').notNullable();
    table.date('date_issued');
    table.text('full_text');
    table.string('pages');
    table.string('download_url');
    table.string('pdf_url');
    table.string('xml_url');
    table.timestamp('last_modified');
    table.timestamps(true, true);
    
    table.index(['package_id']);
    table.index(['congress', 'bill_type', 'bill_number']);
  });

  // Committee reports table
  await knex.schema.createTable('committee_reports', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('package_id').unique().notNullable();
    table.integer('congress');
    table.string('chamber');
    table.string('report_type');
    table.string('report_number');
    table.text('title').notNullable();
    table.date('date_issued');
    table.string('committee');
    table.string('pages');
    table.string('download_url');
    table.string('pdf_url');
    table.timestamp('last_modified');
    table.timestamps(true, true);
    
    table.index(['package_id']);
    table.index(['congress']);
    table.index(['chamber']);
  });

  // Hearings table
  await knex.schema.createTable('hearings', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('package_id').unique().notNullable();
    table.integer('congress');
    table.string('chamber');
    table.text('title').notNullable();
    table.date('date_issued');
    table.string('committee');
    table.string('jacket_number');
    table.text('transcript_text');
    table.string('pages');
    table.string('download_url');
    table.string('pdf_url');
    table.timestamp('last_modified');
    table.timestamps(true, true);
    
    table.index(['package_id']);
    table.index(['congress']);
    table.index(['committee']);
    table.index(['date_issued']);
  });

  // State bills table (OpenStates)
  await knex.schema.createTable('state_bills', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('openstates_id').unique().notNullable();
    table.string('identifier').notNullable(); // e.g., "HB 1234"
    table.text('title').notNullable();
    table.json('classification'); // bill types
    table.json('subject');
    table.string('jurisdiction').notNullable(); // state code
    table.string('session');
    table.string('from_organization');
    table.json('sponsorships');
    table.json('abstracts');
    table.json('other_titles');
    table.json('actions');
    table.date('latest_action_date');
    table.text('latest_action_description');
    table.timestamp('created_at');
    table.timestamp('updated_at');
    table.date('first_action_date');
    table.date('latest_passage_date');
    table.json('versions');
    table.json('documents');
    table.json('sources');
    table.timestamps(true, true);
    
    table.index(['openstates_id']);
    table.index(['jurisdiction']);
    table.index(['identifier']);
    table.index(['session']);
  });

  // State votes table (OpenStates)
  await knex.schema.createTable('state_votes', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('openstates_id').unique().notNullable();
    table.uuid('bill_id').references('id').inTable('state_bills').onDelete('CASCADE');
    table.string('identifier');
    table.text('motion_text');
    table.json('motion_classification');
    table.date('start_date');
    table.string('result');
    table.string('organization');
    table.integer('yes_count').defaultTo(0);
    table.integer('no_count').defaultTo(0);
    table.integer('other_count').defaultTo(0);
    table.json('votes'); // individual member votes
    table.json('sources');
    table.timestamps(true, true);
    
    table.index(['openstates_id']);
    table.index(['bill_id']);
    table.index(['start_date']);
  });

  // State committees table (OpenStates)
  await knex.schema.createTable('state_committees', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('openstates_id').unique().notNullable();
    table.string('name').notNullable();
    table.string('classification');
    table.string('jurisdiction').notNullable();
    table.string('chamber');
    table.string('parent_id');
    table.json('memberships');
    table.json('sources');
    table.timestamps(true, true);
    
    table.index(['openstates_id']);
    table.index(['jurisdiction']);
    table.index(['chamber']);
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTableIfExists('state_committees');
  await knex.schema.dropTableIfExists('state_votes');
  await knex.schema.dropTableIfExists('state_bills');
  await knex.schema.dropTableIfExists('hearings');
  await knex.schema.dropTableIfExists('committee_reports');
  await knex.schema.dropTableIfExists('compiled_bills');
  await knex.schema.dropTableIfExists('federal_register_documents');
  await knex.schema.dropTableIfExists('congressional_records');
}
