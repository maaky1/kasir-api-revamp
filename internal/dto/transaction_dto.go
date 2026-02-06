package dto

import "time"

type Transaction struct {
	ID        uint                `json:"id"`
	Total     uint                `json:"total"`
	CreatedAt time.Time           `json:"created_at"`
	Details   []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            uint   `json:"id"`
	TransactionID uint   `json:"transaction_id"`
	ProductID     uint   `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      uint   `json:"quantity"`
	Subtotal      uint   `json:"subtotal"`
}
