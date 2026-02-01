package dto

import "time"

type Product struct {
	CategoryID *uint  `json:"category_id,omitempty"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
}

type UpdateProduct struct {
	CategoryID *uint   `json:"category_id,omitempty"`
	Name       *string `json:"name,omitempty"`
	Price      *int    `json:"price,omitempty"`
	Stock      *int    `json:"stock,omitempty"`
}

type ProductResponse struct {
	ID         uint      `json:"id"`
	CategoryID uint      `json:"category_id"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	Stock      int       `json:"stock"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductDetailResponse struct {
	ID           uint      `json:"id"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Name         string    `json:"name"`
	Price        int       `json:"price"`
	Stock        int       `json:"stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
