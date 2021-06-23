package model

import (
	"gorm.io/gorm"
)

type Statistics struct {
	gorm.Model
	NumOfLikes        		 uint        `json:"numOfLikes"`
	NumOfDislikes     		 uint        `json:"numOfDislikes"`
	NumOfVisits       		 uint        `json:"numOfVisits"`
	NumOfComments     		 uint        `json:"numOfComments"`
}
