package models

import (
	"time"
)

type Exam struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	TotalScore  float64   `json:"total_score"`
	Duration    int       `json:"duration"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Questions   []Question `json:"questions,omitempty"`
}

type ExamRecord struct {
	ID         int       `json:"id"`
	ExamID     int       `json:"exam_id"`
	UserID     int       `json:"user_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Duration   int       `json:"duration"`
	TotalScore float64   `json:"total_score"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Exam       *Exam     `json:"exam,omitempty"`
	User       *User     `json:"user,omitempty"`
	Answers    []ExamAnswer `json:"answers,omitempty"`
}

type ExamAnswer struct {
	ID          int       `json:"id"`
	RecordID    int       `json:"record_id"`
	QuestionID  int       `json:"question_id"`
	UserAnswer  string    `json:"user_answer"`
	Score       float64   `json:"score"`
	IsCorrect   int       `json:"is_correct"`
	CreatedAt   time.Time `json:"created_at"`
	Question    *Question  `json:"question,omitempty"`
	Record      *ExamRecord `json:"record,omitempty"`
}

type ExamCreateRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"omitempty"`
	Subject     string    `json:"subject" binding:"required"`
	Duration    int       `json:"duration" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	QuestionIDs []int     `json:"question_ids" binding:"required"`
}

type ExamGenerateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	Subject     string `json:"subject" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	QuestionCount int   `json:"question_count" binding:"required"`
	Difficulty   string `json:"difficulty" binding:"omitempty"`
}

type ExamAnswerRequest struct {
	QuestionID int    `json:"question_id" binding:"required"`
	UserAnswer string `json:"user_answer" binding:"required"`
}

type ExamSubmitRequest struct {
	RecordID int                `json:"record_id" binding:"required"`
	Answers  []ExamAnswerRequest `json:"answers" binding:"required"`
}
