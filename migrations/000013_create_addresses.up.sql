CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    title VARCHAR(50) NOT NULL,
    recipient_name VARCHAR(150) NOT NULL,
    phone VARCHAR(30) NOT NULL,

    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address_line TEXT NOT NULL,
    postal_code VARCHAR(30) NOT NULL,

    is_default BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_addresses_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_addresses_user_id
    ON addresses(user_id);

CREATE UNIQUE INDEX idx_addresses_one_default_per_user
    ON addresses(user_id)
    WHERE is_default = TRUE;