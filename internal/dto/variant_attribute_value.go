package dto

type CreateVariantAttributeValueRequest struct {
	VariantID        int64 `json:"variant_id"`
	AttributeValueID int64 `json:"attribute_value_id"`
}

type VariantAttributeValueResponse struct {
	VariantID        int64 `json:"variant_id"`
	AttributeValueID int64 `json:"attribute_value_id"`
}
