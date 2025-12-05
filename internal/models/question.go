package models

import (
	"time"
)

type QuestionBank struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Questions   []Question `json:"questions,omitempty"`
}

type Question struct {
	ID          int       `json:"id"`
	BankID      int       `json:"bank_id"`
	Type        string    `json:"type"`
	Content     string    `json:"content"`
	Options     string    `json:"options"`
	Answer      string    `json:"answer"`
	Score       float64   `json:"score"`
	Difficulty  string    `json:"difficulty"`
	Analysis    string    `json:"analysis"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Bank        *QuestionBank `json:"bank,omitempty"`
}

type QuestionCreateRequest struct {
	BankID     int     `json:"bank_id" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Content    string  `json:"content" binding:"required"`
	Options    string  `json:"options" binding:"required"`
	Answer     string  `json:"answer" binding:"required"`
	Score      float64 `json:"score" binding:"required"`
	Difficulty string  `json:"difficulty" binding:"omitempty"`
	Analysis   string  `json:"analysis" binding:"omitempty"`
}

type QuestionBankCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	Subject     string `json:"subject" binding:"required"`
}
