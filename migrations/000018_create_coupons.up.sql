CREATE TABLE coupons (
    id BIGSERIAL PRIMARY KEY,

    code VARCHAR(50) UNIQUE NOT NULL,

    discount_type VARCHAR(20) NOT NULL,
    discount_value NUMERIC(12,2) NOT NULL,

    minimum_order_amount NUMERIC(12,2),
    maximum_discount NUMERIC(12,2),

    usage_limit INTEGER,
    used_count INTEGER NOT NULL DEFAULT 0,

    starts_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_coupons_discount_type
        CHECK (discount_type IN ('percentage', 'fixed')),

    CONSTRAINT chk_coupons_discount_value
        CHECK (discount_value > 0),

    CONSTRAINT chk_coupons_minimum_order_amount
        CHECK (
            minimum_order_amount IS NULL
            OR minimum_order_amount >= 0
        ),

    CONSTRAINT chk_coupons_maximum_discount
        CHECK (
            maximum_discount IS NULL
            OR maximum_discount >= 0
        ),

    CONSTRAINT chk_coupons_usage_limit
        CHECK (
            usage_limit IS NULL
            OR usage_limit > 0
        ),

    CONSTRAINT chk_coupons_used_count
        CHECK (used_count >= 0),

    CONSTRAINT chk_coupons_dates
        CHECK (expires_at > starts_at)
);

CREATE INDEX idx_coupons_is_active
    ON coupons(is_active);

CREATE INDEX idx_coupons_starts_at
    ON coupons(starts_at);

CREATE INDEX idx_coupons_expires_at
    ON coupons(expires_at);