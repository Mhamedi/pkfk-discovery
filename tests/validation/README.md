# Validation Test Suite

This directory contains the validation test suite for adapter testing with dockerized Postgres and torture-test datasets.

## Overview

The validation test suite provides:
- A dockerized PostgreSQL instance for testing
- Complex test datasets with various PK/FK relationship patterns
- Multiple schemas (ecommerce, hr, inventory) with different relationship types
- Tables with explicit FK constraints
- Tables without explicit FK constraints (for discovery testing)

## Test Database Schemas

### E-commerce Schema
- `customers` - Customer master table
- `orders` - Orders with FK to customers
- `order_items` - Order line items with FKs to orders and products
- `products` - Product catalog
- `categories` - Product categories with self-referencing FK
- `reviews` - Product reviews (no explicit FKs, for discovery)
- `payments` - Payment records (no explicit FK to orders, for discovery)

### HR Schema
- `employees` - Employee table with self-referencing manager FK
- `departments` - Department master
- `locations` - Location master
- `salaries` - Employee salary history

### Inventory Schema
- `warehouses` - Warehouse master
- `items` - Item master
- `inventory` - Many-to-many relationship with composite PK

## Usage

### Start Test Database

```bash
./run_tests.sh
```

This will:
1. Start a PostgreSQL container on port 5433
2. Wait for the database to be ready
3. Create all test schemas and tables
4. Insert test data

### Stop Test Database

```bash
docker-compose -f docker-compose.test.yml down
```

### Connection Details

- **Host**: localhost
- **Port**: 5433
- **Database**: test_db
- **User**: test_user
- **Password**: test_password

### Running Adapter Validation Tests

Use these connection details when testing adapters:

```bash
# Example: Test connection
psql -h localhost -p 5433 -U test_user -d test_db

# Example: Run adapter validation
# (Use your adapter's test connection endpoint)
```

## Test Scenarios

The test suite includes:

1. **Explicit FK Constraints**: Tables with declared FOREIGN KEY constraints
2. **Implicit Relationships**: Tables without explicit FKs that should be discovered
3. **Self-Referencing**: Tables with FK to themselves (employees.manager_id)
4. **Composite Keys**: Tables with composite primary keys (inventory.inventory)
5. **Multiple Schemas**: Cross-schema relationships
6. **Nested Relationships**: Multi-level FK chains (orders -> customers, order_items -> orders/products)

## Expected Discovery Results

When running PK/FK discovery on this database, you should discover:

- `orders.customer_id` -> `customers.customer_id`
- `order_items.order_id` -> `orders.order_id`
- `order_items.product_id` -> `products.product_id`
- `products.category_id` -> `categories.category_id`
- `categories.parent_category_id` -> `categories.category_id` (self-reference)
- `employees.manager_id` -> `employees.employee_id` (self-reference)
- `employees.department_id` -> `departments.department_id`
- `departments.location_id` -> `locations.location_id`
- `salaries.employee_id` -> `employees.employee_id`
- `inventory.warehouse_id` -> `warehouses.warehouse_id`
- `inventory.item_id` -> `items.item_id`

Additionally, discovery should identify potential relationships:
- `reviews.product_id` -> `products.product_id` (no explicit FK)
- `reviews.customer_id` -> `customers.customer_id` (no explicit FK)
- `payments.order_id` -> `orders.order_id` (no explicit FK)

