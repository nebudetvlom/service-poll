package models

import "gorm.io/gorm"

// Question таблица
type Question struct {
	gorm.Model
	Text           string `gorm:"not null"`
	Type           string `gorm:"not null"` // тип вопроса: "single" или "multiple"
	PollID         uint
	PossibleAnswer []PossibleAnswer
	Answers        []Answer `gorm:"foreignKey:QuestionID"`
}
