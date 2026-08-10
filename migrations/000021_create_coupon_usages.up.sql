CREATE TABLE coupon_usages (
    id BIGSERIAL PRIMARY KEY,

    coupon_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,

    used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_coupon_usages_coupon
        FOREIGN KEY (coupon_id)
        REFERENCES coupons(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_coupon_usages_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_coupon_usages_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_coupon_usages_coupon_order
        UNIQUE (coupon_id, order_id)
);

CREATE INDEX idx_coupon_usages_coupon_id
    ON coupon_usages(coupon_id);

CREATE INDEX idx_coupon_usages_user_id
    ON coupon_usages(user_id);