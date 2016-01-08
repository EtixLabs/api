package main

import (
	"github.com/jinzhu/gorm"
)

// Answer contains all informations about an Answer
type Answer struct {
	gorm.Model
	QuestionID uint
	Sentence   string
}