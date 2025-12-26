# 🧭 Survey Application 

🎟️ A simple, modular **Go web application** that lets users create, manage, and take surveys.
This project demonstrates clean architecture with handlers, models, templates, and static assets, all connected to a SQLite database.

---

## 🚀 Features

* Create and manage surveys
* Store questions and responses in a SQLite database
* Modular structure with handlers, models, and views
* HTML templates for user interaction
* Static files for styling and client scripts
* Fully scaffolded Go project (ready for extension)

---

## 🧩 Project Structure

```
.
├── handlers/          # HTTP route handlers (controllers)
├── models/            # Data models and database interactions
├── static/            # CSS, JS, and other static files
├── views/             # HTML templates (survey forms, results, etc.)
├── main.go            # Entry point for the application
├── surveys.db         # SQLite database
├── go.mod             # Go module file
└── go.sum             # Dependencies checksum
```

---

## ⚙️ Installation & Setup

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/survey-app.git
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

The app will start on:

```
http://localhost:8080
```

---

## 🎥 Screenshots

<img width="1088" height="832" alt="image" src="https://github.com/user-attachments/assets/1f0ff7d9-9786-4873-a94a-59660d6856a2" />
<img width="1073" height="636" alt="image" src="https://github.com/user-attachments/assets/0b4127db-cf90-4ad0-8e8b-7bc313b0af68" />

---


## 🧠 Technical Overview

* **Language:** Go (1.21+)
* **Database:** SQLite (`surveys.db`)
* **Routing:** Standard `net/http` or `chi` (depending on your implementation)
* **Templates:** Go’s `html/template` package
* **Static assets:** Served from `/static`

---

## 🧱 Flow

1. **User visits homepage** — chooses a survey or creates one.
2. **Handler** processes requests (`handlers/survey.go`).
3. **Model** queries or inserts survey data (`models/survey.go`).
4. **View** renders a dynamic HTML template (`views/survey.html`).
5. **Static files** handle front-end appearance and interactivity.

---

## 📈 Future Enhancements

* Add user authentication and session management
* Support multiple question types (MCQs, ratings, etc.)
* Store results and analytics
* REST API endpoints for survey management
* Add Dockerfile and CI/CD pipeline

---

## 🧑‍💻 Author

**Keith Thomson**
Computer Science Student • Army Veteran • Web Developer
💡 Focused on Go, Rust, and system design for intelligent web applications.

---

## 🪪 License

MIT License — See [LICENSE](LICENSE) file for details.
