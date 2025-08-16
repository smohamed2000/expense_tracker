package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID   uint    `json:"user_id" gorm:"not null"`
	Type     string  `json:"type" gorm:"not null"` // income or expense
	Category string  `json:"category" gorm:"not null"`
	Amount   float64 `json:"amount" gorm:"not null"`
	Note     string  `json:"note" gorm:"not null"`
}
