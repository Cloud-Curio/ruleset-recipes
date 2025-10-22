import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('users', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.string('email').unique().notNullable();
    table.string('username').unique().notNullable();
    table.string('password_hash').notNullable();
    table.string('first_name').notNullable();
    table.string('last_name').notNullable();
    table.text('bio');
    table.string('avatar_url');
    table.string('location');
    table.string('website');
    table.enum('role', ['user', 'admin', 'moderator']).defaultTo('user');
    table.boolean('is_verified').defaultTo(false);
    table.boolean('is_active').defaultTo(true);
    table.json('preferences').defaultTo('{}');
    table.timestamp('last_login_at');
    table.timestamps(true, true);
    
    // Indexes
    table.index(['email']);
    table.index(['username']);
    table.index(['is_active']);
    table.index(['created_at']);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTable('users');
}

