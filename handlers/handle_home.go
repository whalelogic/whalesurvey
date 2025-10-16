package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/whalelogic/whalesurvey/views/pages"
	"github.com/a-h/templ"
)

func HandleHomePage(c *fiber.Ctx) error {
	return adaptor.HTTPHandler(templ.Handler(pages.HomePage()))(c)
}
