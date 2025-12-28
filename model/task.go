package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Url    string `json:"url" gorm:"primaryKey"`
	Status string `json:"status" gorm:"index"`
}
