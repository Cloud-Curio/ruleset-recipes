import knex, { Knex } from 'knex';
// eslint-disable-next-line @typescript-eslint/no-var-requires
const knexConfig = require('../../knexfile');

const environment = process.env['NODE_ENV'] || 'development';
const config = knexConfig[environment as keyof typeof knexConfig];

// Create database connection
export const db: Knex = knex(config);

// Initialize database connection and run migrations
export async function initializeDatabase(): Promise<void> {
  try {
    // Test the connection
    await db.raw('SELECT 1');
    
    // Run migrations in production/staging
    if (process.env['NODE_ENV'] === 'production' || process.env['NODE_ENV'] === 'staging') {
      await db.migrate.latest();
      console.log('✅ Database migrations completed');
    }
    
    console.log(`✅ Database connected (${environment})`);
  } catch (error) {
    console.error('❌ Database connection failed:', error);
    throw error;
  }
}

// Close database connection
export async function closeDatabase(): Promise<void> {
  try {
    await db.destroy();
    console.log('✅ Database connection closed');
  } catch (error) {
    console.error('❌ Error closing database connection:', error);
    throw error;
  }
}

export default db;

