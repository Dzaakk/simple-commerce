-- ============================================
-- Simple-Commerce Schema
-- PostgreSQL
-- ============================================

-- ============================================
-- CUSTOMER
-- ============================================
CREATE TABLE customers (
    id          UUID PRIMARY KEY,
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
    id            UUID PRIMARY KEY,
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

    CONSTRAINT uq_activation_code UNIQUE (email, code)
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

-- -- customers
-- CREATE INDEX idx_customers_email ON customers(email);

-- -- sellers
-- CREATE INDEX idx_sellers_email ON sellers(email);

-- -- activation_codes
-- CREATE INDEX idx_activation_codes_email ON activation_codes(email);
-- CREATE INDEX idx_activation_codes_expires_at ON activation_codes(expires_at);

-- -- categories
-- CREATE INDEX idx_categories_parent_id ON categories(parent_id);
-- CREATE INDEX idx_categories_slug ON categories(slug);

-- -- products
-- CREATE INDEX idx_products_seller_id ON products(seller_id);
-- CREATE INDEX idx_products_category_id ON products(category_id);
-- CREATE INDEX idx_products_is_active ON products(is_active);

-- -- inventories
-- CREATE INDEX idx_inventories_product_id ON inventories(product_id);

-- -- cart_items
-- CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
-- CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);

-- -- orders
-- CREATE INDEX idx_orders_customer_id ON orders(customer_id);
-- CREATE INDEX idx_orders_status ON orders(status);
-- CREATE INDEX idx_orders_created_at ON orders(created_at);

-- -- order_items
-- CREATE INDEX idx_order_items_order_id ON order_items(order_id);
-- CREATE INDEX idx_order_items_seller_id ON order_items(seller_id);

-- -- transactions
-- CREATE INDEX idx_transactions_order_id ON transactions(order_id);
-- CREATE INDEX idx_transactions_status ON transactions(status);

-- -- email_queues
-- CREATE INDEX idx_email_queues_status_scheduled ON email_queues(status, scheduled_at);
