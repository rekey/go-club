package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Url      string  `json:"url" gorm:"primaryKey"`
	Name     string  `json:"name"`
	Progress float64 `json:"progress"`
	Err      string  `json:"err"`
	Status   int     `json:"status" gorm:"index"`
}
