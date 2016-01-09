package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

// Question contains all informations about a question
type Question struct {
	gorm.Model
	UserID        uint
	Sentence      string
	RightAnswerID uint
	StartedAt     time.Time
}

// getCurrentQuestion returns current question
func getCurrentQuestion(w http.ResponseWriter, r *http.Request) {
	question, err := GetCurrentQuestion()
	if err != nil {
		renderJSON(w, http.StatusNotFound, errCurQuestionNotFound)
		return
	}
	renderJSON(w, http.StatusOK, question)
}

// GetCurrentQuestion returns the current question
func GetCurrentQuestion() (*Question, error) {
	return getCurrentQuestionWithTX(&db)
}

func getCurrentQuestionWithTX(tx *gorm.DB) (*Question, error) {
	question := &Question{}
	err := tx.Order("started_at").Last(question).Error
	return question, err
}

func nextQuestion() error {
	tx := db.Begin()
	defer tx.Commit()
	if err := updateUsersPoints(tx); err != nil {
		return err
	}
	return randomizeQuestion(tx)
}

func updateUsersPoints(tx *gorm.DB) error {
	q, err := getCurrentQuestionWithTX(tx)
	if err != nil {
		return nil
	}
	return tx.Exec("UPDATE `users` JOIN `answer_entries` ON users.id = answer_entries.user_id SET users.points = users.points + 1 WHERE answer_entries.question_id = ? AND answer_entries.answer_id = ?", q.ID, q.RightAnswerID).Error
}

func randomizeQuestion(tx *gorm.DB) error {
	var questions []Question
	if err := tx.Where("started_at = ?", 0).Find(&questions).Error; err != nil {
		return err
	}
	question := &questions[rand.Intn(len(questions))]
	return tx.Model(question).UpdateColumn("started_at", time.Now()).Error
}
