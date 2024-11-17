CREATE TABLE public.customer (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(100) NOT NULL,
    balance NUMERIC(10,2),
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100)
);

CREATE TABLE public.category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100)
);

CREATE TABLE public.product (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(100),
    price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL,
    category_id INT,  
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100),
    FOREIGN KEY (category_id) REFERENCES public.category(id)  
);

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
CREATE TABLE seller (
    id SERIAL PRIMARY KEY,
    seller_name VARCHAR(255) NOT NULL,
    seller_balance NUMERIC(10,2) DEFAULT 0.00,
    created TIMESTAMP NOT NULL,
    created_by VARCHAR(100),
    updated  TIMESTAMP,
    updated_by VARCHAR(100),
);

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
