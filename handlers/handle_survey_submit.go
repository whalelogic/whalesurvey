package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/whalelogic/whalesurvey/models"
)

func HandleSurveySubmit(surveyService *models.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return fiber.ErrBadRequest
		}

		survey, err := surveyService.GetSurvey(uint(id))
		if err != nil {
			return fiber.ErrNotFound
		}

		var answers []models.AnswerInput
		for _, question := range survey.Questions {
			key := fmt.Sprintf("question_%d", question.ID)
			value := c.FormValue(key)
			if value != "" {
				answers = append(answers, models.AnswerInput{
					QuestionID: question.ID,
					AnswerText: value,
				})
			}
		}

		_, err = surveyService.SubmitResponse(uint(id), nil, c.IP(), c.Get("User-Agent"), answers)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.Redirect(fmt.Sprintf("/surveys/%d/results", id))
	}
}