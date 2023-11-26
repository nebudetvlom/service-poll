package main

import (
	"log"
	"net/http"
	"service-poll/pkg/db"
	"service-poll/pkg/handlers"
	"service-poll/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Создаем таблицы
func Migrate(c *gin.Context) {
	if db.DB.Migrator().HasTable(&models.Poll{}) == false {
		db.DB.AutoMigrate(&models.Poll{})
	}
	if db.DB.Migrator().HasTable(&models.Question{}) == false {
		db.DB.AutoMigrate(&models.Question{})
	}
	if db.DB.Migrator().HasTable(&models.PossibleAnswer{}) == false {
		db.DB.AutoMigrate(&models.PossibleAnswer{})
	}
	if db.DB.Migrator().HasTable(&models.Answer{}) == false {
		db.DB.AutoMigrate(&models.Answer{})
	}
	if db.DB.Migrator().HasTable(&models.AnswerPossibleAnswer{}) == false {
		db.DB.AutoMigrate(&models.AnswerPossibleAnswer{})
	}

	c.JSON(200, gin.H{
		"message": "Migrations applied successfully!",
	})
}

// getAnswerPossibleAnswerIDs возвращает идентификаторы возможных ответов для ответа на вопрос
func getAnswerPossibleAnswerIDs(possibleAnswers []models.PossibleAnswer) []uint {
	var ids []uint
	for _, possibleAnswer := range possibleAnswers {
		ids = append(ids, possibleAnswer.ID)
	}
	return ids
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// инициализация DB
	db.InitDB()

	router := gin.Default()

	// GET /check_alive
	// Проверяет, работает ли сервер.
	router.GET("/migrate", Migrate)

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

	// PATCH /poll/:id/question
	// Изменяет вопрос к опросу по указанному идентификатору.
	router.PATCH("/poll/:id/question/:qid", handlers.UpdateQuestion)

	// DELETE /poll/:id/question
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
