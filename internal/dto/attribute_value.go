package dto

type CreateAttributeValueRequest struct {
	AttributeID int64  `json:"attribute_id"`
	Value       string `json:"value"`
}

type UpdateAttributeValueRequest struct {
	Value *string `json:"value,omitempty"`
}

type AttributeValueResponse struct {
	ID          int64  `json:"id"`
	AttributeID int64  `json:"attribute_id"`
	Value       string `json:"value"`
	CreatedAt   string `json:"created_at"`
}
