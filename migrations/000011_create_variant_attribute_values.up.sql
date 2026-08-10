CREATE TABLE variant_attribute_values (
    variant_id BIGINT NOT NULL,
    attribute_value_id BIGINT NOT NULL,

    PRIMARY KEY (variant_id, attribute_value_id),

    CONSTRAINT fk_variant_attribute_values_variant
        FOREIGN KEY (variant_id)
        REFERENCES product_variants(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_variant_attribute_values_attribute_value
        FOREIGN KEY (attribute_value_id)
        REFERENCES attribute_values(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_variant_attribute_values_attribute_value_id
    ON variant_attribute_values(attribute_value_id);