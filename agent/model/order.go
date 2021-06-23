package model

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	Timestamp 	 time.Time	   `json:"timestamp" gorm:"not null"`
	FullPrice    float32       `json:"fullPrice" gorm:"not null"`
	UserId       uint          `json:"userId" gorm:"not null"`
	AgentId      uint          `json:"agentId" gorm:"not null"`
	Items        []Item        `json:"items"`
}

type Item struct {
	gorm.Model
	ProductId       uint       `json:"productId" gorm:"not null"`
	Quantity        uint       `json:"quantity" gorm:"not null"`
	OrderId			uint	   `json:"orderId"`
}
