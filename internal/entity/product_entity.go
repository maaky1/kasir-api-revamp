package entity

import "time"

type Product struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	CategoryID uint      `gorm:"not null;default:1"`
	Name       string    `gorm:"type:text;not null"`
	Price      int       `gorm:"not null"`
	Stock      int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
