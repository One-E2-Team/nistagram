package model

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	ID              uint 			`json:"id" gorm:"primaryKey"`
	Name            string          `json:"name" gorm:"not null"`
	PicturePath     string          `json:"picturePath" gorm:"not null"`
}

type AgentProduct struct {
	gorm.Model
	UserID          uint          `json:"userId" gorm:"not null"`
	ProductID       uint          `json:"productId" gorm:"not null"`
	Quantity        uint          `json:"quantity" gorm:"not null"`
	PricePerItem    float32       `json:"pricePerItem" gorm:"not null"`
	ValidFrom       time.Time     `json:"validFrom" gorm:"not null"`
	IsValid         bool          `json:"isValid" gorm:"not null"`
}