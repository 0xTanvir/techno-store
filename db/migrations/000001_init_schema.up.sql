-- Create brands table
CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    status_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    parent_id INT REFERENCES categories(id) ON DELETE SET NULL,
    sequence INT NOT NULL,
    status_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create suppliers table
CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    status_id INT NOT NULL,
    is_verified_supplier BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    specifications TEXT,
    brand_id INT REFERENCES brands(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    supplier_id INT REFERENCES suppliers(id) ON DELETE CASCADE,
    unit_price DECIMAL(10, 2) NOT NULL,
    discount_price DECIMAL(10, 2),
    tags TEXT,
    status_id INT NOT NULL,
    CONSTRAINT unique_supplier_product UNIQUE(supplier_id, name)
);

-- Create product_stocks table
CREATE TABLE product_stocks (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    stock_quantity INT CHECK (stock_quantity >= 0),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 20 random products insertion
DO $$
DECLARE
    brand_id_1 INT;
    brand_id_2 INT;
    brand_id_3 INT;
    category_id_1 INT;
    category_id_2 INT;
    category_id_3 INT;
    category_id_4 INT;
    supplier_id_1 INT;
    supplier_id_2 INT;
    product_id_1 INT;
    product_id_2 INT;
    product_id_3 INT;
    product_id_4 INT;
    product_id_5 INT;

BEGIN
    -- Insert Brands
    INSERT INTO brands (name, status_id) VALUES ('Brand A', 1) RETURNING id INTO brand_id_1;
    INSERT INTO brands (name, status_id) VALUES ('Brand B', 1) RETURNING id INTO brand_id_2;
    INSERT INTO brands (name, status_id) VALUES ('Brand C', 1) RETURNING id INTO brand_id_3;

    -- Insert Categories
    INSERT INTO categories (name, parent_id, sequence, status_id) VALUES ('Mobile', NULL, 1, 1) RETURNING id INTO category_id_1;
    INSERT INTO categories (name, parent_id, sequence, status_id) VALUES ('IOS', category_id_1, 1, 1) RETURNING id INTO category_id_2;
    INSERT INTO categories (name, parent_id, sequence, status_id) VALUES ('Android', category_id_1, 2, 1) RETURNING id INTO category_id_3;
    INSERT INTO categories (name, parent_id, sequence, status_id) VALUES ('Watch', NULL, 2, 1) RETURNING id INTO category_id_4;

    -- Insert Suppliers
    INSERT INTO suppliers (name, email, phone, status_id, is_verified_supplier) VALUES ('Supplier 1', 's1@example.com', '01230303', 1, TRUE) RETURNING id INTO supplier_id_1;
    INSERT INTO suppliers (name, email, phone, status_id, is_verified_supplier) VALUES ('Supplier 2', 's2@example.com', '02304040', 1, TRUE) RETURNING id INTO supplier_id_2;

    -- Insert Random Products
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 1', 'Description 1', 'Specs 1', brand_id_1, category_id_1, supplier_id_1, 10000.00, 220.00, 'tag', 1) RETURNING id INTO product_id_1;
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 2', 'Description 1', 'Specs 1', brand_id_1, category_id_1, supplier_id_1, 10000.00, 2020.00, 'tag', 1) RETURNING id INTO product_id_2;
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 3', 'Description 1', 'Specs 1', brand_id_1, category_id_1, supplier_id_1, 10050.00, 2000.00, 'tag', 1) RETURNING id INTO product_id_3;
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 4', 'Description 1', 'Specs 1', brand_id_1, category_id_1, supplier_id_1, 1050.00, 2000.00, 'tag', 1) RETURNING id INTO product_id_4;
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 5', 'Description 1', 'Specs 1', brand_id_1, category_id_1, supplier_id_1, 1005.00, 2000.00, 'tag', 1) RETURNING id INTO product_id_5;
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 6', 'Description 1', 'Specs 1', brand_id_2, category_id_2, supplier_id_1, 10000.00, 220.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 7', 'Description 1', 'Specs 1', brand_id_2, category_id_2, supplier_id_1, 10000.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 8', 'Description 1', 'Specs 1', brand_id_2, category_id_2, supplier_id_1, 10000.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 9', 'Description 1', 'Specs 1', brand_id_2, category_id_2, supplier_id_1, 10000.00, 2020.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 10', 'Description 1', 'Specs 1', brand_id_2, category_id_2, supplier_id_1, 10500.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 11', 'Description 1', 'Specs 1', brand_id_2, category_id_3, supplier_id_2, 10000.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 12', 'Description 1', 'Specs 1', brand_id_2, category_id_3, supplier_id_2, 10000.00, 220.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 13', 'Description 1', 'Specs 1', brand_id_1, category_id_3, supplier_id_2, 10000.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 14', 'Description 1', 'Specs 1', brand_id_1, category_id_3, supplier_id_2, 10500.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 15', 'Description 1', 'Specs 1', brand_id_3, category_id_3, supplier_id_2, 10000.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 16', 'Description 1', 'Specs 1', brand_id_3, category_id_4, supplier_id_2, 1300.00, 2020.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 17', 'Description 1', 'Specs 1', brand_id_3, category_id_4, supplier_id_2, 1020.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 18', 'Description 1', 'Specs 1', brand_id_3, category_id_4, supplier_id_2, 1030.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 19', 'Description 1', 'Specs 1', brand_id_3, category_id_4, supplier_id_2, 101.00, 2000.00, 'tag', 1);
    INSERT INTO products (name, description, specifications, brand_id, category_id, supplier_id, unit_price, discount_price, tags, status_id)
        VALUES ('Product 20', 'Description 1', 'Specs 1', brand_id_3, category_id_4, supplier_id_2, 10000.00, 2000.00, 'tag', 1);

    -- Insert Random Products Stocks
    INSERT INTO product_stocks (product_id, stock_quantity) VALUES (product_id_1, 10);
    INSERT INTO product_stocks (product_id, stock_quantity) VALUES (product_id_2, 2);
    INSERT INTO product_stocks (product_id, stock_quantity) VALUES (product_id_3, 7);
    INSERT INTO product_stocks (product_id, stock_quantity) VALUES (product_id_4, 3);
    INSERT INTO product_stocks (product_id, stock_quantity) VALUES (product_id_5, 4);
END $$;