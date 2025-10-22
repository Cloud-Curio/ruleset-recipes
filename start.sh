#!/bin/bash

# Quick Start Script
# Builds and starts the Political Social Network platform

set -e

echo "🚀 Political Social Network - Quick Start"
echo "=========================================="
echo ""

# Check if node is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+ first."
    exit 1
fi

NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 18 ]; then
    echo "❌ Node.js version 18+ is required. You have version $NODE_VERSION."
    exit 1
fi

echo "✅ Node.js $(node -v) detected"
echo ""

# Check if already built
if [ -f "backend/dist/index.js" ] && [ -f "frontend/out/index.html" ]; then
    echo "✅ Build artifacts found"
else
    echo "📦 Installing dependencies and building..."
    npm install
    npm run build
    echo ""
fi

echo "🎯 Starting server..."
echo ""
echo "The server will be available at:"
echo "  🌐 Frontend:  http://localhost:8000"
echo "  📡 API:       http://localhost:8000/api"
echo "  📚 API Docs:  http://localhost:8000/api/docs"
echo "  ❤️  Health:    http://localhost:8000/health"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

cd backend
npm start
