-- ============================================
-- Simple-Commerce Schema
-- PostgreSQL
-- ============================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ============================================
-- CUSTOMER
-- ============================================
CREATE TABLE customers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name   VARCHAR(255) NOT NULL,
    phone       VARCHAR(20),
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_customers_email UNIQUE (email),
    CONSTRAINT chk_customers_status CHECK (
        status IN ('pending', 'active')
    )
);

-- ============================================
-- SELLER
-- ============================================
CREATE TABLE sellers (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    shop_name     VARCHAR(255) NOT NULL,
    phone         VARCHAR(20),
    status        VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_sellers_email UNIQUE (email),
    CONSTRAINT chk_sellers_status CHECK (
        status IN ('pending', 'active')
    )
);

-- ============================================
-- ACTIVATION CODE
-- ============================================
CREATE TABLE activation_codes (
    id          BIGSERIAL PRIMARY KEY,
    code        VARCHAR(64) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    type        VARCHAR(50) NOT NULL,   -- 'email_verification' | 'password_reset'
    user_type   VARCHAR(20) NOT NULL,   -- 'customer' | 'seller'
    expires_at  TIMESTAMP NOT NULL,
    used_at     TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_activation_code UNIQUE (email, code),
    CONSTRAINT chk_activation_codes_type CHECK (
        type IN ('email_verification', 'password_reset')
    ),
    CONSTRAINT chk_activation_codes_user_type CHECK (
        user_type IN ('customer', 'seller')
    )
);

-- ============================================
-- REFRESH TOKEN
-- ============================================
CREATE TABLE refresh_tokens (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID NOT NULL,
    user_type   VARCHAR(20) NOT NULL,
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMP NOT NULL,
    revoked_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_refresh_tokens_token UNIQUE (token_hash),
    CONSTRAINT chk_refresh_tokens_user_type CHECK (
        user_type IN ('customer', 'seller')
    )
);

-- ============================================
-- CATEGORY
-- ============================================
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    parent_id   BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_categories_slug UNIQUE (slug)
);

-- ============================================
-- PRODUCT
-- ============================================
CREATE TABLE products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id   UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name        VARCHAR(255) NOT NULL,
    sku         VARCHAR(100) NOT NULL,
    description TEXT,
    price       DECIMAL(15, 2) NOT NULL CHECK (price >= 0),
    image_url   TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_products_sku UNIQUE (sku)
);

-- ============================================
-- INVENTORY
-- version for optimistic locking (concurrency control)
-- ============================================
CREATE TABLE inventories (
    id                BIGSERIAL PRIMARY KEY,
    product_id        UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    stock_quantity    INT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    reserved_quantity INT NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    version           INT NOT NULL DEFAULT 1,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_inventories_product_id UNIQUE (product_id),
    CONSTRAINT chk_reserved_lte_stock CHECK (reserved_quantity <= stock_quantity)
);

CREATE OR REPLACE FUNCTION create_inventory_for_product()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO inventories (
        product_id,
        stock_quantity,
        reserved_quantity,
        version,
        created_at,
        updated_at
    )
    VALUES (
        NEW.id,
        0,
        0,
        1,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
    ON CONFLICT (product_id) DO NOTHING;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_inventory_after_product_insert
AFTER INSERT ON products
FOR EACH ROW
EXECUTE FUNCTION create_inventory_for_product();

-- ============================================
-- SHOPPING CART
-- ============================================
CREATE TABLE shopping_carts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_shopping_carts_customer_id UNIQUE (customer_id)
);

-- ============================================
-- CART ITEM
-- ============================================
CREATE TABLE cart_items (
    id             BIGSERIAL PRIMARY KEY,
    cart_id        UUID NOT NULL REFERENCES shopping_carts(id) ON DELETE CASCADE,
    product_id     UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity       INT NOT NULL CHECK (quantity > 0),
    price_snapshot DECIMAL(15, 2) NOT NULL CHECK (price_snapshot >= 0),
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_cart_items_cart_product UNIQUE (cart_id, product_id)
);

-- ============================================
-- BUSINESS NUMBER COUNTER
-- Atomic daily counters for human-readable order and transaction numbers.
-- ============================================
CREATE TABLE business_number_counters (
    name         VARCHAR(50) NOT NULL,
    counter_date DATE NOT NULL,
    value        BIGINT NOT NULL CHECK (value > 0),
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_business_number_counters PRIMARY KEY (name, counter_date)
);

-- ============================================
-- ORDER
-- ============================================
CREATE TABLE orders (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number     VARCHAR(50) NOT NULL,
    customer_id      UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    status           VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_amount     DECIMAL(15, 2) NOT NULL CHECK (total_amount >= 0),
    shipping_address TEXT NOT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_orders_order_number UNIQUE (order_number),
    CONSTRAINT chk_orders_status CHECK (
        status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled')
    )
);

-- ============================================
-- ORDER ITEM
-- ============================================
CREATE TABLE order_items (
    id         BIGSERIAL PRIMARY KEY,
    order_id   UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    seller_id  UUID NOT NULL REFERENCES sellers(id) ON DELETE RESTRICT,
    quantity   INT NOT NULL CHECK (quantity > 0),
    price      DECIMAL(15, 2) NOT NULL CHECK (price >= 0),
    subtotal   DECIMAL(15, 2) NOT NULL CHECK (subtotal >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TRANSACTION
-- ============================================
CREATE TABLE transactions (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id           UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    transaction_number VARCHAR(50) NOT NULL,
    payment_method     VARCHAR(50) NOT NULL,
    status             VARCHAR(50) NOT NULL DEFAULT 'pending',
    amount             DECIMAL(15, 2) NOT NULL CHECK (amount >= 0),
    paid_at            TIMESTAMP,
    created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_transactions_number UNIQUE (transaction_number),
    CONSTRAINT uq_transactions_order_id UNIQUE (order_id),
    CONSTRAINT chk_transactions_status CHECK (
        status IN ('pending', 'processing', 'success', 'failed', 'expired', 'refunded')
    ),
    CONSTRAINT chk_transactions_payment_method CHECK (
        payment_method IN ('credit_card', 'bank_transfer', 'e_wallet', 'cod')
    )
);

-- ============================================
-- EMAIL QUEUE
-- ============================================
CREATE TABLE email_queues (
    id              BIGSERIAL PRIMARY KEY,
    recipient_email VARCHAR(255) NOT NULL,
    subject         VARCHAR(255) NOT NULL,
    body_html       TEXT NOT NULL,
    type            VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count     INT NOT NULL DEFAULT 0,
    max_retries     INT NOT NULL DEFAULT 3,
    error_message   TEXT,
    scheduled_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sent_at         TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_email_queues_status CHECK (
        status IN ('pending', 'processing', 'sent', 'failed')
    ),
    CONSTRAINT chk_email_queues_type CHECK (
        type IN ('email_verification', 'password_reset', 'order_confirmation', 'payment_success', 'shipping_update')
    )
);

-- ============================================
-- INDEXES
-- ============================================

-- sellers
CREATE INDEX idx_sellers_shop_name_trgm ON sellers USING gin (shop_name gin_trgm_ops);

-- activation_codes
CREATE INDEX idx_activation_codes_email_user_type_created_at ON activation_codes(email, user_type, created_at DESC);
CREATE INDEX idx_activation_codes_code_active ON activation_codes(code, expires_at) WHERE used_at IS NULL;
CREATE INDEX idx_activation_codes_expires_at ON activation_codes(expires_at);

-- refresh_tokens
CREATE INDEX idx_refresh_tokens_user_id_created_at ON refresh_tokens(user_id, created_at DESC);
CREATE INDEX idx_refresh_tokens_active_token_hash ON refresh_tokens(token_hash) WHERE revoked_at IS NULL;

-- categories
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_is_active ON categories(is_active);

-- products
CREATE INDEX idx_products_seller_id ON products(seller_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_active_created_at_id ON products(created_at DESC, id DESC) WHERE is_active = TRUE;
CREATE INDEX idx_products_active_price_id ON products(price, id) WHERE is_active = TRUE;
CREATE INDEX idx_products_active_name ON products USING gin (name gin_trgm_ops) WHERE is_active = TRUE;

-- cart_items
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);

-- orders
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_customer_status_created_at ON orders(customer_id, status, created_at DESC);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- order_items
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_seller_id ON order_items(seller_id);

-- transactions
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

-- email_queues
CREATE INDEX idx_email_queues_status_scheduled ON email_queues(status, scheduled_at);

-- ============================================
-- SEED DATA
-- Demo password for seeded customer and sellers: password123
-- ============================================

INSERT INTO customers (
    id,
    email,
    password_hash,
    full_name,
    phone,
    status,
    created_at,
    updated_at
)
VALUES (
    '11111111-1111-4111-8111-111111111111',
    'customer.demo@simple-commerce.test',
    '$2a$10$zr.I2oj8buuLxf43LT27weVFjK6OiEjyPhAK7/kUmlAGrOSSdgT8e',
    'Demo Customer',
    '+6281200000001',
    'active',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO UPDATE SET
    full_name = EXCLUDED.full_name,
    phone = EXCLUDED.phone,
    status = EXCLUDED.status,
    updated_at = CURRENT_TIMESTAMP;

INSERT INTO sellers (
    id,
    email,
    password_hash,
    shop_name,
    phone,
    status,
    created_at,
    updated_at
)
VALUES
    (
        '22222222-2222-4222-8222-222222222222',
        'tech.seller@simple-commerce.test',
        '$2a$10$zr.I2oj8buuLxf43LT27weVFjK6OiEjyPhAK7/kUmlAGrOSSdgT8e',
        'Tech Corner',
        '+6281200000002',
        'active',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '33333333-3333-4333-8333-333333333333',
        'lifestyle.seller@simple-commerce.test',
        '$2a$10$zr.I2oj8buuLxf43LT27weVFjK6OiEjyPhAK7/kUmlAGrOSSdgT8e',
        'Lifestyle Market',
        '+6281200000003',
        'active',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (email) DO UPDATE SET
    shop_name = EXCLUDED.shop_name,
    phone = EXCLUDED.phone,
    status = EXCLUDED.status,
    updated_at = CURRENT_TIMESTAMP;

INSERT INTO categories (
    id,
    parent_id,
    name,
    slug,
    is_active,
    created_at,
    updated_at
)
VALUES
    (1, NULL, 'Electronics', 'electronics', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, NULL, 'Fashion', 'fashion', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, NULL, 'Books', 'books', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, NULL, 'Home & Kitchen', 'home-kitchen', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (5, NULL, 'Sports', 'sports', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (slug) DO UPDATE SET
    parent_id = EXCLUDED.parent_id,
    name = EXCLUDED.name,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

SELECT setval(
    pg_get_serial_sequence('categories', 'id'),
    COALESCE((SELECT MAX(id) FROM categories), 1)
);

INSERT INTO products (
    id,
    seller_id,
    category_id,
    name,
    sku,
    description,
    price,
    image_url,
    is_active,
    created_at,
    updated_at
)
VALUES
    (
        '44444444-4444-4444-8444-444444444401',
        '22222222-2222-4222-8222-222222222222',
        1,
        'Laptop Pro 14',
        'TECH-LAPTOP-PRO-14',
        'Portable laptop for development, office work, and content creation.',
        15000000.00,
        'https://example.com/images/laptop-pro-14.jpg',
        TRUE,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '44444444-4444-4444-8444-444444444402',
        '22222222-2222-4222-8222-222222222222',
        1,
        'Wireless Headphones',
        'TECH-HEADPHONES-WIRELESS',
        'Noise-isolating wireless headphones for daily use.',
        850000.00,
        'https://example.com/images/wireless-headphones.jpg',
        TRUE,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '44444444-4444-4444-8444-444444444403',
        '22222222-2222-4222-8222-222222222222',
        1,
        'Mechanical Keyboard',
        'TECH-KEYBOARD-MECHANICAL',
        'Compact mechanical keyboard with tactile switches.',
        1200000.00,
        'https://example.com/images/mechanical-keyboard.jpg',
        TRUE,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '44444444-4444-4444-8444-444444444404',
        '33333333-3333-4333-8333-333333333333',
        2,
        'Everyday T-Shirt',
        'FASHION-TSHIRT-EVERYDAY',
        'Soft cotton t-shirt for casual everyday wear.',
        180000.00,
        'https://example.com/images/everyday-tshirt.jpg',
        TRUE,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '44444444-4444-4444-8444-444444444405',
        '33333333-3333-4333-8333-333333333333',
        3,
        'Practical Go Backend',
        'BOOK-GO-BACKEND',
        'Backend engineering book covering API design, data access, and reliability.',
        250000.00,
        'https://example.com/images/practical-go-backend.jpg',
        TRUE,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '44444444-4444-4444-8444-444444444406',
        '33333333-3333-4333-8333-333333333333',
        4,
        'Coffee Maker',
        'HOME-COFFEE-MAKER',
        'Automatic coffee maker for home kitchens.',
        650000.00,
        'https://example.com/images/coffee-maker.jpg',
        TRUE,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (sku) DO UPDATE SET
    seller_id = EXCLUDED.seller_id,
    category_id = EXCLUDED.category_id,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    image_url = EXCLUDED.image_url,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

UPDATE inventories
SET
    stock_quantity = seed.stock_quantity,
    reserved_quantity = 0,
    version = version + 1,
    updated_at = CURRENT_TIMESTAMP
FROM (
    VALUES
        ('44444444-4444-4444-8444-444444444401'::UUID, 15),
        ('44444444-4444-4444-8444-444444444402'::UUID, 40),
        ('44444444-4444-4444-8444-444444444403'::UUID, 25),
        ('44444444-4444-4444-8444-444444444404'::UUID, 100),
        ('44444444-4444-4444-8444-444444444405'::UUID, 60),
        ('44444444-4444-4444-8444-444444444406'::UUID, 20)
) AS seed(product_id, stock_quantity)
WHERE inventories.product_id = seed.product_id;


INSERT INTO public.sellers
(id, email, password_hash, shop_name, phone, status, created_at, updated_at)
values
('44444444-4444-4444-8444-444444444444'::uuid, 'sports.seller@simple-commerce.test', '$2a$10$zr.I2oj8buuLxf43LT27weVFjK6OiEjyPhAK7/kUmlAGrOSSdgT8e', 'Sports Arena', '+6281200000004', 'active', '2026-05-16 16:10:03.518', '2026-05-16 16:10:03.518'),
('55555555-5555-4555-8555-555555555555'::uuid, 'books.seller@simple-commerce.test', '$2a$10$zr.I2oj8buuLxf43LT27weVFjK6OiEjyPhAK7/kUmlAGrOSSdgT8e', 'Book Haven', '+6281200000005', 'active', '2026-05-16 16:10:03.518', '2026-05-16 16:10:03.518'),
('66666666-6666-4666-8666-666666666666'::uuid, 'beauty.seller@simple-commerce.test', '$2a$10$zr.I2oj8buuLxf43LT27weVFjK6OiEjyPhAK7/kUmlAGrOSSdgT8e', 'Beauty Central', '+6281200000006', 'active', '2026-05-16 16:10:03.518', '2026-05-16 16:10:03.518'),
('77777777-7777-4777-8777-777777777777'::uuid, 'gaming.seller@simple-commerce.test', '$2a$10$zr.I2oj8buuLxf43LT27weVFjK6OiEjyPhAK7/kUmlAGrOSSdgT8e', 'Gaming Zone', '+6281200000007', 'active', '2026-05-16 16:10:03.518', '2026-05-16 16:10:03.518'),
('88888888-8888-4888-8888-888888888888'::uuid, 'groceries.seller@simple-commerce.test', '$2a$10$zr.I2oj8buuLxf43LT27weVFjK6OiEjyPhAK7/kUmlAGrOSSdgT8e', 'Groceries Hub', '+6281200000008', 'active', '2026-05-16 16:10:03.518', '2026-05-16 16:10:03.518');

INSERT INTO public.categories
(id, parent_id, "name", slug, is_active, created_at, updated_at)
values
(6, NULL, 'Beauty', 'beauty', true, '2026-05-16 16:10:03.547', '2026-05-16 16:10:03.547'),
(7, NULL, 'Gaming', 'gaming', true, '2026-05-16 16:10:03.547', '2026-05-16 16:10:03.547'),
(8, NULL, 'Groceries', 'groceries', true, '2026-05-16 16:10:03.547', '2026-05-16 16:10:03.547'),
(9, NULL, 'Office Supplies', 'office-supplies', true, '2026-05-16 16:10:03.547', '2026-05-16 16:10:03.547'),
(10, NULL, 'Automotive', 'automotive', true, '2026-05-16 16:10:03.547', '2026-05-16 16:10:03.547');
