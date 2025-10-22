import { Knex } from 'knex';

export async function up(knex: Knex): Promise<void> {
  // Posts table
  await knex.schema.createTable('posts', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('user_id').references('id').inTable('users').onDelete('CASCADE');
    table.uuid('politician_id').references('id').inTable('politicians').nullable();
    table.uuid('bill_id').references('id').inTable('bills').nullable();
    table.text('content').notNullable();
    table.json('media_urls').defaultTo('[]'); // Array of image/video URLs
    table.enum('post_type', ['text', 'image', 'video', 'poll', 'bill_share', 'vote_share']);
    table.json('poll_options').nullable(); // For poll posts
    table.integer('likes_count').defaultTo(0);
    table.integer('comments_count').defaultTo(0);
    table.integer('shares_count').defaultTo(0);
    table.boolean('is_pinned').defaultTo(false);
    table.boolean('is_deleted').defaultTo(false);
    table.timestamps(true, true);
    
    // Indexes
    table.index(['user_id']);
    table.index(['politician_id']);
    table.index(['bill_id']);
    table.index(['created_at']);
    table.index(['is_deleted']);
  });

  // Comments table
  await knex.schema.createTable('comments', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('post_id').references('id').inTable('posts').onDelete('CASCADE');
    table.uuid('user_id').references('id').inTable('users').onDelete('CASCADE');
    table.uuid('parent_comment_id').references('id').inTable('comments').nullable();
    table.text('content').notNullable();
    table.integer('likes_count').defaultTo(0);
    table.integer('replies_count').defaultTo(0);
    table.boolean('is_deleted').defaultTo(false);
    table.timestamps(true, true);
    
    // Indexes
    table.index(['post_id']);
    table.index(['user_id']);
    table.index(['parent_comment_id']);
    table.index(['created_at']);
  });

  // Likes table
  await knex.schema.createTable('likes', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('user_id').references('id').inTable('users').onDelete('CASCADE');
    table.uuid('post_id').references('id').inTable('posts').onDelete('CASCADE').nullable();
    table.uuid('comment_id').references('id').inTable('comments').onDelete('CASCADE').nullable();
    table.enum('like_type', ['like', 'dislike', 'love', 'angry', 'sad']);
    table.timestamps(true, true);
    
    // Ensure user can only like a post/comment once
    table.unique(['user_id', 'post_id']);
    table.unique(['user_id', 'comment_id']);
    
    // Indexes
    table.index(['user_id']);
    table.index(['post_id']);
    table.index(['comment_id']);
  });

  // Follows table (users following politicians)
  await knex.schema.createTable('follows', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('user_id').references('id').inTable('users').onDelete('CASCADE');
    table.uuid('politician_id').references('id').inTable('politicians').onDelete('CASCADE');
    table.boolean('notifications_enabled').defaultTo(true);
    table.timestamps(true, true);
    
    // Ensure user can only follow a politician once
    table.unique(['user_id', 'politician_id']);
    
    // Indexes
    table.index(['user_id']);
    table.index(['politician_id']);
  });

  // User follows (users following other users)
  await knex.schema.createTable('user_follows', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('follower_id').references('id').inTable('users').onDelete('CASCADE');
    table.uuid('following_id').references('id').inTable('users').onDelete('CASCADE');
    table.timestamps(true, true);
    
    // Ensure user can only follow another user once
    table.unique(['follower_id', 'following_id']);
    
    // Indexes
    table.index(['follower_id']);
    table.index(['following_id']);
  });

  // Notifications table
  await knex.schema.createTable('notifications', (table) => {
    table.uuid('id').primary().defaultTo(knex.raw('gen_random_uuid()'));
    table.uuid('user_id').references('id').inTable('users').onDelete('CASCADE');
    table.enum('type', [
      'like',
      'comment',
      'follow',
      'mention',
      'bill_update',
      'vote_alert',
      'politician_update'
    ]);
    table.string('title').notNullable();
    table.text('message');
    table.json('data').defaultTo('{}'); // Additional notification data
    table.boolean('is_read').defaultTo(false);
    table.timestamps(true, true);
    
    // Indexes
    table.index(['user_id']);
    table.index(['type']);
    table.index(['is_read']);
    table.index(['created_at']);
  });

  return Promise.resolve();
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTable('notifications');
  await knex.schema.dropTable('user_follows');
  await knex.schema.dropTable('follows');
  await knex.schema.dropTable('likes');
  await knex.schema.dropTable('comments');
  await knex.schema.dropTable('posts');
  return Promise.resolve();
}

