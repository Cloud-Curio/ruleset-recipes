import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import morgan from 'morgan';
import compression from 'compression';
import rateLimit from 'express-rate-limit';
import dotenv from 'dotenv';
import { createServer } from 'http';

// Import routes
import authRoutes from './routes/auth';
import userRoutes from './routes/users';
import politicianRoutes from './routes/politicians';
import billRoutes from './routes/bills';
import voteRoutes from './routes/votes';
import socialRoutes from './routes/social';
import analyticsRoutes from './routes/analytics';

// Import middleware
import { errorHandler } from './middleware/errorHandler';
import { notFound } from './middleware/notFound';
import { authMiddleware } from './middleware/auth';

// Import database
import { initializeDatabase } from './database/connection';

// Import services
import { initializeRedis } from './services/redis';
import { startBackgroundJobs } from './jobs/scheduler';

// Load environment variables
dotenv.config();

const app = express();
const PORT = process.env.PORT || 8000;

// Security middleware
app.use(helmet({
  contentSecurityPolicy: {
    directives: {
      defaultSrc: ["'self'"],
      styleSrc: ["'self'", "'unsafe-inline'"],
      scriptSrc: ["'self'"],
      imgSrc: ["'self'", "data:", "https:"],
    },
  },
}));

// CORS configuration
app.use(cors({
  origin: process.env.FRONTEND_URL || 'http://localhost:3000',
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization'],
}));

// Rate limiting
const limiter = rateLimit({
  windowMs: parseInt(process.env.RATE_LIMIT_WINDOW_MS || '900000'), // 15 minutes
  max: parseInt(process.env.RATE_LIMIT_MAX_REQUESTS || '100'),
  message: {
    error: 'Too many requests from this IP, please try again later.',
  },
  standardHeaders: true,
  legacyHeaders: false,
});
app.use('/api/', limiter);

// Body parsing middleware
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true, limit: '10mb' }));

// Compression middleware
app.use(compression());

// Logging middleware
app.use(morgan(process.env.NODE_ENV === 'production' ? 'combined' : 'dev'));

// Health check endpoint
app.get('/health', (req, res) => {
  res.status(200).json({
    status: 'OK',
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    environment: process.env.NODE_ENV,
  });
});

// API routes
app.use('/api/auth', authRoutes);
app.use('/api/users', authMiddleware, userRoutes);
app.use('/api/politicians', politicianRoutes);
app.use('/api/bills', billRoutes);
app.use('/api/votes', voteRoutes);
app.use('/api/social', authMiddleware, socialRoutes);
app.use('/api/analytics', analyticsRoutes);

// API documentation endpoint
app.get('/api', (req, res) => {
  res.json({
    name: 'Political Social Network API',
    version: '1.0.0',
    description: 'API for political social networking platform with advanced analytics',
    endpoints: {
      auth: '/api/auth',
      users: '/api/users',
      politicians: '/api/politicians',
      bills: '/api/bills',
      votes: '/api/votes',
      social: '/api/social',
      analytics: '/api/analytics',
    },
    documentation: '/api/docs',
    health: '/health',
  });
});

// Error handling middleware
app.use(notFound);
app.use(errorHandler);

// Initialize services and start server
async function startServer() {
  try {
    // Initialize database connection
    await initializeDatabase();
    console.log('✅ Database connected successfully');

    // Initialize Redis connection
    await initializeRedis();
    console.log('✅ Redis connected successfully');

    // Start background jobs
    if (process.env.NODE_ENV !== 'test') {
      await startBackgroundJobs();
      console.log('✅ Background jobs started successfully');
    }

    // Create HTTP server
    const server = createServer(app);

    // Start server
    server.listen(PORT, () => {
      console.log(`🚀 Server running on port ${PORT}`);
      console.log(`📊 Environment: ${process.env.NODE_ENV}`);
      console.log(`🌐 API URL: http://localhost:${PORT}/api`);
      console.log(`❤️  Health check: http://localhost:${PORT}/health`);
    });

    // Graceful shutdown
    process.on('SIGTERM', () => {
      console.log('SIGTERM received, shutting down gracefully');
      server.close(() => {
        console.log('Process terminated');
        process.exit(0);
      });
    });

    process.on('SIGINT', () => {
      console.log('SIGINT received, shutting down gracefully');
      server.close(() => {
        console.log('Process terminated');
        process.exit(0);
      });
    });

  } catch (error) {
    console.error('❌ Failed to start server:', error);
    process.exit(1);
  }
}

// Start the server
if (require.main === module) {
  startServer();
}

export default app;

