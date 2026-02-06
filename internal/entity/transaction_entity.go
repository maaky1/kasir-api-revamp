package entity

import "time"

type Transaction struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	TotalAmount uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type TransactionDetail struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	TransactionID uint      `gorm:"not null"`
	ProductID     uint      `gorm:"not null"`
	Quantity      uint      `gorm:"not null"`
	SubTotal      uint      `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}
