package entity

import "time"

type Category struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:text;not null;unique"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
