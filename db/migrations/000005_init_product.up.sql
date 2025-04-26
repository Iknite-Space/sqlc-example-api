-- 5. PRODUCT CATEGORIES TABLE
CREATE TABLE product_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 1. PRODUCTS TABLE
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES product_categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL,
    regular_price DECIMAL(10,2), -- Regular price for single products
    sale_price DECIMAL(10,2), -- Sale price for single products
    sku VARCHAR(100) UNIQUE, -- SKU for single products
    stock_id INT,
    main_image_url VARCHAR(255), 
    product_gallery_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. PRODUCT VARIATIONS TABLE
CREATE TABLE product_variations (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    variation_name VARCHAR(100) NOT NULL,
    variation_value VARCHAR(100) NOT NULL,
    regular_price DECIMAL(10,2) NOT NULL,
    sale_price DECIMAL(10,2),
    sku VARCHAR(100) UNIQUE,
    stock_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. STOCK TABLE
CREATE TABLE stock (
    id SERIAL PRIMARY KEY,
    quantity INT NOT NULL,
    low_stock_threshold INT DEFAULT 5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. PRODUCT GALLERY TABLE
CREATE TABLE product_gallery (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    image_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

