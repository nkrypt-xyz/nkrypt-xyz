#!/bin/bash

# Test Go Server Compliance Script
# This script helps verify that the Go rewrite is compliant with the Node.js implementation

set -e

echo "==================================="
echo "Go Server Compliance Test Runner"
echo "==================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if MinIO is running
echo "Checking MinIO..."
if ! curl -s http://localhost:9000/minio/health/live > /dev/null 2>&1; then
    echo -e "${RED}✗ MinIO is not running on localhost:9000${NC}"
    echo ""
    echo "Please start MinIO with:"
    echo "  cd ../web-server"
    echo "  docker-compose -f deploy/docker/docker-compose.dev.yml up -d minio"
    exit 1
fi
echo -e "${GREEN}✓ MinIO is running${NC}"

# Check if Go server is running
echo "Checking Go server..."
if ! curl -s http://localhost:9041/healthz > /dev/null 2>&1; then
    echo -e "${RED}✗ Go server is not running on localhost:9041${NC}"
    echo ""
    echo "Please start the Go server with:"
    echo "  cd ../web-server"
    echo "  make run"
    exit 1
fi
echo -e "${GREEN}✓ Go server is running${NC}"
echo ""

# Set MinIO environment variables
export MINIO_ENDPOINT=${MINIO_ENDPOINT:-localhost}
export MINIO_PORT=${MINIO_PORT:-9000}
export MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY:-minioadmin}
export MINIO_SECRET_KEY=${MINIO_SECRET_KEY:-minioadmin}
export MINIO_BUCKET_NAME=${MINIO_BUCKET_NAME:-nkrypt-blobs}
export MINIO_USE_SSL=${MINIO_USE_SSL:-false}

echo "MinIO Configuration:"
echo "  Endpoint: $MINIO_ENDPOINT:$MINIO_PORT"
echo "  Bucket: $MINIO_BUCKET_NAME"
echo "  SSL: $MINIO_USE_SSL"
echo ""

# Run tests
echo "Running compliance tests..."
echo "==================================="
echo ""

# Run all tests
if npm test; then
    echo ""
    echo "==================================="
    echo -e "${GREEN}✓ All compliance tests passed!${NC}"
    echo "==================================="
    echo ""
    echo "The Go server is compliant with the Node.js implementation."
    exit 0
else
    echo ""
    echo "==================================="
    echo -e "${RED}✗ Some compliance tests failed${NC}"
    echo "==================================="
    echo ""
    echo "There are discrepancies between the Go and Node.js implementations."
    echo "Review the test output above for details."
    exit 1
fi
