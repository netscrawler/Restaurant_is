package dto

type OrderFilter struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Status    string `json:"status"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}
