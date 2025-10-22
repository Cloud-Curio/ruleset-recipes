import swaggerJsdoc from 'swagger-jsdoc';

const options: swaggerJsdoc.Options = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'Political Social Network API',
      version: '1.0.0',
      description: `
        A comprehensive API for political social networking platform with advanced analytics.
        
        ## Features
        - User authentication and authorization
        - Political entity profiles (politicians, bills, votes)
        - Social networking functionality (posts, comments, follows)
        - Advanced analytics (voting patterns, similarity scores)
        - Real-time data integration from government APIs
        
        ## Security
        - JWT-based authentication
        - API rate limiting
        - Input validation and sanitization
        - CORS protection
        - Helmet security headers
        - XSS and SQL injection prevention
        
        ## Data Sources
        - Congress.gov API
        - GovInfo.gov API
        - OpenStates API
      `,
      contact: {
        name: 'API Support',
        url: 'https://github.com/political-social-network',
      },
      license: {
        name: 'MIT',
        url: 'https://opensource.org/licenses/MIT',
      },
    },
    servers: [
      {
        url: 'http://localhost:8000',
        description: 'Development server',
      },
      {
        url: 'https://api.political-social-network.com',
        description: 'Production server',
      },
    ],
    components: {
      securitySchemes: {
        bearerAuth: {
          type: 'http',
          scheme: 'bearer',
          bearerFormat: 'JWT',
          description: 'JWT authorization token',
        },
        apiKey: {
          type: 'apiKey',
          in: 'header',
          name: 'X-API-Key',
          description: 'API Key for server-to-server communication',
        },
      },
      schemas: {
        User: {
          type: 'object',
          properties: {
            id: { type: 'string', format: 'uuid' },
            email: { type: 'string', format: 'email' },
            username: { type: 'string' },
            firstName: { type: 'string' },
            lastName: { type: 'string' },
            bio: { type: 'string' },
            avatarUrl: { type: 'string', format: 'uri' },
            role: { type: 'string', enum: ['user', 'admin', 'moderator'] },
            isVerified: { type: 'boolean' },
            createdAt: { type: 'string', format: 'date-time' },
            updatedAt: { type: 'string', format: 'date-time' },
          },
        },
        Politician: {
          type: 'object',
          properties: {
            id: { type: 'string', format: 'uuid' },
            firstName: { type: 'string' },
            lastName: { type: 'string' },
            fullName: { type: 'string' },
            party: { type: 'string' },
            state: { type: 'string' },
            chamber: { type: 'string', enum: ['house', 'senate', 'governor', 'state_house', 'state_senate'] },
            office: { type: 'string' },
            imageUrl: { type: 'string', format: 'uri' },
            inOffice: { type: 'boolean' },
            influenceScore: { type: 'number', format: 'float' },
            bipartisanScore: { type: 'number', format: 'float' },
            attendanceRate: { type: 'number', format: 'float' },
          },
        },
        Bill: {
          type: 'object',
          properties: {
            id: { type: 'string', format: 'uuid' },
            billNumber: { type: 'string' },
            congress: { type: 'integer' },
            title: { type: 'string' },
            summary: { type: 'string' },
            status: {
              type: 'string',
              enum: ['introduced', 'referred', 'reported', 'passed_house', 'passed_senate', 'to_president', 'signed', 'vetoed', 'failed'],
            },
            introducedDate: { type: 'string', format: 'date' },
            controversyScore: { type: 'number', format: 'float' },
            bipartisanSupport: { type: 'number', format: 'float' },
          },
        },
        Vote: {
          type: 'object',
          properties: {
            id: { type: 'string', format: 'uuid' },
            voteId: { type: 'string' },
            chamber: { type: 'string', enum: ['house', 'senate'] },
            votePosition: { type: 'string', enum: ['yes', 'no', 'present', 'not_voting'] },
            voteDate: { type: 'string', format: 'date-time' },
            result: { type: 'string', enum: ['passed', 'failed', 'agreed_to', 'rejected'] },
          },
        },
        Post: {
          type: 'object',
          properties: {
            id: { type: 'string', format: 'uuid' },
            content: { type: 'string' },
            postType: { type: 'string', enum: ['text', 'image', 'video', 'poll', 'bill_share', 'vote_share'] },
            likesCount: { type: 'integer' },
            commentsCount: { type: 'integer' },
            sharesCount: { type: 'integer' },
            createdAt: { type: 'string', format: 'date-time' },
          },
        },
        Error: {
          type: 'object',
          properties: {
            success: { type: 'boolean', example: false },
            error: { type: 'string' },
            message: { type: 'string' },
          },
        },
        PaginatedResponse: {
          type: 'object',
          properties: {
            success: { type: 'boolean', example: true },
            data: { type: 'array', items: {} },
            pagination: {
              type: 'object',
              properties: {
                page: { type: 'integer' },
                limit: { type: 'integer' },
                total: { type: 'integer' },
                totalPages: { type: 'integer' },
                hasNext: { type: 'boolean' },
                hasPrev: { type: 'boolean' },
              },
            },
          },
        },
      },
      responses: {
        UnauthorizedError: {
          description: 'Access token is missing or invalid',
          content: {
            'application/json': {
              schema: { $ref: '#/components/schemas/Error' },
            },
          },
        },
        NotFoundError: {
          description: 'Resource not found',
          content: {
            'application/json': {
              schema: { $ref: '#/components/schemas/Error' },
            },
          },
        },
        ValidationError: {
          description: 'Invalid input',
          content: {
            'application/json': {
              schema: { $ref: '#/components/schemas/Error' },
            },
          },
        },
        ServerError: {
          description: 'Internal server error',
          content: {
            'application/json': {
              schema: { $ref: '#/components/schemas/Error' },
            },
          },
        },
      },
    },
    security: [
      {
        bearerAuth: [],
      },
    ],
  },
  apis: ['./src/routes/*.ts', './src/index.ts'],
};

export const swaggerSpec = swaggerJsdoc(options);
