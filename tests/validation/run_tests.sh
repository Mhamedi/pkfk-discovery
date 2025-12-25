#!/bin/bash
set -e

echo "Starting validation test suite..."

# Start test database
echo "Starting test PostgreSQL database..."
docker-compose -f docker-compose.test.yml up -d postgres-test

# Wait for database to be ready
echo "Waiting for database to be ready..."
timeout=30
counter=0
until docker-compose -f docker-compose.test.yml exec -T postgres-test pg_isready -U test_user -d test_db > /dev/null 2>&1; do
    sleep 1
    counter=$((counter + 1))
    if [ $counter -ge $timeout ]; then
        echo "Database failed to start"
        docker-compose -f docker-compose.test.yml logs postgres-test
        exit 1
    fi
done

echo "Database is ready"

# Setup test database schema
echo "Setting up test database schema..."
docker-compose -f docker-compose.test.yml exec -T postgres-test psql -U test_user -d test_db < setup_test_db.sql

echo "Test database setup complete!"
echo ""
echo "Connection details:"
echo "  Host: localhost"
echo "  Port: 5433"
echo "  Database: test_db"
echo "  User: test_user"
echo "  Password: test_password"
echo ""
echo "To stop the test database, run:"
echo "  docker-compose -f docker-compose.test.yml down"

