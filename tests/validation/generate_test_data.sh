#!/bin/bash
# Script to generate complex test data for PK/FK discovery testing

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Generating complex test data for PK/FK discovery...${NC}"

# Container name from docker-compose
CONTAINER_NAME="pkfk-discover-postgres"
DB_NAME="pkfk_discovery"
DB_USER="pkfk"
DB_PASSWORD="pkfk_dev_password"

# Check if PostgreSQL container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "Error: PostgreSQL container ($CONTAINER_NAME) is not running."
    echo "Please start it with: docker compose -f deploy/docker-compose.yml up -d postgres"
    exit 1
fi

# Wait for PostgreSQL to be ready
echo -e "${BLUE}Waiting for PostgreSQL to be ready...${NC}"
until docker exec "$CONTAINER_NAME" pg_isready -U "$DB_USER" >/dev/null 2>&1; do
    echo "Waiting for PostgreSQL..."
    sleep 1
done

# Execute the SQL script directly in the container
echo -e "${BLUE}Executing test data generation script...${NC}"
docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" < "$(dirname "$0")/generate_complex_test_data.sql"

# Count records
echo -e "${BLUE}Counting generated records...${NC}"
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" <<EOF
SELECT 
    'locations' as table_name, COUNT(*) as record_count FROM locations
UNION ALL SELECT 'warehouses', COUNT(*) FROM warehouses
UNION ALL SELECT 'departments', COUNT(*) FROM departments
UNION ALL SELECT 'employees', COUNT(*) FROM employees
UNION ALL SELECT 'suppliers', COUNT(*) FROM suppliers
UNION ALL SELECT 'categories', COUNT(*) FROM categories
UNION ALL SELECT 'products', COUNT(*) FROM products
UNION ALL SELECT 'product_variants', COUNT(*) FROM product_variants
UNION ALL SELECT 'customers', COUNT(*) FROM customers
UNION ALL SELECT 'customer_addresses', COUNT(*) FROM customer_addresses
UNION ALL SELECT 'orders', COUNT(*) FROM orders
UNION ALL SELECT 'order_items', COUNT(*) FROM order_items
UNION ALL SELECT 'supplier_products', COUNT(*) FROM supplier_products
UNION ALL SELECT 'reviews', COUNT(*) FROM reviews
UNION ALL SELECT 'shipments', COUNT(*) FROM shipments
UNION ALL SELECT 'shipment_items', COUNT(*) FROM shipment_items
UNION ALL SELECT 'inventory', COUNT(*) FROM inventory
UNION ALL SELECT 'price_history', COUNT(*) FROM price_history
UNION ALL SELECT 'promotions', COUNT(*) FROM promotions
UNION ALL SELECT 'promotion_products', COUNT(*) FROM promotion_products
ORDER BY table_name;
EOF

echo -e "${GREEN}Test data generation completed!${NC}"
echo -e "${BLUE}This dataset includes:${NC}"
echo "  - Composite primary keys (multi-column PKs)"
echo "  - Composite foreign keys (multi-column FKs)"
echo "  - Self-referencing relationships (hierarchies)"
echo "  - Multi-level foreign key chains"
echo "  - Nullable foreign keys (orphaned records)"
echo "  - Circular dependencies (indirect)"
echo "  - Various data types and constraints"
echo ""
echo -e "${BLUE}Use this dataset to test your PK/FK discovery solution!${NC}"

