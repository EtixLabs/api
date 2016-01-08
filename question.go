package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Question contains all informations about a question
type Question struct {
	gorm.Model
	UserID        uint
	Sentence      string
	RightAnswerID uint `sql:"unique"`
	StartedAt     time.Time
}
