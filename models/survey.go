package models

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Survey represents a survey with questions and responses
type Survey struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Questions   []Question `json:"questions,omitempty" gorm:"foreignKey:SurveyID"`
	Responses   []Response `json:"responses,omitempty" gorm:"foreignKey:SurveyID"`
}

// Question represents a survey question
type Question struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SurveyID  uint      `json:"survey_id" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"` // text, multiple_choice, checkbox, rating, etc.
	Question  string    `json:"question" gorm:"not null"`
	Required  bool      `json:"required" gorm:"default:false"`
	Order     int       `json:"order" gorm:"not null"`
	Options   []Option  `json:"options,omitempty" gorm:"foreignKey:QuestionID"`
	Answers   []Answer  `json:"answers,omitempty" gorm:"foreignKey:QuestionID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Option represents multiple choice options for questions
type Option struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	QuestionID uint      `json:"question_id" gorm:"not null"`
	Text       string    `json:"text" gorm:"not null"`
	Order      int       `json:"order"`
	CreatedAt  time.Time `json:"created_at"`
}

// Response represents a complete survey response from a user
type Response struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	SurveyID    uint       `json:"survey_id" gorm:"not null"`
	UserID      *string    `json:"user_id,omitempty"` // Optional for anonymous responses
	IPAddress   string     `json:"ip_address"`
	UserAgent   string     `json:"user_agent"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	IsComplete  bool       `json:"is_complete" gorm:"default:false"`
	Answers     []Answer   `json:"answers,omitempty" gorm:"foreignKey:ResponseID"`
	Survey      Survey     `json:"survey,omitempty" gorm:"foreignKey:SurveyID"`
}

// Answer represents an answer to a specific question
type Answer struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ResponseID uint      `json:"response_id" gorm:"not null"`
	QuestionID uint      `json:"question_id" gorm:"not null"`
	AnswerText string    `json:"answer_text"`         // For text answers
	OptionID   *uint     `json:"option_id,omitempty"` // For multiple choice
	Rating     *int      `json:"rating,omitempty"`    // For rating questions
	CreatedAt  time.Time `json:"created_at"`
	Question   Question  `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Option     *Option   `json:"option,omitempty" gorm:"foreignKey:OptionID"`
}

// Database represents the database connection
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&Survey{}, &Question{}, &Option{}, &Response{}, &Answer{})
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

// SurveyStats represents aggregated statistics for a survey
type SurveyStats struct {
	SurveyID       uint    `json:"survey_id"`
	TotalResponses int     `json:"total_responses"`
	CompletionRate float64 `json:"completion_rate"`
	AverageTime    float64 `json:"average_time_seconds"`
}

// QuestionStats represents statistics for a specific question
type QuestionStats struct {
	QuestionID    uint         `json:"question_id"`
	QuestionType  string       `json:"question_type"`
	ResponseCount int          `json:"response_count"`
	OptionStats   []OptionStat `json:"option_stats,omitempty"`
	TextAnswers   []string     `json:"text_answers,omitempty"`
	RatingStats   *RatingStats `json:"rating_stats,omitempty"`
}

// OptionStat represents statistics for a multiple choice option
type OptionStat struct {
	OptionID   uint    `json:"option_id"`
	OptionText string  `json:"option_text"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// RatingStats represents statistics for rating questions
type RatingStats struct {
	Average float64 `json:"average"`
	Min     int     `json:"min"`
	Max     int     `json:"max"`
	Count   int     `json:"count"`
}
