CREATE TABLE product_variants (
    id BIGSERIAL PRIMARY KEY,

    product_id BIGINT NOT NULL,

    sku VARCHAR(100) UNIQUE NOT NULL,

    price NUMERIC(12,2) NOT NULL,
    discount_price NUMERIC(12,2) NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_product_variants_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_product_variants_price
        CHECK (price >= 0),

    CONSTRAINT chk_product_variants_discount_price
        CHECK (
            discount_price IS NULL
            OR (
                discount_price >= 0
                AND discount_price <= price
            )
        )
);

CREATE INDEX idx_product_variants_product_id
    ON product_variants(product_id);