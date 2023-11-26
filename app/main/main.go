package main

import (
	"log"
	"net/http"
	"service-poll/migrations"
	"service-poll/pkg/db"
	"service-poll/pkg/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// инициализация DB
	db.InitDB()

	router := gin.Default()

	// POST /migrate
	// Создаёт базу данных и таблицы + наполняет их фикстурами
	router.POST("/migrate", migrations.Migrate)

	// GET /check_alive
	// Проверяет, работает ли сервер.
	router.GET("/check_alive", CheckAlive)

	// POST /poll/create
	// Создает новый опрос.
	router.POST("/poll/create", handlers.CreatePoll)

	// GET /poll/:id
	// Получает опрос с вопросами по указанному идентификатору.
	router.GET("/poll/:id", handlers.GetPoll)

	// PATCH /poll/:id
	// Изменяет опрос по указанному идентификатору.
	router.PATCH("/poll/:id", handlers.UpdatePoll)

	// DELETE /poll/:id
	// Удаляет опрос по указанному идентификатору.
	router.DELETE("/poll/:id", handlers.DeletePoll)

	// GET /poll/:id/results
	// Получает результаты опроса по указанному идентификатору.
	router.GET("/poll/:id/results", handlers.GetPollResults)

	// POST /poll/:id/question
	// Добавляет вопрос к опросу по указанному идентификатору.
	router.POST("/poll/:id/question", handlers.AddQuestion)

	// PATCH /poll/:id/question/:qid
	// Изменяет вопрос к опросу по указанному идентификатору.
	router.PATCH("/poll/:id/question/:qid", handlers.UpdateQuestion)

	// DELETE /poll/:id/question/:qid
	// Удаляет вопрос к опросу по указанному идентификатору.
	router.DELETE("/poll/:id/question/:qid", handlers.DeleteQuestion)

	// POST /poll/:id/answer
	// Регистрирует ответ на вопрос к опросу по указанному идентификатору.
	router.POST("/poll/:id/answer", handlers.RegisterAnswer)

	router.HandleMethodNotAllowed = true
	router.Run(":5000")
}

// CheckAlive возвращает статус 200, чтобы проверить, работает ли сервер.
func CheckAlive(c *gin.Context) {
	c.Status(http.StatusOK)
}
