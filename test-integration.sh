#!/bin/bash

# Integration Test Script
# Tests that the backend can serve the frontend and API endpoints work

set -e

echo "🧪 Running Integration Tests..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

test_count=0
pass_count=0

run_test() {
  test_count=$((test_count + 1))
  echo -n "  Test $test_count: $1... "
  if eval "$2" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PASS${NC}"
    pass_count=$((pass_count + 1))
    return 0
  else
    echo -e "${RED}✗ FAIL${NC}"
    return 1
  fi
}

echo "📦 Checking build artifacts..."
run_test "Backend compiled" "test -f backend/dist/index.js"
run_test "Frontend built" "test -f frontend/out/index.html"
run_test "Swagger config exists" "test -f backend/dist/config/swagger.js"
run_test "Routes compiled" "test -d backend/dist/routes"
run_test "Middleware compiled" "test -d backend/dist/middleware"

echo ""
echo "📚 Checking documentation..."
run_test "API docs exist" "test -f docs/API.md"
run_test "Deployment docs exist" "test -f docs/DEPLOYMENT.md"
run_test "Security docs exist" "test -f docs/SECURITY.md"
run_test "README updated" "grep -q 'Swagger' README.md"

echo ""
echo "🔧 Checking configuration..."
run_test "Frontend configured for export" "grep -q 'output.*export' frontend/next.config.js"
run_test "Backend has Swagger" "grep -q 'swagger' backend/dist/index.js"
run_test "Security headers configured" "grep -q 'helmet' backend/dist/index.js"

echo ""
echo "================================================"
echo "  Tests Passed: $pass_count / $test_count"
echo "================================================"

if [ $pass_count -eq $test_count ]; then
  echo -e "${GREEN}✅ All tests passed!${NC}"
  exit 0
else
  echo -e "${RED}❌ Some tests failed${NC}"
  exit 1
fi
