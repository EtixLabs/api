package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	commandTVFunc = map[string]func(*SlackCommandRequest, *User) *SlackCommandResponse{
		"help":     slackCommandTVHelp,
		"question": slackCommandTVQuestion,
		"answer":   slackCommandTVAnswer,
		"status":   slackCommandTVStatus,
	}
	commandTVUsage = "Valid commands: help, question, answer, status."
)

// SlackCommandRequest contains the data of slack command request.
type SlackCommandRequest struct {
	Token       string `schema:"token"`
	UserID      string `schema:"user_id"`
	Command     string `schema:"command"`
	Text        string `schema:"text"`
	ResponseURL string `schema:"response_url"`
}

// SlackCommandResponse contains the data of slack command request.
type SlackCommandResponse struct {
	Text string `json:"text"`
}

func slackCommandTV(w http.ResponseWriter, r *http.Request) {
	var req SlackCommandRequest
	if err := decodeRequestForm(r, &req); err != nil {
		renderJSON(w, http.StatusBadRequest, Error{err.Error()})
		return
	}
	if req.Token != slackOutgoingToken {
		renderJSON(w, http.StatusBadRequest, errInvalidToken)
		return
	}
	user, err := GetUserBySlackID(req.UserID)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, Error{err.Error()})
		return
	}
	resp, err := getCommandTVResponse(&req, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getCommandTVResponse(req *SlackCommandRequest, user *User) (*http.Response, error) {
	cmdStr := req.Text
	if i := strings.IndexRune(req.Text, ' '); i != -1 {
		cmdStr = cmdStr[:i]
	}
	var slackResp *SlackCommandResponse
	if cmdFunc, ok := commandTVFunc[cmdStr]; ok {
		slackResp = cmdFunc(req, user)
	} else {
		text := fmt.Sprintf("Invalid command %q.\n%s", cmdStr, commandTVUsage)
		slackResp = &SlackCommandResponse{Text: text}
	}
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(slackResp); err != nil {
		return nil, err
	}
	return http.Post(req.ResponseURL, ContentJSON, body)
}

func slackCommandTVHelp(req *SlackCommandRequest, user *User) *SlackCommandResponse {
	return &SlackCommandResponse{Text: commandTVUsage}
}

func slackCommandTVQuestion(req *SlackCommandRequest, user *User) *SlackCommandResponse {
	return &SlackCommandResponse{Text: commandTVUsage}
}

func slackCommandTVAnswer(req *SlackCommandRequest, user *User) *SlackCommandResponse {
	resp := &SlackCommandResponse{}
	question, err := GetCurrentQuestion()
	if err != nil {
		resp.Text = fmt.Sprintf("Error: Can't get current question: %v", err)
		return resp
	}
	answers, err := GetAnswersByQuestionID(question.ID)
	if err != nil {
		resp.Text = fmt.Sprintf("Error: Can't get answers: %v", err)
		return resp
	}
	answerIndex, _ := strconv.Atoi(req.Text[len("answer "):])
	if answerIndex <= 0 || answerIndex > len(answers) {
		resp.Text = fmt.Sprintf("Invalid answer index.\nThere is %d possible answers.\nSee help and status for more details", len(answers))
		return resp
	}
	answer := answers[answerIndex]
	answerEntry := AnswerEntry{UserID: user.ID, QuestionID: question.ID, AnswerID: answer.ID}
	if err := InsertOrUpdateDB(answerEntry, answerEntry); err != nil {
		resp.Text = fmt.Sprintf("Error: Can't add your answers: %v", err)
		return resp
	}
	resp.Text = fmt.Sprintf("Answer Added.\n%s %s", question.Sentence, answer.Sentence)
	return resp
}

func slackCommandTVStatus(req *SlackCommandRequest, user *User) *SlackCommandResponse {
	return &SlackCommandResponse{Text: commandTVUsage}
}
