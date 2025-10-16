package handlers

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/whalelogic/whalesurvey/models"
	"github.com/whalelogic/whalesurvey/views/pages"
)

func HandleSurveyDetailPage(surveyService *models.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return fiber.ErrBadRequest
		}
		survey, err := surveyService.GetSurvey(uint(id))
		if err != nil {
			// Handle survey not found error
			return fiber.ErrNotFound
		}

		return adaptor.HTTPHandler(templ.Handler(pages.SurveyDetailPage(*survey)))(c)
	}
}