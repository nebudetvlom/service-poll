package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Answer таблица
type Answer struct {
	gorm.Model
	UserID          uint
	AccessToken     string
	QuestionID      uint
	PollID          uint
	PossibleAnswers []PossibleAnswer `gorm:"many2many:answer_possible_answers"`
}

// Poll таблица
type Poll struct {
	gorm.Model
	Title     string `gorm:"not null"`
	URL       string `gorm:"not null"`
	Questions []Question
	Answers   []Answer `gorm:"foreignKey:PollID"`
}

// Question таблица
type Question struct {
	gorm.Model
	Text           string `gorm:"not null"`
	Type           string `gorm:"not null"` // тип вопроса: "single" или "multiple"
	PollID         uint
	PossibleAnswer []PossibleAnswer
	Answers        []Answer `gorm:"foreignKey:QuestionID"`
}

// PossibleAnswer таблица
type PossibleAnswer struct {
	gorm.Model
	Text       string `gorm:"not null"`
	QuestionID uint
	Answers    []Answer `gorm:"many2many:answer_possible_answers"`
}

var DB *gorm.DB // глобальная переменная для доступа к базе данных
var err error

// Создаем таблицы
func Migrate(c *gin.Context) {
	if DB.Migrator().HasTable(&Poll{}) == false {
		DB.AutoMigrate(&Poll{})
	}
	if DB.Migrator().HasTable(&Question{}) == false {
		DB.AutoMigrate(&Question{})
	}
	if DB.Migrator().HasTable(&PossibleAnswer{}) == false {
		DB.AutoMigrate(&PossibleAnswer{})
	}
	if DB.Migrator().HasTable(&Answer{}) == false {
		DB.AutoMigrate(&Answer{})
	}

	c.JSON(200, gin.H{
		"message": "Migrations applied successfully!",
	})
}

func initDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbTimezone := os.Getenv("DB_TIMEZONE")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUser, dbPassword, dbTimezone)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	var dbNameOnServer string
	queryResult := DB.Raw("SELECT datname FROM pg_database WHERE datname = ?", dbName).Scan(&dbNameOnServer)
	if queryResult.Error != nil {
		log.Panicf("Error executing query: %v", queryResult.Error)
	} else if queryResult.RowsAffected == 0 {
		createDatabaseCommand := fmt.Sprintf("CREATE DATABASE \"%s\"", dbName)
		DB.Exec(createDatabaseCommand)
	}

	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbTimezone)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// инициализация DB
	initDB()

	router := gin.Default()

	// router.Use(DisableOptionsMiddleware())

	// router.HandleMethodNotAllowed = true
	// router.NoRoute(NotFoundMiddleware)
	// router.NoMethod(WrongMethodMiddleware)

	// GET /check_alive
	// Проверяет, работает ли сервер.
	router.GET("/migrate", Migrate)

	// GET /check_alive
	// Проверяет, работает ли сервер.
	router.GET("/check_alive", CheckAlive)

	// POST /poll/create
	// Создает новый опрос.
	router.POST("/poll/create", CreatePoll)

	// GET /poll/:id
	// Получает опрос с вопросами по указанному идентификатору.
	router.GET("/poll/:id", GetPoll)

	// PATCH /poll/:id
	// Изменяет опрос по указанному идентификатору.
	router.PATCH("/poll/:id", UpdatePoll)

	// DELETE /poll/:id
	// Удаляет опрос по указанному идентификатору.
	router.DELETE("/poll/:id", DeletePoll)

	// GET /poll/:id/results
	// Получает результаты опроса по указанному идентификатору.
	router.GET("/poll/:id/results", GetPollResults)

	// POST /poll/:id/question
	// Добавляет вопрос к опросу по указанному идентификатору.
	router.POST("/poll/:id/question", AddQuestion)

	// PATCH /poll/:id/question
	// Изменяет вопрос к опросу по указанному идентификатору.
	router.PATCH("/poll/:id/question", UpdateQuestion)

	// DELETE /poll/:id/question
	// Удаляет вопрос к опросу по указанному идентификатору.
	router.DELETE("/poll/:id/question", DeleteQuestion)

	// POST /poll/:id/answer
	// Регистрирует ответ на вопрос к опросу по указанному идентификатору.
	router.POST("/poll/:id/answer", RegisterAnswer)

	router.HandleMethodNotAllowed = true
	router.Run(":5000")
}

// CheckAlive возвращает статус 200, чтобы проверить, работает ли сервер.
func CheckAlive(c *gin.Context) {
	c.Status(http.StatusOK)
}

// GetPoll возвращает опрос с вопросами по указанному идентификатору.
func GetPoll(c *gin.Context) {
	// ... логика обработки запроса ...
}

// CreatePoll создает новый опрос.
func CreatePoll(c *gin.Context) {
	var pollData struct {
		Title string `json:"title" binding:"required"`
		URL   string `json:"url" binding:"required"`
	}

	// Извлекаем данные из тела запроса
	if err := c.ShouldBindJSON(&pollData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем новый объект Poll
	newPoll := Poll{
		Title: pollData.Title,
		URL:   pollData.URL,
	}

	// Сохраняем новый объект в базе данных
	if err := DB.Create(&newPoll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
		return
	}

	c.JSON(http.StatusOK, newPoll)
}

// UpdatePoll изменяет опрос по указанному идентификатору.
func UpdatePoll(c *gin.Context) {
	// ... логика обработки запроса ...
}

// DeletePoll удаляет опрос по указанному идентификатору.
func DeletePoll(c *gin.Context) {
	// ... логика обработки запроса ...
}

// GetPollResults возвращает результаты опроса по указанному идентификатору.
func GetPollResults(c *gin.Context) {
	// ... логика обработки запроса ...
}

func AddQuestion(c *gin.Context) {
	// Получаем идентификатор опроса из параметра запроса
	pollID := c.Param("id")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll Poll
	if err := DB.First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Извлекаем данные вопроса из тела запроса
	var questionData struct {
		Text    string   `json:"text" binding:"required"`
		Type    string   `json:"type" binding:"required,oneof=single multiple"`
		Answers []string `json:"answers" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&questionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка значения Type
	if questionData.Type != "single" && questionData.Type != "multiple" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be either 'single' or 'multiple'"})
		return
	}

	// Создаем новый объект Question
	newQuestion := Question{
		Text:   questionData.Text,
		Type:   questionData.Type,
		PollID: existingPoll.ID,
	}

	// Сохраняем новый вопрос в базе данных
	if err := DB.Create(&newQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	// Сохраняем варианты ответов
	var createdAnswers []PossibleAnswer
	for _, answerText := range questionData.Answers {
		possibleAnswer := PossibleAnswer{
			Text:       answerText,
			QuestionID: newQuestion.ID,
		}

		if err := DB.Create(&possibleAnswer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create possible answer"})
			return
		}
		createdAnswers = append(createdAnswers, possibleAnswer)
	}

	// newQuestion.Answers = createdAnswers
	c.JSON(http.StatusOK, newQuestion)
}

// UpdateQuestion изменяет вопрос к опросу по указанному идентификатору.
func UpdateQuestion(c *gin.Context) {
	// ... логика обработки запроса ...
}

// DeleteQuestion удаляет вопрос к опросу по указанному идентификатору.
func DeleteQuestion(c *gin.Context) {
	// ... логика обработки запроса ...
}

func RegisterAnswer(c *gin.Context) {
	// Получаем идентификатор опроса из параметра запроса
	pollID := c.Param("id")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll Poll
	if err := DB.First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	type AnswerRequest struct {
		AccessToken string `json:"access_token" binding:"required"`
		UserID      uint   `json:"user_id" binding:"required"`
		Results     []struct {
			QuestionID uint   `json:"question_id" binding:"required"`
			Answers    []uint `json:"answers" binding:"required,min=1"`
		} `json:"results" binding:"required"`
	}

	var answerRequest AnswerRequest

	if err := c.ShouldBindJSON(&answerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ваши проверки и логика обработки ответа

	for _, result := range answerRequest.Results {
		// Проверяем, существует ли вопрос с указанным ID в опросе
		var existingQuestion Question
		if err := DB.First(&existingQuestion, result.QuestionID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question not found"})
			return
		}

		// Проверяем тип вопроса
		if existingQuestion.Type == "single" && len(result.Answers) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "For single type question, exactly one answer is required"})
			return
		} else if existingQuestion.Type == "multiple" && len(result.Answers) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "For multiple type question, at least two answers are required"})
			return
		}

		// Собираем все возможные ответы для данного вопроса
		var possibleAnswers []PossibleAnswer
		for _, answerID := range result.Answers {
			var existingPossibleAnswer PossibleAnswer
			if err := DB.Where("question_id = ?", result.QuestionID).First(&existingPossibleAnswer, answerID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Possible Answer (%d) not found for the specified question (%d)", result.QuestionID, answerID)})
				return
			}
			possibleAnswers = append(possibleAnswers, existingPossibleAnswer)
		}

		pollIDInt, _ := strconv.Atoi(pollID)

		// Создаем новый объект Answer и связываем с вариантами ответов
		newAnswer := Answer{
			UserID:          answerRequest.UserID,
			AccessToken:     answerRequest.AccessToken,
			QuestionID:      result.QuestionID,
			PollID:          uint(pollIDInt),
			PossibleAnswers: possibleAnswers,
		}

		if err := DB.Create(&newAnswer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register answer"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Answers registered successfully!"})
}
