# 🧭 Survey Application

🎟️ A simple, modular **Go web application** that lets users create, manage, and analyze surveys.
This project demonstrates clean architecture using handlers, models, templates, and static assets, all backed by a SQLite database.

---

## 🚀 Features

* **Survey Management:** Create, read, update, and delete (CRUD) surveys.
* **Dynamic Questions:** Support for text answers, multiple choice, and ratings.
* **Data Persistence:** Robust storage using SQLite.
* **Analytics:** Built-in calculation of survey statistics (response counts, rating averages).
* **Modular Design:** Clean separation of concerns (Handlers, Services, Models, Views).
* **Frontend:** HTML templates with static CSS/JS support.

## ⚙️ Technical Overview

* **Language:** Go (1.21+)
* **Database:** SQLite (`surveys.db`)
* **Routing:** Standard `net/http` (compatible with `chi` or `mux`)
* **Templating:** Go’s `html/template` engine
* **Architecture:** Service-Repository pattern

---

## 🧩 Project Structure

```bash
.
├── handlers/          # HTTP route handlers (controllers)
├── models/            # Data structs and database logic
├── services/          # Business logic (SurveyService)
├── static/            # CSS, JS, and images
├── views/             # HTML templates (forms, results, dashboard)
├── main.go            # Application entry point
├── surveys.db         # SQLite database file
├── go.mod             # Go module definition
└── go.sum             # Checksums
```

## 🏗️ Architecture & Core Logic

The application relies on strong typing and a centralized service layer to handle business logic.

<details>
<summary><b>Click to view Core Types & API Signatures</b></summary>

### Key Data Models
The application uses strict struct definitions for database mapping and JSON responses.

```go
// Example: Input structure for submitting answers
type AnswerInput struct {
    QuestionID uint   `json:"question_id"`
    AnswerText string `json:"answer_text,omitempty"`
    OptionID   *uint  `json:"option_id,omitempty"`
    Rating     *int   `json:"rating,omitempty"`
}

// Service Layer wrapper
type SurveyService struct {
    db *Database
}
```

**Other Core Types:**
`Survey`, `Question`, `Answer`, `Response`, `Option`, `Database`
**Analytics Types:**
`SurveyStats`, `QuestionStats`, `OptionStat`, `RatingStats`

### Service API
The `SurveyService` acts as the bridge between the HTTP handlers and the Database.

**Survey Management**
```go
func NewSurveyService(db *Database) *SurveyService
func (s *SurveyService) CreateSurvey(title, description string) (*Survey, error) 
func (s *SurveyService) GetSurvey(id uint) (*Survey, error)
func (s *SurveyService) UpdateSurvey(id uint, title, description string, isActive bool) error
func (s *SurveyService) GetAllSurveys() ([]Survey, error)
func (s *SurveyService) DeleteSurvey(id uint) error
```

**Questions & Responses**
```go
func (s *SurveyService) AddQuestion(surveyID uint, questionType, question string, required bool, options []string) (*Question, error)
func (s *SurveyService) SubmitResponse(surveyID uint, userID *string, ipAddress, userAgent string, answers []AnswerInput) (*Response, error)
```

**Analytics & Stats**
```go
func (s *SurveyService) GetSurveyStats(surveyID uint) (*SurveyStats, error)
func (s *SurveyService) GetQuestionStats(surveyID uint) ([]QuestionStats, error)
func (s *SurveyService) getOptionStats(questionID uint) []OptionStat
func (s *SurveyService) getRatingStats(questionID uint) *RatingStats 
func (s *SurveyService) getTextAnswers(questionID uint, limit int) []string
```

</details>

---

## ⚡ Installation & Setup

### 1. Clone the repository

```bash
git clone [https://github.com/](https://github.com/)<your-username>/survey-app.git
cd survey-app
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the application

```bash
go run main.go
```

The app will start on: **`http://localhost:8080`**

---

## 🎥 Screenshots

| Dashboard View | Survey View |
|:---:|:---:|
| <img src="https://github.com/user-attachments/assets/1f0ff7d9-9786-4873-a94a-59660d6856a2" alt="Dashboard" width="100%"> | <img src="https://github.com/user-attachments/assets/0b4127db-cf90-4ad0-8e8b-7bc313b0af68" alt="Survey Form" width="100%"> |

---

## 📈 Future Roadmap

* [ ] User authentication (OAuth/Session auth)
* [ ] Export results to CSV/PDF
* [ ] Admin dashboard for advanced analytics
* [ ] Docker support & CI/CD pipeline
* [ ] REST API documentation (Swagger/OpenAPI)

---

## 🧑‍💻 Author

**Keith Thomson**
*Computer Science Student • Army Veteran • Cloud Developer*
💡 Focused on Go, Rust, and system design for data applications.

---

## 🪪 License

MIT License — See [LICENSE](LICENSE) file for details.
