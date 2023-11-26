package models

import "gorm.io/gorm"

// AnswerPossibleAnswer представляет связь между ответами и возможными ответами
type AnswerPossibleAnswer struct {
	gorm.Model
	AnswerID         uint
	PossibleAnswerID uint
}

// Answer таблица
type Answer struct {
	gorm.Model
	UserID          uint
	AccessToken     string
	QuestionID      uint
	PollID          uint
	PossibleAnswers []PossibleAnswer `gorm:"many2many:answer_possible_answers"`
}

// PossibleAnswer таблица
type PossibleAnswer struct {
	gorm.Model
	Text       string `gorm:"not null"`
	QuestionID uint
	Answers    []Answer `gorm:"many2many:answer_possible_answers"`
}
