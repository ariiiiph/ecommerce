CREATE TABLE attribute_values (
    id BIGSERIAL PRIMARY KEY,

    attribute_id BIGINT NOT NULL,
    value VARCHAR(100) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_attribute_values_attribute
        FOREIGN KEY (attribute_id)
        REFERENCES attributes(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_attribute_values_attribute_value
        UNIQUE (attribute_id, value)
);

CREATE INDEX idx_attribute_values_attribute_id
    ON attribute_values(attribute_id);