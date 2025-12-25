-- Setup test database with torture-test datasets
-- This creates complex schemas with various PK/FK relationships

-- Create schemas
CREATE SCHEMA IF NOT EXISTS ecommerce;
CREATE SCHEMA IF NOT EXISTS hr;
CREATE SCHEMA IF NOT EXISTS inventory;

-- E-commerce schema with complex relationships
CREATE TABLE ecommerce.customers (
    customer_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ecommerce.orders (
    order_id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10,2),
    status VARCHAR(50),
    FOREIGN KEY (customer_id) REFERENCES ecommerce.customers(customer_id)
);

CREATE TABLE ecommerce.order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2),
    FOREIGN KEY (order_id) REFERENCES ecommerce.orders(order_id)
);

CREATE TABLE ecommerce.products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2),
    category_id INTEGER
);

CREATE TABLE ecommerce.categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_category_id INTEGER,
    FOREIGN KEY (parent_category_id) REFERENCES ecommerce.categories(category_id)
);

-- Add FK constraint after products table exists
ALTER TABLE ecommerce.order_items 
    ADD CONSTRAINT fk_order_items_product 
    FOREIGN KEY (product_id) REFERENCES ecommerce.products(product_id);

ALTER TABLE ecommerce.products 
    ADD CONSTRAINT fk_products_category 
    FOREIGN KEY (category_id) REFERENCES ecommerce.categories(category_id);

-- HR schema with self-referencing and complex relationships
CREATE TABLE hr.employees (
    employee_id SERIAL PRIMARY KEY,
    employee_number VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    manager_id INTEGER,
    department_id INTEGER,
    hire_date DATE,
    FOREIGN KEY (manager_id) REFERENCES hr.employees(employee_id)
);

CREATE TABLE hr.departments (
    department_id SERIAL PRIMARY KEY,
    department_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    location_id INTEGER
);

CREATE TABLE hr.locations (
    location_id SERIAL PRIMARY KEY,
    location_code VARCHAR(50) UNIQUE NOT NULL,
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100)
);

ALTER TABLE hr.employees 
    ADD CONSTRAINT fk_employees_department 
    FOREIGN KEY (department_id) REFERENCES hr.departments(department_id);

ALTER TABLE hr.departments 
    ADD CONSTRAINT fk_departments_location 
    FOREIGN KEY (location_id) REFERENCES hr.locations(location_id);

CREATE TABLE hr.salaries (
    salary_id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL,
    amount DECIMAL(10,2),
    effective_date DATE,
    FOREIGN KEY (employee_id) REFERENCES hr.employees(employee_id)
);

-- Inventory schema with composite keys and many-to-many
CREATE TABLE inventory.warehouses (
    warehouse_id SERIAL PRIMARY KEY,
    warehouse_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100),
    address TEXT
);

CREATE TABLE inventory.items (
    item_id SERIAL PRIMARY KEY,
    sku VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255),
    description TEXT
);

CREATE TABLE inventory.inventory (
    warehouse_id INTEGER NOT NULL,
    item_id INTEGER NOT NULL,
    quantity INTEGER DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (warehouse_id, item_id),
    FOREIGN KEY (warehouse_id) REFERENCES inventory.warehouses(warehouse_id),
    FOREIGN KEY (item_id) REFERENCES inventory.items(item_id)
);

-- Insert test data
INSERT INTO ecommerce.customers (email, first_name, last_name) VALUES
    ('john.doe@example.com', 'John', 'Doe'),
    ('jane.smith@example.com', 'Jane', 'Smith'),
    ('bob.jones@example.com', 'Bob', 'Jones');

INSERT INTO ecommerce.categories (name, parent_category_id) VALUES
    ('Electronics', NULL),
    ('Computers', 1),
    ('Phones', 1),
    ('Clothing', NULL);

INSERT INTO ecommerce.products (name, description, price, category_id) VALUES
    ('Laptop', 'High-performance laptop', 999.99, 2),
    ('Smartphone', 'Latest smartphone', 699.99, 3),
    ('T-Shirt', 'Cotton t-shirt', 19.99, 4);

INSERT INTO ecommerce.orders (customer_id, total_amount, status) VALUES
    (1, 999.99, 'completed'),
    (2, 699.99, 'pending'),
    (1, 19.99, 'completed');

INSERT INTO ecommerce.order_items (order_id, product_id, quantity, price) VALUES
    (1, 1, 1, 999.99),
    (2, 2, 1, 699.99),
    (3, 3, 1, 19.99);

INSERT INTO hr.locations (location_code, city, country) VALUES
    ('NYC', 'New York', 'USA'),
    ('LON', 'London', 'UK'),
    ('TOK', 'Tokyo', 'Japan');

INSERT INTO hr.departments (department_code, name, location_id) VALUES
    ('ENG', 'Engineering', 1),
    ('SALES', 'Sales', 2),
    ('HR', 'Human Resources', 1);

INSERT INTO hr.employees (employee_number, first_name, last_name, email, department_id, hire_date) VALUES
    ('EMP001', 'Alice', 'Johnson', 'alice@example.com', 1, '2020-01-15'),
    ('EMP002', 'Bob', 'Williams', 'bob@example.com', 1, '2021-03-20'),
    ('EMP003', 'Carol', 'Brown', 'carol@example.com', 2, '2019-06-10');

UPDATE hr.employees SET manager_id = 1 WHERE employee_id = 2;

INSERT INTO hr.salaries (employee_id, amount, effective_date) VALUES
    (1, 100000.00, '2024-01-01'),
    (2, 85000.00, '2024-01-01'),
    (3, 90000.00, '2024-01-01');

INSERT INTO inventory.warehouses (warehouse_code, name, address) VALUES
    ('WH001', 'Main Warehouse', '123 Main St'),
    ('WH002', 'Secondary Warehouse', '456 Oak Ave');

INSERT INTO inventory.items (sku, name, description) VALUES
    ('ITEM001', 'Widget A', 'Standard widget'),
    ('ITEM002', 'Widget B', 'Premium widget'),
    ('ITEM003', 'Gadget X', 'Electronic gadget');

INSERT INTO inventory.inventory (warehouse_id, item_id, quantity) VALUES
    (1, 1, 100),
    (1, 2, 50),
    (2, 1, 75),
    (2, 3, 25);

-- Create some tables without explicit FK constraints (for discovery testing)
CREATE TABLE ecommerce.reviews (
    review_id SERIAL PRIMARY KEY,
    product_id INTEGER,
    customer_id INTEGER,
    rating INTEGER,
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    -- Note: No FK constraints, should be discovered
);

CREATE TABLE ecommerce.payments (
    payment_id SERIAL PRIMARY KEY,
    order_id INTEGER,
    amount DECIMAL(10,2),
    payment_method VARCHAR(50),
    transaction_id VARCHAR(100)
    -- Note: No FK constraint to orders, should be discovered
);

