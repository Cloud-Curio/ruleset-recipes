# API Documentation

This document provides comprehensive information about the Political Social Network API.

## 🌐 Base URL

- **Development**: `http://localhost:8000`
- **Production**: `https://api.political-social-network.com`

## 📚 Interactive Documentation

The API documentation is available through Swagger UI:

- **Swagger UI**: `http://localhost:8000/api/docs`
- **OpenAPI JSON**: `http://localhost:8000/api/docs.json`

## 🔐 Authentication

The API uses JWT (JSON Web Token) based authentication.

## 🛡️ Security Features

### Rate Limiting

- **Window**: 15 minutes
- **Max Requests**: 100 per IP
- **Response**: 429 Too Many Requests when limit exceeded

### Security Headers (via Helmet)

- Content Security Policy
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- Referrer-Policy: strict-origin-when-cross-origin

### CORS

CORS is configured to allow requests from the frontend URL:
- Development: `http://localhost:3000`
- Production: Configured via `FRONTEND_URL` environment variable

## 📋 API Endpoints

For detailed endpoint documentation, visit the Swagger UI at `/api/docs`

### Main Endpoint Groups

- `/api/auth` - Authentication endpoints
- `/api/users` - User management
- `/api/politicians` - Political entity data
- `/api/bills` - Legislative bills
- `/api/votes` - Voting records
- `/api/social` - Social features
- `/api/analytics` - Analytics and insights

## 📊 Response Format

All API responses follow a consistent format with `success`, `data`, `error`, and optional `pagination` fields.

## 🔒 Security Reporting

If you discover a security vulnerability, please create an issue on GitHub with the "security" label.
