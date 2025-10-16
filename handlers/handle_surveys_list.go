package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/whalelogic/whalesurvey/models"
	"github.com/whalelogic/whalesurvey/views/pages"
	"github.com/a-h/templ"
)

func HandleSurveysListPage(surveyService *models.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		surveys, err := surveyService.GetAllSurveys()
		if err != nil {
			// In a real app, you'd want to log this error
			return fiber.ErrInternalServerError
		}
		return adaptor.HTTPHandler(templ.Handler(pages.SurveysListPage(surveys)))(c)
	}
}