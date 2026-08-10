CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,

    order_id BIGINT UNIQUE NOT NULL,

    amount NUMERIC(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,

    provider VARCHAR(30) NOT NULL,
    transaction_id VARCHAR(255) UNIQUE,

    status VARCHAR(30) NOT NULL,

    paid_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_payments_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_payments_amount
        CHECK (amount > 0),

    CONSTRAINT chk_payments_provider
        CHECK (
            provider IN ('mock', 'stripe')
        ),

    CONSTRAINT chk_payments_status
        CHECK (
            status IN (
                'pending',
                'processing',
                'paid',
                'failed',
                'refunded',
                'cancelled'
            )
        )
);

CREATE INDEX idx_payments_status
    ON payments(status);