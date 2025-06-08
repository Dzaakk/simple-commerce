CREATE TABLE public.customer (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    gender SMALLINT,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    profile_picture TEXT,
    date_of_birth DATE,
    balance NUMERIC(10,2) DEFAULT 0,
    last_login TIMESTAMP,        
    status SMALLINT,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL DEFAULT 0, 
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by INT               
    -- email_verified BOOLEAN DEFAULT FALSE,
    -- phone_verified BOOLEAN DEFAULT FALSE,
    -- Optional future fields:
    -- referral_code VARCHAR(50),
);

CREATE INDEX idx_customer_username ON public.customer (username);
CREATE INDEX idx_customer_email ON public.customer (email);
CREATE INDEX idx_customer_phone_number ON public.customer (phone_number);

CREATE TABLE public.category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100)
);
CREATE INDEX idx_category_name ON public.category (name);

CREATE TABLE public.product (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(100),
    price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL,
    category_id INT,  
    seller_id INT,
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100),
    FOREIGN KEY (category_id) REFERENCES public.category(id)  
);

CREATE INDEX idx_product_name ON public.product (product_name);
CREATE INDEX idx_product_price ON public.product (price);
CREATE INDEX idx_product_category_id ON public.product (category_id);
CREATE INDEX idx_product_seller_id ON public.product (seller_id);

CREATE TABLE public.shopping_cart (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    status varchar(20) NOT NULL,
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100),
    FOREIGN KEY (customer_id) REFERENCES public.customer(id)
);
CREATE INDEX idx_shopping_cart_customer_id ON public.shopping_cart (customer_id);

CREATE TABLE public.shopping_cart_item (
    cart_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100),
    PRIMARY KEY (cart_id, product_id),
    FOREIGN KEY (cart_id) REFERENCES shopping_cart(id),
    FOREIGN KEY (product_id) REFERENCES public.product(id)
);
CREATE INDEX idx_shopping_cart_item_cart_id ON public.shopping_cart_item (cart_id);
CREATE INDEX idx_shopping_cart_item_product_id ON public.shopping_cart_item (product_id);

CREATE TABLE public.transaction (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    cart_id INT NOT NULL,
    total_amount NUMERIC(10, 2) NOT NULL,
    transaction_date TIMESTAMP,
    status VARCHAR(50),
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100),
    FOREIGN KEY (customer_id) REFERENCES public.customer(id),
    FOREIGN KEY (cart_id) REFERENCES public.shopping_cart(id)
);
CREATE INDEX idx_transaction_customer_id ON public.transaction (customer_id);
CREATE INDEX idx_transaction_cart_id ON public.transaction (cart_id);
CREATE INDEX idx_transaction_status ON public.transaction (status);

CREATE TABLE public.history_transaction(
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    product_name VARCHAR(100),
    price NUMERIC(10, 2) NOT NULL,
    quantity INT NOT NULL,
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100),
    status varchar(20) NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES public.customer(id)
);
CREATE INDEX idx_history_transaction_customer_id ON public.history_transaction (customer_id);
CREATE INDEX idx_history_transaction_status ON public.history_transaction (status);


CREATE TABLE seller (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance NUMERIC(10,2) DEFAULT 0.00,
    status varchar(1) NOT NULL, 
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100),
);
CREATE INDEX idx_seller_username ON public.seller (username);
CREATE INDEX idx_seller_email ON public.seller (email);

-- timestamp indexing
-- CREATE INDEX idx_customer_created ON public.customer (created);
-- CREATE INDEX idx_product_created ON public.product (created);
-- CREATE INDEX idx_shopping_cart_created ON public.shopping_cart (created);
-- CREATE INDEX idx_transaction_created ON public.transaction (created);
-- CREATE INDEX idx_history_transaction_created ON public.history_transaction (created);
-- CREATE INDEX idx_seller_created ON public.seller (created);

CREATE TABLE customer_activation_code (
    customer_id BIGINT,
    code_activation VARCHAR(6),
    is_used BOOLEAN,
    created TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
)

CREATE TABLE seller_activation_code (
    seller_id BIGINT,
    code_activation VARCHAR(6),
    is_used BOOLEAN,
    created TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
)

INSERT INTO public.category (name, created_by, created)
VALUES ('Electronics', 'Admin', now()),
       ('Clothing', 'Admin', now()),
       ('Books', 'Admin', now()),
       ('Home & Kitchen', 'Admin', now());

-- Products for Electronics category
INSERT INTO public.product (product_name, price, stock, category_id, created_by, created)
VALUES ('Laptop', 15000000, 10, 1, 'Admin', now()),
       ('Smartphone', 7000000, 20, 1, 'Admin', now()),
       ('Headphones', 800000, 30, 1, 'Admin', now()),
       ('Tablet', 4000000, 15, 1, 'Admin', now());

-- Products for Clothing category
INSERT INTO public.product (product_name, price, stock, category_id, created_by, created)
VALUES ('T-Shirt', 200000, 50, 2, 'Admin', now()),
       ('Jeans', 500000, 30, 2, 'Admin', now()),
       ('Sneakers', 1000000, 25, 2, 'Admin', now()),
       ('Dress', 900000, 20, 2, 'Admin', now());

-- Products for Books category
INSERT INTO public.product (product_name, price, stock, category_id, created_by, created)
VALUES ('Programming Book', 150000, 100, 3, 'Admin', now()),
       ('Fiction Book', 80000, 80, 3, 'Admin', now()),
       ('Self-Help Book', 75000, 75, 3, 'Admin', now()),
       ('Cookbook', 50000, 60, 3, 'Admin', now());

-- Products for Home & Kitchen category
INSERT INTO public.product (product_name, price, stock, category_id, created_by, created)
VALUES ('Coffee Maker', 5000000, 15, 4, 'Admin', now()),
       ('Vacuum Cleaner', 2500000, 10, 4, 'Admin', now()),
       ('Knife Set', 700000, 20, 4, 'Admin', now()),
       ('Cookware Set', 600000, 12, 4, 'Admin', now());


