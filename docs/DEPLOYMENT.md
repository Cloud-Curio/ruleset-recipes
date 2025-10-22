# Deployment Guide

This guide explains how to deploy the Political Social Network platform in different configurations.

## 🏗️ Architecture Overview

The platform consists of three main components:

1. **Frontend** (Next.js) - Portable static site
2. **Backend** (Express.js) - API server
3. **Database** (PostgreSQL) + **Cache** (Redis)

## 📦 Deployment Options

### Option 1: Unified Deployment (Backend serves Frontend)

The backend can serve the frontend as static files. This is the simplest deployment option.

#### Build Steps

```bash
# Install dependencies
npm install

# Build everything
npm run build

# This will:
# 1. Build shared types
# 2. Build backend
# 3. Build frontend and export to static files
```

#### Start Server

```bash
cd backend
npm start

# Server runs on http://localhost:8000
# Frontend accessible at http://localhost:8000/
# API accessible at http://localhost:8000/api
# API Docs at http://localhost:8000/api/docs
```

#### Environment Variables

```bash
# Backend .env
PORT=8000
NODE_ENV=production
DATABASE_URL=postgresql://user:pass@localhost:5432/dbname
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
FRONTEND_URL=http://localhost:8000
```

### Option 2: Separate Deployment

Deploy frontend and backend separately for better scalability.

#### Frontend Deployment

The frontend is built as a static site and can be deployed to:
- Vercel
- Netlify
- AWS S3 + CloudFront
- GitHub Pages
- Any static hosting service

```bash
cd frontend
npm run build

# The 'out' directory contains the static site
# Upload to your hosting service
```

#### Backend Deployment

Deploy the backend to:
- AWS EC2/ECS
- Google Cloud Run
- Heroku
- DigitalOcean
- Any Node.js hosting service

```bash
cd backend
npm run build
npm start
```

## 🐳 Docker Deployment

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# Services:
# - Frontend: http://localhost:3000
# - Backend: http://localhost:8000
# - PostgreSQL: localhost:5432
# - Redis: localhost:6379
```

### Production Docker Build

```bash
# Build backend image
cd backend
docker build -t political-backend .

# Run backend
docker run -p 8000:8000 \
  -e DATABASE_URL=postgresql://... \
  -e REDIS_URL=redis://... \
  political-backend
```

## ☸️ Kubernetes Deployment

Example Kubernetes manifests are provided in `k8s/` directory:

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/ingress.yaml
```

## 🔐 Security Checklist

Before deploying to production:

- [ ] Set strong JWT_SECRET
- [ ] Use HTTPS (SSL/TLS certificates)
- [ ] Configure CORS properly
- [ ] Set secure database credentials
- [ ] Enable rate limiting
- [ ] Configure firewall rules
- [ ] Set up monitoring and logging
- [ ] Regular security updates
- [ ] Implement backup strategy
- [ ] Use environment variables for secrets

## 📊 Monitoring

### Health Checks

- Backend health: `http://your-domain:8000/health`
- API status: `http://your-domain:8000/api`

### Recommended Monitoring Tools

- Application: Sentry, New Relic, DataDog
- Infrastructure: Prometheus + Grafana
- Logs: ELK Stack, CloudWatch

## 🚀 Performance Optimization

### Backend

- Enable gzip compression (already configured)
- Use Redis caching
- Configure connection pooling for PostgreSQL
- Use CDN for static assets
- Implement API response caching

### Frontend

- Static export for fast loading
- PWA capabilities enabled
- Image optimization
- Code splitting
- Lazy loading

## 🔄 CI/CD Pipeline

Example GitHub Actions workflow:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: '18'
      - run: npm install
      - run: npm run build
      - run: npm test
      # Deploy steps...
```

## 🆘 Troubleshooting

### Frontend not loading

Check that:
1. Frontend was built: `cd frontend && npm run build`
2. Backend is serving static files
3. CORS is configured correctly

### API not accessible

Check that:
1. Backend is running
2. PORT is configured correctly
3. Firewall allows traffic
4. Database connection is established

### Database connection issues

Check that:
1. DATABASE_URL is correct
2. Database is running
3. Migrations are up to date

## 📚 Additional Resources

- [API Documentation](./API.md)
- [Security Best Practices](./SECURITY.md)
- [Main README](../README.md)
