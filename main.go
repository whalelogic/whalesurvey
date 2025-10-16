package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/whalelogic/whalesurvey/handlers"
	"github.com/whalelogic/whalesurvey/models"
)

func main() {
	// Initialize database
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "surveys.db"
	}

	db, err := models.NewDatabase(dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	surveyService := models.NewSurveyService(db)

	// Add some sample data
	if surveys, _ := surveyService.GetAllSurveys(); len(surveys) == 0 {
		log.Println("No surveys found, creating sample data...")
		survey, err := surveyService.CreateSurvey("Favorite Programming Languages", "A survey to find out which programming languages people like the most.")
		if err != nil {
			log.Println("Error creating sample survey:", err)
		}

		questions := []struct {
			Question string
			Options  []string
		}{
			{"What is your favorite systems programming language?", []string{"Rust", "C++", "Go", "C"}},
			{"What is your favorite web backend language?", []string{"Go", "Python", "JavaScript/Node.js", "Rust"}},
			{"What is your favorite web frontend framework?", []string{"React", "Vue", "Svelte", "HTMX"}},
		}

		for _, q := range questions {
			_, err := surveyService.AddQuestion(survey.ID, "multiple_choice", q.Question, true, q.Options)
			if err != nil {
				log.Println("Error adding question:", err)
			}
		}
	}


	app := fiber.New()

	app.Static("/static", "./static")

	app.Get("/", handlers.HandleHomePage)
	app.Get("/surveys", handlers.HandleSurveysListPage(surveyService))
	app.Get("/surveys/:id", handlers.HandleSurveyDetailPage(surveyService))
	app.Post("/surveys/:id/submit", handlers.HandleSurveySubmit(surveyService))
	app.Get("/surveys/:id/results", handlers.HandleSurveyResultsPage(surveyService))

	log.Fatal(app.Listen(":8080"))
}
