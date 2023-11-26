package models

import "gorm.io/gorm"

// Poll таблица
type Poll struct {
	gorm.Model
	Title     string `gorm:"not null"`
	URL       string `gorm:"not null"`
	Questions []Question
	Answers   []Answer `gorm:"foreignKey:PollID"`
}
