-- Complex Test Data Generator for PK/FK Discovery
-- This script creates a challenging dataset with ~250 records across multiple tables
-- Designed to test edge cases in PK/FK discovery:
-- - Composite keys (multi-column PKs and FKs)
-- - Self-referencing relationships
-- - Multi-level foreign key chains
-- - Nullable foreign keys
-- - Circular dependencies (indirect)
-- - Orphaned records
-- - Missing relationships

-- Drop existing tables if they exist (in reverse dependency order)
DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS product_variants CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS employees CASCADE;
DROP TABLE IF EXISTS departments CASCADE;
DROP TABLE IF EXISTS locations CASCADE;
DROP TABLE IF EXISTS suppliers CASCADE;
DROP TABLE IF EXISTS supplier_products CASCADE;
DROP TABLE IF EXISTS reviews CASCADE;
DROP TABLE IF EXISTS shipments CASCADE;
DROP TABLE IF EXISTS shipment_items CASCADE;
DROP TABLE IF EXISTS warehouses CASCADE;
DROP TABLE IF EXISTS inventory CASCADE;
DROP TABLE IF EXISTS price_history CASCADE;
DROP TABLE IF EXISTS customer_addresses CASCADE;
DROP TABLE IF EXISTS promotions CASCADE;
DROP TABLE IF EXISTS promotion_products CASCADE;

-- ============================================================================
-- BASE TABLES (No foreign keys)
-- ============================================================================

-- Locations table (self-referencing for regions)
CREATE TABLE locations (
    location_id SERIAL PRIMARY KEY,
    location_code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    parent_location_id INTEGER REFERENCES locations(location_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Departments table (self-referencing for hierarchy)
CREATE TABLE departments (
    dept_id SERIAL PRIMARY KEY,
    dept_code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    parent_dept_id INTEGER REFERENCES departments(dept_id) ON DELETE SET NULL,
    location_id INTEGER REFERENCES locations(location_id) ON DELETE SET NULL,
    budget DECIMAL(12,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Categories table (self-referencing for category hierarchy)
CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    category_code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    parent_category_id INTEGER REFERENCES categories(category_id) ON DELETE SET NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Warehouses table
CREATE TABLE warehouses (
    warehouse_id SERIAL PRIMARY KEY,
    warehouse_code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    location_id INTEGER REFERENCES locations(location_id) ON DELETE SET NULL,
    capacity INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Suppliers table
CREATE TABLE suppliers (
    supplier_id SERIAL PRIMARY KEY,
    supplier_code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    contact_email VARCHAR(255),
    location_id INTEGER REFERENCES locations(location_id) ON DELETE SET NULL,
    rating DECIMAL(3,2) CHECK (rating >= 0 AND rating <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- COMPOSITE KEY TABLES
-- ============================================================================

-- Products table (simple PK, but will have composite FKs)
CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    product_code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    category_id INTEGER REFERENCES categories(category_id) ON DELETE SET NULL,
    base_price DECIMAL(10,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Product Variants (composite PK: product_id + variant_code)
CREATE TABLE product_variants (
    product_id INTEGER NOT NULL,
    variant_code VARCHAR(20) NOT NULL,
    size VARCHAR(20),
    color VARCHAR(30),
    weight DECIMAL(8,2),
    additional_price DECIMAL(10,2) DEFAULT 0,
    PRIMARY KEY (product_id, variant_code),
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);

-- Customers table
CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    customer_code VARCHAR(10) NOT NULL UNIQUE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    registration_date DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Customer Addresses (composite key: customer_id + address_type)
CREATE TABLE customer_addresses (
    customer_id INTEGER NOT NULL,
    address_type VARCHAR(20) NOT NULL, -- 'billing', 'shipping', 'home'
    street_address VARCHAR(200) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(50),
    postal_code VARCHAR(20),
    country VARCHAR(50) DEFAULT 'USA',
    PRIMARY KEY (customer_id, address_type),
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE CASCADE
);

-- Employees table (self-referencing for manager hierarchy)
CREATE TABLE employees (
    employee_id SERIAL PRIMARY KEY,
    employee_code VARCHAR(10) NOT NULL UNIQUE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) UNIQUE,
    manager_id INTEGER REFERENCES employees(employee_id) ON DELETE SET NULL,
    dept_id INTEGER REFERENCES departments(dept_id) ON DELETE SET NULL,
    hire_date DATE NOT NULL,
    salary DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- COMPLEX RELATIONSHIP TABLES
-- ============================================================================

-- Orders table (references customers)
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    order_number VARCHAR(20) NOT NULL UNIQUE,
    customer_id INTEGER REFERENCES customers(customer_id) ON DELETE SET NULL,
    order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    total_amount DECIMAL(12,2),
    shipping_address_type VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Composite FK to customer_addresses
    FOREIGN KEY (customer_id, shipping_address_type) 
        REFERENCES customer_addresses(customer_id, address_type) 
        ON DELETE SET NULL
);

-- Order Items (composite FK to product_variants)
CREATE TABLE order_items (
    order_id INTEGER NOT NULL,
    line_item_number INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    variant_code VARCHAR(20),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) DEFAULT 0,
    PRIMARY KEY (order_id, line_item_number),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id, variant_code) 
        REFERENCES product_variants(product_id, variant_code) 
        ON DELETE RESTRICT
);

-- Supplier Products (many-to-many with composite key)
CREATE TABLE supplier_products (
    supplier_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    supplier_product_code VARCHAR(50) NOT NULL,
    cost DECIMAL(10,2) NOT NULL,
    min_order_quantity INTEGER DEFAULT 1,
    lead_time_days INTEGER,
    PRIMARY KEY (supplier_id, product_id, supplier_product_code),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(supplier_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);

-- Reviews (references products and customers)
CREATE TABLE reviews (
    review_id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
    customer_id INTEGER REFERENCES customers(customer_id) ON DELETE SET NULL,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    review_text TEXT,
    review_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    helpful_count INTEGER DEFAULT 0
);

-- Shipments table
CREATE TABLE shipments (
    shipment_id SERIAL PRIMARY KEY,
    shipment_number VARCHAR(20) NOT NULL UNIQUE,
    order_id INTEGER REFERENCES orders(order_id) ON DELETE SET NULL,
    warehouse_id INTEGER REFERENCES warehouses(warehouse_id) ON DELETE SET NULL,
    carrier VARCHAR(50),
    tracking_number VARCHAR(100),
    shipped_date TIMESTAMP,
    delivered_date TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending'
);

-- Shipment Items (composite FK)
CREATE TABLE shipment_items (
    shipment_id INTEGER NOT NULL,
    item_sequence INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    line_item_number INTEGER NOT NULL,
    quantity_shipped INTEGER NOT NULL,
    PRIMARY KEY (shipment_id, item_sequence),
    FOREIGN KEY (shipment_id) REFERENCES shipments(shipment_id) ON DELETE CASCADE,
    FOREIGN KEY (order_id, line_item_number) 
        REFERENCES order_items(order_id, line_item_number) 
        ON DELETE RESTRICT
);

-- Inventory (composite FK to product_variants and warehouses)
CREATE TABLE inventory (
    warehouse_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    variant_code VARCHAR(20) NOT NULL,
    quantity_on_hand INTEGER NOT NULL DEFAULT 0,
    reorder_point INTEGER DEFAULT 10,
    max_stock INTEGER,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (warehouse_id, product_id, variant_code),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(warehouse_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id, variant_code) 
        REFERENCES product_variants(product_id, variant_code) 
        ON DELETE CASCADE
);

-- Price History (tracks price changes over time)
CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
    variant_code VARCHAR(20),
    effective_date DATE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    changed_by_employee_id INTEGER REFERENCES employees(employee_id) ON DELETE SET NULL,
    FOREIGN KEY (product_id, variant_code) 
        REFERENCES product_variants(product_id, variant_code) 
        ON DELETE CASCADE
);

-- Promotions table
CREATE TABLE promotions (
    promotion_id SERIAL PRIMARY KEY,
    promotion_code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    discount_percentage DECIMAL(5,2),
    created_by_employee_id INTEGER REFERENCES employees(employee_id) ON DELETE SET NULL
);

-- Promotion Products (many-to-many)
CREATE TABLE promotion_products (
    promotion_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    variant_code VARCHAR(20),
    PRIMARY KEY (promotion_id, product_id, variant_code),
    FOREIGN KEY (promotion_id) REFERENCES promotions(promotion_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id, variant_code) 
        REFERENCES product_variants(product_id, variant_code) 
        ON DELETE CASCADE
);

-- ============================================================================
-- GENERATE TEST DATA (~250 records total)
-- ============================================================================

-- Insert Locations (15 records)
INSERT INTO locations (location_code, name, parent_location_id) VALUES
('LOC001', 'North America', NULL),
('LOC002', 'United States', 1),
('LOC003', 'Canada', 1),
('LOC004', 'Europe', NULL),
('LOC005', 'United Kingdom', 4),
('LOC006', 'Germany', 4),
('LOC007', 'France', 4),
('LOC008', 'California', 2),
('LOC009', 'New York', 2),
('LOC010', 'Texas', 2),
('LOC011', 'Ontario', 3),
('LOC012', 'London', 5),
('LOC013', 'Berlin', 6),
('LOC014', 'Paris', 7),
('LOC015', 'Asia Pacific', NULL);

-- Insert Warehouses (8 records)
INSERT INTO warehouses (warehouse_code, name, location_id, capacity) VALUES
('WH001', 'Main Warehouse', 8, 10000),
('WH002', 'East Coast Hub', 9, 8000),
('WH003', 'West Coast Distribution', 8, 12000),
('WH004', 'Texas Storage', 10, 6000),
('WH005', 'UK Warehouse', 12, 5000),
('WH006', 'European Hub', 13, 7000),
('WH007', 'France Distribution', 14, 4500),
('WH008', 'Canada Storage', 11, 5500);

-- Insert Departments (12 records with hierarchy)
INSERT INTO departments (dept_code, name, parent_dept_id, location_id, budget) VALUES
('DEPT001', 'Executive', NULL, 2, 5000000),
('DEPT002', 'Sales', NULL, 2, 2000000),
('DEPT003', 'Marketing', NULL, 2, 1500000),
('DEPT004', 'Operations', NULL, 2, 3000000),
('DEPT005', 'IT', NULL, 2, 2500000),
('DEPT006', 'Sales North', 2, 8, 800000),
('DEPT007', 'Sales South', 2, 10, 750000),
('DEPT008', 'Digital Marketing', 3, 2, 600000),
('DEPT009', 'Warehouse Operations', 4, 8, 1200000),
('DEPT010', 'Logistics', 4, 9, 900000),
('DEPT011', 'Software Development', 5, 2, 1800000),
('DEPT012', 'Infrastructure', 5, 2, 700000);

-- Insert Employees (20 records with manager hierarchy)
INSERT INTO employees (employee_code, first_name, last_name, email, manager_id, dept_id, hire_date, salary) VALUES
('EMP001', 'John', 'Smith', 'john.smith@company.com', NULL, 1, '2020-01-15', 250000),
('EMP002', 'Sarah', 'Johnson', 'sarah.johnson@company.com', 1, 2, '2020-03-20', 180000),
('EMP003', 'Michael', 'Brown', 'michael.brown@company.com', 1, 3, '2020-05-10', 170000),
('EMP004', 'Emily', 'Davis', 'emily.davis@company.com', 1, 4, '2020-02-01', 190000),
('EMP005', 'David', 'Wilson', 'david.wilson@company.com', 1, 5, '2020-04-05', 200000),
('EMP006', 'Lisa', 'Anderson', 'lisa.anderson@company.com', 2, 6, '2021-06-15', 120000),
('EMP007', 'Robert', 'Taylor', 'robert.taylor@company.com', 2, 7, '2021-07-20', 115000),
('EMP008', 'Jennifer', 'Martinez', 'jennifer.martinez@company.com', 3, 8, '2021-08-10', 110000),
('EMP009', 'William', 'Garcia', 'william.garcia@company.com', 4, 9, '2021-09-01', 130000),
('EMP010', 'Jessica', 'Rodriguez', 'jessica.rodriguez@company.com', 4, 10, '2021-10-15', 125000),
('EMP011', 'James', 'Lee', 'james.lee@company.com', 5, 11, '2021-11-20', 140000),
('EMP012', 'Amanda', 'White', 'amanda.white@company.com', 5, 12, '2022-01-10', 135000),
('EMP013', 'Christopher', 'Harris', 'christopher.harris@company.com', 6, 6, '2022-03-15', 95000),
('EMP014', 'Melissa', 'Clark', 'melissa.clark@company.com', 6, 6, '2022-04-20', 92000),
('EMP015', 'Daniel', 'Lewis', 'daniel.lewis@company.com', 7, 7, '2022-05-10', 90000),
('EMP016', 'Michelle', 'Walker', 'michelle.walker@company.com', 8, 8, '2022-06-15', 88000),
('EMP017', 'Matthew', 'Hall', 'matthew.hall@company.com', 9, 9, '2022-07-20', 105000),
('EMP018', 'Nicole', 'Allen', 'nicole.allen@company.com', 10, 10, '2022-08-10', 100000),
('EMP019', 'Andrew', 'Young', 'andrew.young@company.com', 11, 11, '2022-09-15', 115000),
('EMP020', 'Stephanie', 'King', 'stephanie.king@company.com', 12, 12, '2022-10-20', 110000);

-- Insert Suppliers (10 records)
INSERT INTO suppliers (supplier_code, name, contact_email, location_id, rating) VALUES
('SUP001', 'Global Supplies Inc', 'contact@globalsupplies.com', 8, 4.5),
('SUP002', 'Quality Materials Co', 'info@qualitymaterials.com', 9, 4.2),
('SUP003', 'Premium Products Ltd', 'sales@premiumproducts.com', 10, 4.8),
('SUP004', 'Reliable Source Corp', 'contact@reliablesource.com', 8, 4.0),
('SUP005', 'Best Value Suppliers', 'info@bestvalue.com', 9, 3.8),
('SUP006', 'European Imports', 'sales@europeanimports.com', 12, 4.6),
('SUP007', 'Asian Manufacturing', 'contact@asianmfg.com', 15, 4.3),
('SUP008', 'North American Goods', 'info@nagoods.com', 11, 4.1),
('SUP009', 'Tech Components Inc', 'sales@techcomponents.com', 8, 4.7),
('SUP010', 'Industrial Supplies', 'contact@industrialsupplies.com', 10, 3.9);

-- Insert Categories (15 records with hierarchy)
INSERT INTO categories (category_code, name, parent_category_id, description) VALUES
('CAT001', 'Electronics', NULL, 'Electronic devices and components'),
('CAT002', 'Computers', 1, 'Desktop and laptop computers'),
('CAT003', 'Mobile Devices', 1, 'Smartphones and tablets'),
('CAT004', 'Accessories', 1, 'Electronic accessories'),
('CAT005', 'Clothing', NULL, 'Apparel and fashion items'),
('CAT006', 'Men''s Clothing', 5, 'Clothing for men'),
('CAT007', 'Women''s Clothing', 5, 'Clothing for women'),
('CAT008', 'Footwear', 5, 'Shoes and boots'),
('CAT009', 'Home & Garden', NULL, 'Home improvement and garden items'),
('CAT010', 'Furniture', 9, 'Home and office furniture'),
('CAT011', 'Kitchenware', 9, 'Kitchen tools and appliances'),
('CAT012', 'Sports & Outdoors', NULL, 'Sports equipment and outdoor gear'),
('CAT013', 'Books', NULL, 'Books and publications'),
('CAT014', 'Toys & Games', NULL, 'Toys and board games'),
('CAT015', 'Automotive', NULL, 'Car parts and accessories');

-- Insert Products (25 records)
INSERT INTO products (product_code, name, category_id, base_price, description) VALUES
('PROD001', 'Laptop Pro 15"', 2, 1299.99, 'High-performance laptop'),
('PROD002', 'Laptop Air 13"', 2, 999.99, 'Lightweight laptop'),
('PROD003', 'Smartphone X', 3, 799.99, 'Latest smartphone model'),
('PROD004', 'Tablet Pro', 3, 599.99, 'Professional tablet'),
('PROD005', 'Wireless Mouse', 4, 29.99, 'Ergonomic wireless mouse'),
('PROD006', 'Keyboard Mechanical', 4, 89.99, 'Mechanical keyboard'),
('PROD007', 'Men''s T-Shirt', 6, 19.99, 'Cotton t-shirt'),
('PROD008', 'Men''s Jeans', 6, 49.99, 'Classic fit jeans'),
('PROD009', 'Women''s Dress', 7, 59.99, 'Elegant dress'),
('PROD010', 'Women''s Blouse', 7, 39.99, 'Professional blouse'),
('PROD011', 'Running Shoes', 8, 89.99, 'Athletic running shoes'),
('PROD012', 'Dress Shoes', 8, 129.99, 'Formal dress shoes'),
('PROD013', 'Office Chair', 10, 199.99, 'Ergonomic office chair'),
('PROD014', 'Dining Table', 10, 499.99, 'Wooden dining table'),
('PROD015', 'Coffee Maker', 11, 79.99, 'Programmable coffee maker'),
('PROD016', 'Basketball', 12, 24.99, 'Official size basketball'),
('PROD017', 'Tennis Racket', 12, 89.99, 'Professional tennis racket'),
('PROD018', 'Mystery Novel', 13, 14.99, 'Bestselling mystery novel'),
('PROD019', 'Board Game', 14, 34.99, 'Family board game'),
('PROD020', 'Car Battery', 15, 129.99, 'Automotive battery'),
('PROD021', 'Tire Set', 15, 399.99, 'Set of 4 tires'),
('PROD022', 'USB-C Cable', 4, 12.99, 'Fast charging cable'),
('PROD023', 'Monitor 27"', 2, 299.99, '4K monitor'),
('PROD024', 'Headphones', 4, 149.99, 'Noise-cancelling headphones'),
('PROD025', 'Backpack', 4, 49.99, 'Laptop backpack');

-- Insert Product Variants (50 records - composite keys)
INSERT INTO product_variants (product_id, variant_code, size, color, weight, additional_price) VALUES
-- Laptop variants
(1, 'LAP15-BLK-512', '15"', 'Black', 2.1, 0),
(1, 'LAP15-SLV-512', '15"', 'Silver', 2.1, 0),
(1, 'LAP15-BLK-1TB', '15"', 'Black', 2.1, 200),
(2, 'LAP13-BLK-256', '13"', 'Black', 1.3, 0),
(2, 'LAP13-GLD-256', '13"', 'Gold', 1.3, 100),
-- Smartphone variants
(3, 'PHX-BLK-128', NULL, 'Black', 0.2, 0),
(3, 'PHX-WHT-128', NULL, 'White', 0.2, 0),
(3, 'PHX-BLK-256', NULL, 'Black', 0.2, 100),
(4, 'TAB-GRY-64', '10.5"', 'Gray', 0.5, 0),
(4, 'TAB-BLK-128', '10.5"', 'Black', 0.5, 50),
-- Clothing variants
(7, 'TSH-S-M-BLK', 'S', 'Black', NULL, 0),
(7, 'TSH-M-M-BLK', 'M', 'Black', NULL, 0),
(7, 'TSH-L-M-BLK', 'L', 'Black', NULL, 0),
(7, 'TSH-XL-M-BLK', 'XL', 'Black', NULL, 0),
(8, 'JNS-32-BLU', '32', 'Blue', NULL, 0),
(8, 'JNS-34-BLU', '34', 'Blue', NULL, 0),
(8, 'JNS-36-BLU', '36', 'Blue', NULL, 0),
(9, 'DRS-S-RED', 'S', 'Red', NULL, 0),
(9, 'DRS-M-RED', 'M', 'Red', NULL, 0),
(9, 'DRS-L-RED', 'L', 'Red', NULL, 0),
(10, 'BLS-S-WHT', 'S', 'White', NULL, 0),
(10, 'BLS-M-WHT', 'M', 'White', NULL, 0),
-- Shoe variants
(11, 'RUN-8-BLK', '8', 'Black', 0.8, 0),
(11, 'RUN-9-BLK', '9', 'Black', 0.8, 0),
(11, 'RUN-10-BLK', '10', 'Black', 0.8, 0),
(11, 'RUN-8-WHT', '8', 'White', 0.8, 0),
(12, 'DSH-9-BRN', '9', 'Brown', 1.2, 0),
(12, 'DSH-10-BRN', '10', 'Brown', 1.2, 0),
(12, 'DSH-9-BLK', '9', 'Black', 1.2, 0),
-- Furniture variants
(13, 'CHR-BLK', NULL, 'Black', 15.0, 0),
(13, 'CHR-BRN', NULL, 'Brown', 15.0, 0),
(13, 'CHR-GRY', NULL, 'Gray', 15.0, 0),
(14, 'TBL-6-SEA', '6 seats', 'Natural', 50.0, 0),
(14, 'TBL-8-SEA', '8 seats', 'Natural', 60.0, 100),
-- Electronics variants
(15, 'CMK-BLK', NULL, 'Black', 3.5, 0),
(15, 'CMK-WHT', NULL, 'White', 3.5, 0),
(16, 'BBL-STD', 'Standard', 'Orange', 0.6, 0),
(17, 'TRK-ADULT', 'Adult', NULL, 0.3, 0),
(17, 'TRK-JUNIOR', 'Junior', NULL, 0.25, -10),
(18, 'BOOK-HC', 'Hardcover', NULL, 0.5, 5),
(18, 'BOOK-PB', 'Paperback', NULL, 0.3, 0),
(19, 'GAME-STD', 'Standard', NULL, 1.5, 0),
(20, 'BAT-STD', 'Standard', NULL, 15.0, 0),
(20, 'BAT-PREM', 'Premium', NULL, 15.0, 30),
(21, 'TIRE-16', '16"', 'Black', 20.0, 0),
(21, 'TIRE-17', '17"', 'Black', 22.0, 50),
(5, 'MOUSE-WRLS-BLK', NULL, 'Black', 0.1, 0),
(5, 'MOUSE-WRLS-WHT', NULL, 'White', 0.1, 0),
(6, 'KEY-MECH-BLK', NULL, 'Black', 0.8, 0),
(6, 'KEY-MECH-WHT', NULL, 'White', 0.8, 0),
(22, 'CBL-6FT', '6ft', 'Black', 0.1, 0),
(22, 'CBL-10FT', '10ft', 'Black', 0.15, 5),
(23, 'MON-27-BLK', '27"', 'Black', 8.0, 0),
(24, 'HD-BLK', NULL, 'Black', 0.3, 0),
(24, 'HD-WHT', NULL, 'White', 0.3, 0),
(25, 'BPK-BLK', NULL, 'Black', 1.2, 0),
(25, 'BPK-BLU', NULL, 'Blue', 1.2, 0);

-- Insert Customers (20 records)
INSERT INTO customers (customer_code, first_name, last_name, email, phone, registration_date) VALUES
('CUST001', 'Alice', 'Johnson', 'alice.johnson@email.com', '555-0101', '2023-01-15'),
('CUST002', 'Bob', 'Smith', 'bob.smith@email.com', '555-0102', '2023-02-20'),
('CUST003', 'Carol', 'Williams', 'carol.williams@email.com', '555-0103', '2023-03-10'),
('CUST004', 'David', 'Brown', 'david.brown@email.com', '555-0104', '2023-04-05'),
('CUST005', 'Emma', 'Jones', 'emma.jones@email.com', '555-0105', '2023-05-12'),
('CUST006', 'Frank', 'Garcia', 'frank.garcia@email.com', '555-0106', '2023-06-18'),
('CUST007', 'Grace', 'Miller', 'grace.miller@email.com', '555-0107', '2023-07-22'),
('CUST008', 'Henry', 'Davis', 'henry.davis@email.com', '555-0108', '2023-08-30'),
('CUST009', 'Ivy', 'Rodriguez', 'ivy.rodriguez@email.com', '555-0109', '2023-09-15'),
('CUST010', 'Jack', 'Martinez', 'jack.martinez@email.com', '555-0110', '2023-10-20'),
('CUST011', 'Karen', 'Hernandez', 'karen.hernandez@email.com', '555-0111', '2023-11-05'),
('CUST012', 'Larry', 'Lopez', 'larry.lopez@email.com', '555-0112', '2023-12-10'),
('CUST013', 'Mary', 'Wilson', 'mary.wilson@email.com', '555-0113', '2024-01-15'),
('CUST014', 'Nancy', 'Anderson', 'nancy.anderson@email.com', '555-0114', '2024-02-20'),
('CUST015', 'Oscar', 'Thomas', 'oscar.thomas@email.com', '555-0115', '2024-03-10'),
('CUST016', 'Patricia', 'Taylor', 'patricia.taylor@email.com', '555-0116', '2024-04-05'),
('CUST017', 'Quinn', 'Moore', 'quinn.moore@email.com', '555-0117', '2024-05-12'),
('CUST018', 'Rachel', 'Jackson', 'rachel.jackson@email.com', '555-0118', '2024-06-18'),
('CUST019', 'Steve', 'White', 'steve.white@email.com', '555-0119', '2024-07-22'),
('CUST020', 'Tina', 'Harris', 'tina.harris@email.com', '555-0120', '2024-08-30');

-- Insert Customer Addresses (40 records - composite keys)
INSERT INTO customer_addresses (customer_id, address_type, street_address, city, state, postal_code, country) VALUES
(1, 'billing', '123 Main St', 'Los Angeles', 'CA', '90001', 'USA'),
(1, 'shipping', '123 Main St', 'Los Angeles', 'CA', '90001', 'USA'),
(2, 'billing', '456 Oak Ave', 'New York', 'NY', '10001', 'USA'),
(2, 'shipping', '789 Pine Rd', 'New York', 'NY', '10002', 'USA'),
(3, 'billing', '321 Elm St', 'Chicago', 'IL', '60601', 'USA'),
(3, 'shipping', '321 Elm St', 'Chicago', 'IL', '60601', 'USA'),
(4, 'billing', '654 Maple Dr', 'Houston', 'TX', '77001', 'USA'),
(4, 'shipping', '654 Maple Dr', 'Houston', 'TX', '77001', 'USA'),
(5, 'billing', '987 Cedar Ln', 'Phoenix', 'AZ', '85001', 'USA'),
(5, 'shipping', '987 Cedar Ln', 'Phoenix', 'AZ', '85001', 'USA'),
(6, 'billing', '147 Birch Way', 'Philadelphia', 'PA', '19101', 'USA'),
(6, 'shipping', '258 Spruce Ct', 'Philadelphia', 'PA', '19102', 'USA'),
(7, 'billing', '369 Willow St', 'San Antonio', 'TX', '78201', 'USA'),
(7, 'shipping', '369 Willow St', 'San Antonio', 'TX', '78201', 'USA'),
(8, 'billing', '741 Ash Ave', 'San Diego', 'CA', '92101', 'USA'),
(8, 'shipping', '852 Poplar Rd', 'San Diego', 'CA', '92102', 'USA'),
(9, 'billing', '963 Cherry Blvd', 'Dallas', 'TX', '75201', 'USA'),
(9, 'shipping', '963 Cherry Blvd', 'Dallas', 'TX', '75201', 'USA'),
(10, 'billing', '159 Walnut Dr', 'San Jose', 'CA', '95101', 'USA'),
(10, 'shipping', '159 Walnut Dr', 'San Jose', 'CA', '95101', 'USA'),
(11, 'billing', '357 Chestnut Ln', 'Austin', 'TX', '78701', 'USA'),
(11, 'shipping', '357 Chestnut Ln', 'Austin', 'TX', '78701', 'USA'),
(12, 'billing', '741 Hickory Way', 'Jacksonville', 'FL', '32201', 'USA'),
(12, 'shipping', '852 Sycamore Ct', 'Jacksonville', 'FL', '32202', 'USA'),
(13, 'billing', '963 Magnolia St', 'Fort Worth', 'TX', '76101', 'USA'),
(13, 'shipping', '963 Magnolia St', 'Fort Worth', 'TX', '76101', 'USA'),
(14, 'billing', '159 Dogwood Ave', 'Columbus', 'OH', '43201', 'USA'),
(14, 'shipping', '258 Redwood Rd', 'Columbus', 'OH', '43202', 'USA'),
(15, 'billing', '369 Cypress Blvd', 'Charlotte', 'NC', '28201', 'USA'),
(15, 'shipping', '369 Cypress Blvd', 'Charlotte', 'NC', '28201', 'USA'),
(16, 'billing', '741 Fir Dr', 'San Francisco', 'CA', '94101', 'USA'),
(16, 'shipping', '852 Hemlock Ln', 'San Francisco', 'CA', '94102', 'USA'),
(17, 'billing', '963 Juniper Way', 'Indianapolis', 'IN', '46201', 'USA'),
(17, 'shipping', '963 Juniper Way', 'Indianapolis', 'IN', '46201', 'USA'),
(18, 'billing', '159 Larch Ct', 'Seattle', 'WA', '98101', 'USA'),
(18, 'shipping', '258 Mahogany St', 'Seattle', 'WA', '98102', 'USA'),
(19, 'billing', '369 Teak Ave', 'Denver', 'CO', '80201', 'USA'),
(19, 'shipping', '369 Teak Ave', 'Denver', 'CO', '80201', 'USA'),
(20, 'billing', '741 Rosewood Rd', 'Washington', 'DC', '20001', 'USA'),
(20, 'shipping', '852 Ebony Blvd', 'Washington', 'DC', '20002', 'USA');

-- Insert Orders (30 records)
INSERT INTO orders (order_number, customer_id, order_date, status, total_amount, shipping_address_type) VALUES
('ORD001', 1, '2024-01-20 10:30:00', 'completed', 1299.99, 'shipping'),
('ORD002', 2, '2024-02-15 14:20:00', 'completed', 1049.98, 'shipping'),
('ORD003', 3, '2024-03-10 09:15:00', 'pending', 599.99, 'shipping'),
('ORD004', 4, '2024-04-05 16:45:00', 'completed', 129.98, 'shipping'),
('ORD005', 5, '2024-05-12 11:30:00', 'shipped', 89.99, 'shipping'),
('ORD006', 6, '2024-06-18 13:20:00', 'completed', 199.99, 'shipping'),
('ORD007', 7, '2024-07-22 15:10:00', 'completed', 49.99, 'shipping'),
('ORD008', 8, '2024-08-30 10:00:00', 'pending', 59.99, 'shipping'),
('ORD009', 9, '2024-09-15 14:30:00', 'completed', 179.98, 'shipping'),
('ORD010', 10, '2024-10-20 09:45:00', 'shipped', 299.99, 'shipping'),
('ORD011', 11, '2024-11-05 12:15:00', 'completed', 149.99, 'shipping'),
('ORD012', 12, '2024-12-10 16:20:00', 'completed', 34.99, 'shipping'),
('ORD013', 13, '2024-01-25 10:30:00', 'pending', 399.99, 'shipping'),
('ORD014', 14, '2024-02-28 14:20:00', 'completed', 79.99, 'shipping'),
('ORD015', 15, '2024-03-15 09:15:00', 'completed', 24.99, 'shipping'),
('ORD016', 16, '2024-04-20 16:45:00', 'shipped', 89.99, 'shipping'),
('ORD017', 17, '2024-05-25 11:30:00', 'completed', 14.99, 'shipping'),
('ORD018', 18, '2024-06-30 13:20:00', 'pending', 49.99, 'shipping'),
('ORD019', 19, '2024-07-05 15:10:00', 'completed', 129.99, 'shipping'),
('ORD020', 20, '2024-08-10 10:00:00', 'completed', 12.99, 'shipping'),
('ORD021', 1, '2024-09-20 14:30:00', 'completed', 29.99, 'shipping'),
('ORD022', 2, '2024-10-25 09:45:00', 'shipped', 199.99, 'shipping'),
('ORD023', 3, '2024-11-30 12:15:00', 'completed', 89.99, 'shipping'),
('ORD024', 4, '2024-12-05 16:20:00', 'pending', 149.99, 'shipping'),
('ORD025', 5, '2024-01-10 10:30:00', 'completed', 499.99, 'shipping'),
('ORD026', 6, '2024-02-15 14:20:00', 'completed', 79.99, 'shipping'),
('ORD027', 7, '2024-03-20 09:15:00', 'shipped', 24.99, 'shipping'),
('ORD028', 8, '2024-04-25 16:45:00', 'completed', 89.99, 'shipping'),
('ORD029', 9, '2024-05-30 11:30:00', 'pending', 34.99, 'shipping'),
('ORD030', 10, '2024-06-05 13:20:00', 'completed', 129.99, 'shipping');

-- Insert Order Items (60 records - composite FKs)
INSERT INTO order_items (order_id, line_item_number, product_id, variant_code, quantity, unit_price, discount) VALUES
(1, 1, 1, 'LAP15-BLK-512', 1, 1299.99, 0),
(2, 1, 2, 'LAP13-BLK-256', 1, 999.99, 0),
(2, 2, 22, 'CBL-6FT', 1, 12.99, 0),
(2, 3, 24, 'HD-BLK', 1, 149.99, 0),
(3, 1, 4, 'TAB-GRY-64', 1, 599.99, 0),
(4, 1, 12, 'DSH-9-BRN', 1, 129.99, 0),
(5, 1, 6, 'KEY-MECH-BLK', 1, 89.99, 0),
(6, 1, 13, 'CHR-BLK', 1, 199.99, 0),
(7, 1, 8, 'JNS-32-BLU', 1, 49.99, 0),
(8, 1, 9, 'DRS-S-RED', 1, 59.99, 0),
(9, 1, 11, 'RUN-9-BLK', 1, 89.99, 0),
(9, 2, 11, 'RUN-10-BLK', 1, 89.99, 0),
(10, 1, 23, 'MON-27-BLK', 1, 299.99, 0),
(11, 1, 24, 'HD-BLK', 1, 149.99, 0),
(12, 1, 19, 'GAME-STD', 1, 34.99, 0),
(13, 1, 21, 'TIRE-16', 1, 399.99, 0),
(14, 1, 15, 'CMK-BLK', 1, 79.99, 0),
(15, 1, 16, 'BBL-STD', 1, 24.99, 0),
(16, 1, 17, 'TRK-ADULT', 1, 89.99, 0),
(17, 1, 18, 'BOOK-PB', 1, 14.99, 0),
(18, 1, 25, 'BPK-BLK', 1, 49.99, 0),
(19, 1, 20, 'BAT-STD', 1, 129.99, 0),
(20, 1, 22, 'CBL-6FT', 1, 12.99, 0),
(21, 1, 5, 'MOUSE-WRLS-BLK', 1, 29.99, 0),
(22, 1, 13, 'CHR-BRN', 1, 199.99, 0),
(23, 1, 6, 'KEY-MECH-BLK', 1, 89.99, 0),
(24, 1, 24, 'HD-WHT', 1, 149.99, 0),
(25, 1, 14, 'TBL-6-SEA', 1, 499.99, 0),
(26, 1, 15, 'CMK-WHT', 1, 79.99, 0),
(27, 1, 16, 'BBL-STD', 1, 24.99, 0),
(28, 1, 17, 'TRK-JUNIOR', 1, 79.99, 10),
(29, 1, 19, 'GAME-STD', 1, 34.99, 0),
(30, 1, 20, 'BAT-PREM', 1, 159.99, 30),
(1, 2, 22, 'CBL-10FT', 1, 17.99, 0),
(3, 2, 22, 'CBL-6FT', 1, 12.99, 0),
(4, 2, 7, 'TSH-M-M-BLK', 2, 39.98, 0),
(5, 2, 5, 'MOUSE-WRLS-BLK', 2, 59.98, 0),
(6, 2, 25, 'BPK-BLU', 1, 49.99, 0),
(7, 2, 7, 'TSH-L-M-BLK', 1, 19.99, 0),
(8, 2, 10, 'BLS-S-WHT', 1, 39.99, 0),
(10, 2, 22, 'CBL-6FT', 1, 12.99, 0),
(11, 2, 22, 'CBL-10FT', 1, 17.99, 0),
(12, 2, 18, 'BOOK-HC', 1, 19.99, 5),
(13, 2, 20, 'BAT-STD', 1, 129.99, 0),
(14, 2, 22, 'CBL-6FT', 1, 12.99, 0),
(15, 2, 17, 'TRK-ADULT', 1, 89.99, 0),
(16, 2, 22, 'CBL-6FT', 1, 12.99, 0),
(17, 2, 18, 'BOOK-PB', 1, 14.99, 0),
(18, 2, 25, 'BPK-BLK', 1, 49.99, 0),
(19, 2, 22, 'CBL-10FT', 1, 17.99, 0),
(20, 2, 22, 'CBL-6FT', 2, 25.98, 0),
(21, 2, 24, 'HD-BLK', 1, 149.99, 0),
(22, 2, 25, 'BPK-BLU', 1, 49.99, 0),
(23, 2, 22, 'CBL-10FT', 1, 17.99, 0),
(24, 2, 22, 'CBL-6FT', 1, 12.99, 0),
(25, 2, 13, 'CHR-GRY', 4, 799.96, 0),
(26, 2, 22, 'CBL-6FT', 1, 12.99, 0),
(27, 2, 17, 'TRK-JUNIOR', 1, 79.99, 10),
(28, 2, 22, 'CBL-6FT', 1, 12.99, 0),
(29, 2, 18, 'BOOK-PB', 2, 29.98, 0),
(30, 2, 22, 'CBL-10FT', 1, 17.99, 0);

-- Insert Supplier Products (25 records - composite keys)
INSERT INTO supplier_products (supplier_id, product_id, supplier_product_code, cost, min_order_quantity, lead_time_days) VALUES
(1, 1, 'SUP001-PROD001-LAP15', 900.00, 10, 14),
(2, 2, 'SUP002-PROD002-LAP13', 700.00, 10, 14),
(3, 3, 'SUP003-PROD003-PHX', 500.00, 20, 21),
(4, 4, 'SUP004-PROD004-TAB', 400.00, 15, 14),
(5, 5, 'SUP005-PROD005-MOUSE-WRLS', 15.00, 100, 7),
(6, 6, 'SUP006-PROD006-KEY', 50.00, 50, 10),
(7, 7, 'SUP007-PROD007-TSH', 8.00, 200, 14),
(8, 8, 'SUP008-PROD008-JNS', 25.00, 100, 21),
(9, 9, 'SUP009-PROD009-DRS', 30.00, 50, 14),
(10, 10, 'SUP010-PROD010-BLS', 18.00, 100, 14),
(1, 11, 'SUP001-PROD011-RUN', 50.00, 50, 21),
(2, 12, 'SUP002-PROD012-DSH', 70.00, 30, 21),
(3, 13, 'SUP003-PROD013-CHR', 120.00, 20, 14),
(4, 14, 'SUP004-PROD014-TBL', 300.00, 5, 30),
(5, 15, 'SUP005-PROD015-CMK', 45.00, 50, 10),
(6, 16, 'SUP006-PROD016-BBL', 12.00, 100, 7),
(7, 17, 'SUP007-PROD017-TRK', 50.00, 30, 14),
(8, 18, 'SUP008-PROD018-BOOK', 8.00, 500, 7),
(9, 19, 'SUP009-PROD019-GAME', 18.00, 100, 14),
(10, 20, 'SUP010-PROD020-BAT', 80.00, 20, 14),
(1, 21, 'SUP001-PROD021-TIRE', 250.00, 10, 30),
(2, 22, 'SUP002-PROD022-CBL', 6.00, 200, 7),
(3, 23, 'SUP003-PROD023-MON', 180.00, 15, 14),
(4, 24, 'SUP004-PROD024-HD', 90.00, 30, 14),
(5, 25, 'SUP005-PROD025-BPK', 25.00, 50, 10);

-- Insert Reviews (20 records)
INSERT INTO reviews (product_id, customer_id, rating, review_text, helpful_count) VALUES
(1, 1, 5, 'Excellent laptop, very fast!', 12),
(2, 2, 4, 'Good value for money', 8),
(3, 3, 5, 'Love this phone!', 15),
(4, 4, 4, 'Great tablet for work', 6),
(5, 5, 5, 'MOUSE-WRLS-BLK', 20),
(6, 6, 4, 'Nice keyboard, good build quality', 10),
(7, 7, 3, 'T-shirt is okay, but runs small', 4),
(8, 8, 5, 'Perfect fit, great quality', 9),
(9, 9, 4, 'Beautiful dress, runs true to size', 7),
(10, 10, 5, 'Love this blouse!', 11),
(11, 11, 5, 'Comfortable running shoes', 18),
(12, 12, 4, 'Good dress shoes', 5),
(13, 13, 5, 'Very comfortable chair', 14),
(14, 14, 4, 'Nice table, assembly required', 6),
(15, 15, 5, 'Great coffee maker', 13),
(16, 16, 4, 'Good basketball', 3),
(17, 17, 5, 'Excellent racket', 8),
(18, 18, 5, 'Great mystery novel', 22),
(19, 19, 4, 'Fun family game', 15),
(20, 20, 3, 'Battery works but short lifespan', 2);

-- Insert Shipments (15 records)
INSERT INTO shipments (shipment_number, order_id, warehouse_id, carrier, tracking_number, shipped_date, delivered_date, status) VALUES
('SHIP001', 1, 1, 'FedEx', 'FX123456789', '2024-01-22 10:00:00', '2024-01-24 14:30:00', 'delivered'),
('SHIP002', 2, 1, 'UPS', 'UPS987654321', '2024-02-17 09:00:00', '2024-02-19 16:00:00', 'delivered'),
('SHIP003', 5, 2, 'USPS', 'US123456789', '2024-05-14 11:00:00', NULL, 'in_transit'),
('SHIP004', 6, 1, 'FedEx', 'FX234567890', '2024-06-20 10:00:00', '2024-06-22 12:00:00', 'delivered'),
('SHIP005', 10, 3, 'UPS', 'UPS876543210', '2024-10-22 09:00:00', NULL, 'in_transit'),
('SHIP006', 11, 1, 'FedEx', 'FX345678901', '2024-11-07 10:00:00', '2024-11-09 15:00:00', 'delivered'),
('SHIP007', 15, 2, 'USPS', 'US234567890', '2024-03-17 11:00:00', '2024-03-19 10:00:00', 'delivered'),
('SHIP008', 16, 1, 'FedEx', 'FX456789012', '2024-04-27 10:00:00', NULL, 'in_transit'),
('SHIP009', 19, 3, 'UPS', 'UPS765432109', '2024-07-07 09:00:00', '2024-07-09 14:00:00', 'delivered'),
('SHIP010', 22, 1, 'FedEx', 'FX567890123', '2024-10-27 10:00:00', NULL, 'in_transit'),
('SHIP011', 27, 2, 'USPS', 'US345678901', '2024-03-22 11:00:00', '2024-03-24 13:00:00', 'delivered'),
('SHIP012', 28, 1, 'FedEx', 'FX678901234', '2024-04-27 10:00:00', NULL, 'in_transit'),
('SHIP013', 30, 3, 'UPS', 'UPS654321098', '2024-06-07 09:00:00', '2024-06-09 16:00:00', 'delivered'),
('SHIP014', 4, 1, 'FedEx', 'FX789012345', '2024-04-07 10:00:00', '2024-04-09 11:00:00', 'delivered'),
('SHIP015', 7, 2, 'USPS', 'US456789012', '2024-07-24 11:00:00', '2024-07-26 10:00:00', 'delivered');

-- Insert Shipment Items (25 records - composite FKs)
INSERT INTO shipment_items (shipment_id, item_sequence, order_id, line_item_number, quantity_shipped) VALUES
(1, 1, 1, 1, 1),
(1, 2, 1, 2, 1),
(2, 1, 2, 1, 1),
(2, 2, 2, 2, 1),
(2, 3, 2, 3, 1),
(3, 1, 5, 1, 1),
(3, 2, 5, 2, 2),
(4, 1, 6, 1, 1),
(4, 2, 6, 2, 1),
(5, 1, 10, 1, 1),
(5, 2, 10, 2, 1),
(6, 1, 11, 1, 1),
(6, 2, 11, 2, 1),
(7, 1, 15, 1, 1),
(7, 2, 15, 2, 1),
(8, 1, 16, 1, 1),
(8, 2, 16, 2, 1),
(9, 1, 19, 1, 1),
(9, 2, 19, 2, 1),
(10, 1, 22, 1, 1),
(10, 2, 22, 2, 1),
(11, 1, 27, 1, 1),
(11, 2, 27, 2, 1),
(12, 1, 28, 1, 1),
(12, 2, 28, 2, 1),
(13, 1, 30, 1, 1),
(13, 2, 30, 2, 1);

-- Insert Inventory (30 records - composite keys)
INSERT INTO inventory (warehouse_id, product_id, variant_code, quantity_on_hand, reorder_point, max_stock) VALUES
(1, 1, 'LAP15-BLK-512', 25, 10, 100),
(1, 1, 'LAP15-SLV-512', 20, 10, 100),
(1, 2, 'LAP13-BLK-256', 30, 15, 150),
(1, 3, 'PHX-BLK-128', 50, 20, 200),
(1, 5, 'MOUSE-WRLS-BLK', 200, 50, 500),
(1, 5, 'MOUSE-WRLS-WHT', 150, 50, 500),
(1, 6, 'KEY-MECH-BLK', 75, 25, 200),
(2, 4, 'TAB-GRY-64', 40, 15, 150),
(2, 7, 'TSH-S-M-BLK', 100, 30, 300),
(2, 8, 'JNS-32-BLU', 60, 20, 150),
(2, 11, 'RUN-9-BLK', 80, 25, 200),
(2, 13, 'CHR-BLK', 35, 10, 100),
(3, 9, 'DRS-S-RED', 45, 15, 150),
(3, 10, 'BLS-S-WHT', 70, 25, 200),
(3, 12, 'DSH-9-BRN', 50, 15, 150),
(3, 15, 'CMK-BLK', 90, 30, 250),
(3, 16, 'BBL-STD', 150, 50, 500),
(4, 17, 'TRK-ADULT', 55, 20, 150),
(4, 18, 'BOOK-PB', 200, 50, 1000),
(4, 19, 'GAME-STD', 120, 40, 300),
(4, 20, 'BAT-STD', 40, 15, 150),
(4, 21, 'TIRE-16', 25, 10, 100),
(5, 5, 'MOUSE-WRLS-BLK', 200, 50, 500),
(5, 23, 'MON-27-BLK', 30, 10, 100),
(5, 24, 'HD-WHT', 65, 25, 200),
(5, 25, 'BPK-BLK', 85, 30, 250),
(6, 1, 'LAP15-BLK-1TB', 15, 5, 50),
(6, 2, 'LAP13-GLD-256', 20, 10, 100),
(6, 3, 'PHX-WHT-128', 35, 15, 150),
(7, 4, 'TAB-BLK-128', 25, 10, 100),
(8, 22, 'CBL-6FT', 150, 50, 500);

-- Insert Price History (20 records)
INSERT INTO price_history (product_id, variant_code, effective_date, price, changed_by_employee_id) VALUES
(1, 'LAP15-BLK-512', '2024-01-01', 1399.99, 5),
(1, 'LAP15-BLK-512', '2024-06-01', 1299.99, 5),
(2, 'LAP13-BLK-256', '2024-01-01', 1099.99, 5),
(2, 'LAP13-BLK-256', '2024-05-01', 999.99, 5),
(3, 'PHX-BLK-128', '2024-01-01', 899.99, 5),
(3, 'PHX-BLK-128', '2024-08-01', 799.99, 5),
(4, 'TAB-GRY-64', '2024-01-01', 649.99, 5),
(4, 'TAB-GRY-64', '2024-07-01', 599.99, 5),
(5, 'MOUSE-WRLS-BLK', '2024-01-01', 34.99, 5),
(5, 'MOUSE-WRLS-BLK', '2024-03-01', 29.99, 5),
(6, 'KEY-MECH-BLK', '2024-01-01', 99.99, 5),
(6, 'KEY-MECH-BLK', '2024-04-01', 89.99, 5),
(7, 'TSH-S-M-BLK', '2024-01-01', 24.99, 5),
(7, 'TSH-S-M-BLK', '2024-02-01', 19.99, 5),
(8, 'JNS-32-BLU', '2024-01-01', 59.99, 5),
(8, 'JNS-32-BLU', '2024-03-01', 49.99, 5),
(9, 'DRS-S-RED', '2024-01-01', 69.99, 5),
(9, 'DRS-S-RED', '2024-04-01', 59.99, 5),
(10, 'BLS-S-WHT', '2024-01-01', 44.99, 5),
(10, 'BLS-S-WHT', '2024-05-01', 39.99, 5);

-- Insert Promotions (5 records)
INSERT INTO promotions (promotion_code, name, start_date, end_date, discount_percentage, created_by_employee_id) VALUES
('PROMO001', 'Summer Sale', '2024-06-01', '2024-08-31', 15.00, 3),
('PROMO002', 'Back to School', '2024-08-01', '2024-09-30', 20.00, 3),
('PROMO003', 'Holiday Special', '2024-11-01', '2024-12-31', 25.00, 3),
('PROMO004', 'New Year Clearance', '2024-12-26', '2025-01-31', 30.00, 3),
('PROMO005', 'Spring Promotion', '2024-03-01', '2024-05-31', 10.00, 3);

-- Insert Promotion Products (15 records - composite keys)
INSERT INTO promotion_products (promotion_id, product_id, variant_code) VALUES
(1, 5, 'MOUSE-WRLS-BLK'),
(1, 6, 'KEY-MECH-BLK'),
(1, 7, 'TSH-S-M-BLK'),
(2, 1, 'LAP15-BLK-512'),
(2, 2, 'LAP13-BLK-256'),
(2, 4, 'TAB-GRY-64'),
(3, 13, 'CHR-BLK'),
(3, 14, 'TBL-6-SEA'),
(3, 23, 'MON-27-BLK'),
(4, 8, 'JNS-32-BLU'),
(4, 9, 'DRS-S-RED'),
(4, 10, 'BLS-S-WHT'),
(5, 15, 'CMK-BLK'),
(5, 16, 'BBL-STD'),
(5, 17, 'TRK-ADULT');

-- Summary: Total records inserted
-- Locations: 15
-- Warehouses: 8
-- Departments: 12
-- Employees: 20
-- Suppliers: 10
-- Categories: 15
-- Products: 25
-- Product Variants: 50
-- Customers: 20
-- Customer Addresses: 40
-- Orders: 30
-- Order Items: 60
-- Supplier Products: 25
-- Reviews: 20
-- Shipments: 15
-- Shipment Items: 25
-- Inventory: 30
-- Price History: 20
-- Promotions: 5
-- Promotion Products: 15
-- TOTAL: ~424 records

-- Note: Some records intentionally have NULL foreign keys to test edge cases
-- Some relationships are circular (indirect) through multiple tables
-- Composite keys test multi-column PK/FK relationships
-- Self-referencing tables test hierarchical relationships

