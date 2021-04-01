CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS cart
(
    id       UUID NOT NULL,
    buyer_id UUID NOT NULL,
    CONSTRAINT pk_cart PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS product_line
(
    id      UUID DEFAULT uuid_generate_v4(),
    product_id UUID  NOT NULL,
    cart_id UUID  NOT NULL,
    name    text  NOT NULL,
    price   FLOAT NOT NULL,
    qty     INT   NOT NULL,
    CONSTRAINT pk_product_line PRIMARY KEY (id),
    CONSTRAINT fk_product_cart FOREIGN KEY (cart_id) REFERENCES cart (id)
);
