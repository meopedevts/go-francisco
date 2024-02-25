CREATE TABLE  customers (
    id INT4 NOT NULL,
    "limit" INT4 NOT NULL,
    balance INT4 NOT NULL DEFAULT 0,
    CONSTRAINT customers_pk PRIMARY KEY (id)
);

CREATE TABLE transactions (
    id SERIAL NOT NULL,
    customer_id INT4 NOT NULL,
    amount INT4 NOT NULL DEFAULT 0,
    "type" CHAR NULL,
    description VARCHAR(50) NULL,
    created_at TIMESTAMP NULL,
    CONSTRAINT transactions_pk PRIMARY KEY (id)
);

INSERT INTO customers (id, "limit", balance) VALUES(1, 100000, 0);
INSERT INTO customers (id, "limit", balance) VALUES(2, 80000, 0);
INSERT INTO customers (id, "limit", balance) VALUES(3, 1000000, 0);
INSERT INTO customers (id, "limit", balance) VALUES(4, 10000000, 0);
INSERT INTO customers (id, "limit", balance) VALUES(5, 500000, 0);