package model

import (
	"gorm.io/gorm"
	"time"
)

type Timestamp struct {
	gorm.Model
	Time []time.Time `json:"timestamp"`
}