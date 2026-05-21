package models

import "time"

type User struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	TelegramChatID int64     `json:"telegram_chat_id,omitempty"`
	ActivationCode string    `json:"activation_code,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	PurgeAt        time.Time `json:"purge_at"`
	Active         bool      `json:"active"`
}

type StudySession struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Subject   string    `json:"subject"`
	Topic     string    `json:"topic"`
	StudyType string    `json:"study_type"`
	Stars     int       `json:"stars,omitempty"`
	Correct   int       `json:"correct,omitempty"`
	Wrong     int       `json:"wrong,omitempty"`
	Net       float64   `json:"net,omitempty"`
	StudyDate string    `json:"study_date"`
	CreatedAt time.Time `json:"created_at"`
}

type Todo struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Subject   string    `json:"subject"`
	Topic     string    `json:"topic"`
	TodoType  string    `json:"todo_type"`
	DueDate   string    `json:"due_date"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}
