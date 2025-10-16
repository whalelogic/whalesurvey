package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// SurveyService provides methods for survey operations
type SurveyService struct {
	db *Database
}

// NewSurveyService creates a new survey service
func NewSurveyService(db *Database) *SurveyService {
	return &SurveyService{db: db}
}

// CreateSurvey creates a new survey
func (s *SurveyService) CreateSurvey(title, description string) (*Survey, error) {
	survey := &Survey{
		Title:       title,
		Description: description,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(survey).Error; err != nil {
		return nil, err
	}

	return survey, nil
}

// GetSurvey retrieves a survey by ID with questions and options
func (s *SurveyService) GetSurvey(id uint) (*Survey, error) {
	var survey Survey
	err := s.db.Preload("Questions.Options").Where("id = ?", id).First(&survey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("survey not found")
		}
		return nil, err
	}
	return &survey, nil
}

// GetAllSurveys retrieves all surveys
func (s *SurveyService) GetAllSurveys() ([]Survey, error) {
	var surveys []Survey
	err := s.db.Preload("Questions").Order("created_at desc").Find(&surveys).Error
	return surveys, err
}

// UpdateSurvey updates a survey
func (s *SurveyService) UpdateSurvey(id uint, title, description string, isActive bool) error {
	updates := map[string]interface{}{
		"title":       title,
		"description": description,
		"is_active":   isActive,
		"updated_at":  time.Now(),
	}

	return s.db.Model(&Survey{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteSurvey deletes a survey and all related data
func (s *SurveyService) DeleteSurvey(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete answers first
		if err := tx.Where("response_id IN (SELECT id FROM responses WHERE survey_id = ?)", id).Delete(&Answer{}).Error; err != nil {
			return err
		}

		// Delete responses
		if err := tx.Where("survey_id = ?", id).Delete(&Response{}).Error; err != nil {
			return err
		}

		// Delete options
		if err := tx.Where("question_id IN (SELECT id FROM questions WHERE survey_id = ?)", id).Delete(&Option{}).Error; err != nil {
			return err
		}

		// Delete questions
		if err := tx.Where("survey_id = ?", id).Delete(&Question{}).Error; err != nil {
			return err
		}

		// Delete survey
		return tx.Delete(&Survey{}, id).Error
	})
}

// AddQuestion adds a question to a survey
func (s *SurveyService) AddQuestion(surveyID uint, questionType, question string, required bool, options []string) (*Question, error) {
	// Get the highest order for this survey
	var maxOrder int
	s.db.Model(&Question{}).Where("survey_id = ?", surveyID).Select("COALESCE(MAX(\"order\"), 0)").Scan(&maxOrder)

	q := &Question{
		SurveyID:  surveyID,
		Type:      questionType,
		Question:  question,
		Required:  required,
		Order:     maxOrder + 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(q).Error; err != nil {
		return nil, err
	}

	// Add options if provided
	if len(options) > 0 {
		for i, optionText := range options {
			option := &Option{
				QuestionID: q.ID,
				Text:       optionText,
				Order:      i + 1,
				CreatedAt:  time.Now(),
			}
			if err := s.db.Create(option).Error; err != nil {
				return nil, err
			}
		}
	}

	return q, nil
}

// SubmitResponse submits a complete survey response
func (s *SurveyService) SubmitResponse(surveyID uint, userID *string, ipAddress, userAgent string, answers []AnswerInput) (*Response, error) {
	var response *Response

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create response
		response = &Response{
			SurveyID:    surveyID,
			UserID:      userID,
			IPAddress:   ipAddress,
			UserAgent:   userAgent,
			StartedAt:   time.Now(),
			CompletedAt: &time.Time{},
			IsComplete:  true,
		}
		*response.CompletedAt = time.Now()

		if err := tx.Create(response).Error; err != nil {
			return err
		}

		// Create answers
		for _, answerInput := range answers {
			answer := &Answer{
				ResponseID: response.ID,
				QuestionID: answerInput.QuestionID,
				AnswerText: answerInput.AnswerText,
				OptionID:   answerInput.OptionID,
				Rating:     answerInput.Rating,
				CreatedAt:  time.Now(),
			}
			if err := tx.Create(answer).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// AnswerInput represents input for an answer
type AnswerInput struct {
	QuestionID uint   `json:"question_id"`
	AnswerText string `json:"answer_text,omitempty"`
	OptionID   *uint  `json:"option_id,omitempty"`
	Rating     *int   `json:"rating,omitempty"`
}

// GetSurveyStats gets statistics for a survey
func (s *SurveyService) GetSurveyStats(surveyID uint) (*SurveyStats, error) {
	var stats SurveyStats
	stats.SurveyID = surveyID

	// Total responses
	var totalResponses int64
	s.db.Model(&Response{}).Where("survey_id = ? AND is_complete = ?", surveyID, true).Count(&totalResponses)
	stats.TotalResponses = int(totalResponses)

	// Completion rate (simplified - assuming all started responses are counted)
	var totalStarted int64
	s.db.Model(&Response{}).Where("survey_id = ?", surveyID).Count(&totalStarted)

	if totalStarted > 0 {
		stats.CompletionRate = float64(totalResponses) / float64(totalStarted) * 100
	}

	// Average completion time
	var avgTime float64
	s.db.Model(&Response{}).
		Where("survey_id = ? AND completed_at IS NOT NULL", surveyID).
		Select("AVG(julianday(completed_at) - julianday(started_at)) * 86400").
		Scan(&avgTime)
	stats.AverageTime = avgTime

	return &stats, nil
}

// GetQuestionStats gets statistics for all questions in a survey
func (s *SurveyService) GetQuestionStats(surveyID uint) ([]QuestionStats, error) {
	var questions []Question
	err := s.db.Where("survey_id = ?", surveyID).Preload("Options").Find(&questions).Error
	if err != nil {
		return nil, err
	}

	var stats []QuestionStats
	for _, question := range questions {
		stat := QuestionStats{
			QuestionID:   question.ID,
			QuestionType: question.Type,
		}

		// Count responses for this question
		var responseCount int64
		s.db.Model(&Answer{}).
			Joins("JOIN responses ON answers.response_id = responses.id").
			Where("answers.question_id = ? AND responses.is_complete = ?", question.ID, true).
			Count(&responseCount)
		stat.ResponseCount = int(responseCount)

		switch question.Type {
		case "multiple_choice", "checkbox":
			stat.OptionStats = s.getOptionStats(question.ID)
		case "rating":
			stat.RatingStats = s.getRatingStats(question.ID)
		case "text":
			stat.TextAnswers = s.getTextAnswers(question.ID, 10) // Get last 10 answers
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func (s *SurveyService) getOptionStats(questionID uint) []OptionStat {
	var options []Option
	s.db.Where("question_id = ?", questionID).Find(&options)

	var stats []OptionStat
	for _, option := range options {
		var count int64
		s.db.Model(&Answer{}).
			Joins("JOIN responses ON answers.response_id = responses.id").
			Where("answers.option_id = ? AND responses.is_complete = ?", option.ID, true).
			Count(&count)

		stat := OptionStat{
			OptionID:   option.ID,
			OptionText: option.Text,
			Count:      int(count),
		}
		stats = append(stats, stat)
	}

	// Calculate percentages
	total := 0
	for _, stat := range stats {
		total += stat.Count
	}
	if total > 0 {
		for i := range stats {
			stats[i].Percentage = float64(stats[i].Count) / float64(total) * 100
		}
	}

	return stats
}

func (s *SurveyService) getRatingStats(questionID uint) *RatingStats {
	var ratings []int
	s.db.Model(&Answer{}).
		Joins("JOIN responses ON answers.response_id = responses.id").
		Where("answers.question_id = ? AND responses.is_complete = ? AND answers.rating IS NOT NULL", questionID, true).
		Pluck("answers.rating", &ratings)

	if len(ratings) == 0 {
		return nil
	}

	sum := 0
	min := ratings[0]
	max := ratings[0]
	for _, rating := range ratings {
		sum += rating
		if rating < min {
			min = rating
		}
		if rating > max {
			max = rating
		}
	}

	return &RatingStats{
		Average: float64(sum) / float64(len(ratings)),
		Min:     min,
		Max:     max,
		Count:   len(ratings),
	}
}

func (s *SurveyService) getTextAnswers(questionID uint, limit int) []string {
	var answers []string
	s.db.Model(&Answer{}).
		Joins("JOIN responses ON answers.response_id = responses.id").
		Where("answers.question_id = ? AND responses.is_complete = ? AND answers.answer_text != ''", questionID, true).
		Order("answers.created_at DESC").
		Limit(limit).
		Pluck("answers.answer_text", &answers)

	return answers
}
