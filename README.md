# Political Social Network Platform

A comprehensive social media platform designed specifically for political entities, combining Facebook-like functionality with advanced political data analytics.

## 🏛️ Overview

This platform creates detailed profiles for political entities using data from:
- **api.congress.gov** - Congressional data and voting records
- **govinfo.gov/api** - Government information and bulk data
- **openstates.org** - State-level political data

## ✨ Key Features

### Social Media Functionality
- 📱 Mobile-responsive Facebook-like interface
- 👥 User profiles and politician profiles
- 📝 Posts, comments, likes, and shares
- 🔔 Real-time notifications
- 📰 Personalized news feeds

### Political Analytics
- 🗳️ Voting record analysis and visualization
- 📊 Bill summaries using NLP processing
- 🔍 Cosine similarity calculations between politicians
- 📈 KPIs and statistical measures
- 🏷️ Automated bill categorization and binning

### Data Integration
- 🔄 Automated data ingestion from government APIs
- 📋 Comprehensive political entity profiles
- 🗂️ Membership and committee data
- 📜 Legislative history tracking

## 🛠️ Technology Stack

### Frontend
- **Next.js 14** with TypeScript
- **Tailwind CSS** for styling
- **React Query** for data fetching
- **Chart.js/D3.js** for data visualization
- **PWA** capabilities for mobile experience

### Backend
- **Node.js** with Express.js
- **TypeScript** for type safety
- **PostgreSQL** for primary database
- **Redis** for caching and sessions
- **Bull Queue** for background jobs

### Analytics & NLP
- **Natural Language Processing** for bill analysis
- **Cosine similarity** algorithms
- **Text summarization** and categorization
- **Statistical analysis** and KPI generation

### Infrastructure
- **Docker** containerization
- **GitHub Actions** CI/CD
- **AWS/Vercel** deployment
- **Monitoring** and logging

## 📁 Project Structure

```
├── frontend/                 # Next.js React application
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   ├── pages/          # Next.js pages
│   │   ├── hooks/          # Custom React hooks
│   │   ├── styles/         # CSS and styling
│   │   └── utils/          # Utility functions
│   └── package.json
├── backend/                 # Express.js API server
│   ├── src/
│   │   ├── routes/         # API route handlers
│   │   ├── controllers/    # Business logic
│   │   ├── services/       # External service integrations
│   │   ├── models/         # Database models
│   │   ├── middleware/     # Express middleware
│   │   ├── jobs/           # Background job processors
│   │   ├── analytics/      # Analytics and NLP processing
│   │   └── database/       # Database migrations and seeds
│   └── package.json
├── shared/                  # Shared types and utilities
├── docs/                   # Documentation
├── docker-compose.yml      # Development environment
└── package.json           # Root package.json
```

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- PostgreSQL 14+ (optional for full features)
- Redis 6+ (optional for caching)
- Docker (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd political-social-network
   ```

2. **Quick Start (Recommended)**
   ```bash
   ./start.sh
   ```
   
   This script will:
   - Install all dependencies
   - Build the entire project
   - Start the unified server
   
   Access the application at http://localhost:8000

3. **Manual Setup**
   ```bash
   # Install dependencies
   npm install
   
   # Build the project
   npm run build
   
   # Start the server
   cd backend && npm start
   ```

4. **Set up environment variables (optional)**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

### Running the Application

#### Option 1: Unified Server (Backend serves Frontend)

The simplest way to run the application:

```bash
cd backend
npm start
```

- **Frontend**: http://localhost:8000
- **API**: http://localhost:8000/api
- **API Docs**: http://localhost:8000/api/docs
- **Health Check**: http://localhost:8000/health

#### Option 2: Separate Development Servers

For development with hot-reload:

```bash
# Terminal 1 - Backend
cd backend && npm run dev

# Terminal 2 - Frontend
cd frontend && npm run dev
```

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8000

#### Option 3: Docker Compose

For full stack with database:

```bash
docker-compose up -d
```

This starts:
- Frontend (port 3000)
- Backend (port 8000)
- PostgreSQL (port 5432)
- Redis (port 6379)

## 📊 Data Sources

### Congress.gov API
- Bills and resolutions
- Voting records
- Committee information
- Member data

### GovInfo.gov API
- Federal Register documents
- Congressional documents
- Bulk data downloads

### OpenStates API
- State legislature data
- State bills and votes
- State legislator information

## 📖 Documentation

- **[API Documentation](docs/API.md)** - Complete API reference and examples
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Deployment options and configurations  
- **[Security Guide](docs/SECURITY.md)** - Security best practices and features
- **[Swagger UI](http://localhost:8000/api/docs)** - Interactive API documentation (when server is running)

## 🔧 Development

### Running Tests
```bash
npm test                    # Run all tests
npm run test:frontend      # Frontend tests only
npm run test:backend       # Backend tests only
```

### Database Operations
```bash
npm run migrate            # Run migrations
npm run migrate:rollback   # Rollback migrations
npm run seed              # Seed database
```

### Background Jobs
```bash
npm run jobs:start        # Start job processors
npm run jobs:dashboard    # View job dashboard
```

## 📈 Analytics Features

### NLP Processing
- **Bill Summarization**: Automatic extraction of key points from legislation
- **Topic Classification**: Categorization of bills by policy area
- **Sentiment Analysis**: Analysis of political statements and positions

### Similarity Calculations
- **Voting Pattern Analysis**: Compare politicians based on voting history
- **Policy Position Similarity**: Cosine similarity on policy stances
- **Coalition Detection**: Identify political alliances and opposition

### KPIs and Metrics
- **Participation Rates**: Voting attendance and engagement
- **Bipartisan Index**: Measure of cross-party collaboration
- **Influence Scores**: Impact and leadership metrics
- **Consistency Ratings**: Alignment with stated positions

## 🔐 Security

The platform implements comprehensive security measures:

- **JWT-based authentication** with secure token management
- **Rate limiting** to prevent abuse (100 requests per 15 minutes per IP)
- **Helmet.js** for security headers (CSP, XSS protection, etc.)
- **Input validation and sanitization** on all endpoints
- **CORS configuration** to control cross-origin requests
- **SQL injection prevention** through parameterized queries
- **XSS protection** via content security policy
- **OpenAPI/Swagger documentation** for API transparency

For detailed security information, see [SECURITY.md](docs/SECURITY.md)

## 📱 Mobile Optimization

- Responsive design for all screen sizes
- Touch-friendly interface elements
- Optimized loading for mobile networks
- Progressive Web App (PWA) features
- Offline capability for core features

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For questions and support:
- Review the [API Documentation](docs/API.md)
- Check the [Deployment Guide](docs/DEPLOYMENT.md)
- Read the [Security Guide](docs/SECURITY.md)
- Visit the [Swagger API Docs](http://localhost:8000/api/docs) (when running)
- Create an issue on GitHub
- Check the documentation in `/docs`

---

**Built with ❤️ for political transparency and civic engagement**

