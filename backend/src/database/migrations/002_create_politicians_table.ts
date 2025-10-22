import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('politicians', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('bioguide_id').unique(); // Congress.gov bioguide ID
    table.string('govtrack_id').unique(); // GovTrack ID
    table.string('opensecrets_id'); // OpenSecrets ID
    table.string('votesmart_id'); // Vote Smart ID
    table.string('fec_ids').nullable(); // JSON array of FEC IDs
    table.string('first_name').notNullable();
    table.string('last_name').notNullable();
    table.string('middle_name');
    table.string('suffix');
    table.string('nickname');
    table.string('full_name').notNullable();
    table.enum('gender', ['M', 'F', 'Other']);
    table.date('date_of_birth');
    table.string('party').notNullable(); // Democratic, Republican, Independent, etc.
    table.string('state').notNullable(); // Two-letter state code
    table.string('district'); // For House members
    table.enum('chamber', ['house', 'senate', 'governor', 'state_house', 'state_senate']);
    table.string('office'); // Current office title
    table.string('phone');
    table.string('website');
    table.string('contact_form');
    table.string('twitter_account');
    table.string('facebook_account');
    table.string('youtube_account');
    table.text('biography');
    table.string('image_url');
    table.boolean('in_office').defaultTo(true);
    table.date('term_start');
    table.date('term_end');
    table.integer('terms_served').defaultTo(0);
    table.json('committees').defaultTo('[]'); // Array of committee memberships
    table.json('leadership_roles').defaultTo('[]'); // Array of leadership positions
    table.decimal('influence_score', 5, 2).defaultTo(0); // Calculated influence metric
    table.decimal('bipartisan_score', 5, 2).defaultTo(0); // Bipartisan cooperation score
    table.decimal('attendance_rate', 5, 2).defaultTo(0); // Voting attendance percentage
    table.integer('bills_sponsored').defaultTo(0);
    table.integer('bills_cosponsored').defaultTo(0);
    table.json('policy_positions').defaultTo('{}'); // NLP-derived policy positions
    table.json('voting_patterns').defaultTo('{}'); // Aggregated voting statistics
    table.timestamp('last_updated');
    table.timestamps(true, true);
    
    // Indexes
    table.index(['bioguide_id']);
    table.index(['party']);
    table.index(['state']);
    table.index(['chamber']);
    table.index(['in_office']);
    table.index(['influence_score']);
    table.index(['last_name', 'first_name']);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('politicians');
}

