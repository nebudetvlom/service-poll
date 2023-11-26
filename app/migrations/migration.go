package migrations

import (
	"errors"
	"math/rand"
	"service-poll/pkg/db"
	"service-poll/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Создаем таблицы и заполняем их фикстурами
func Migrate(c *gin.Context) {
	if db.DB.Migrator().HasTable(&models.Poll{}) == false {
		db.DB.AutoMigrate(&models.Poll{})
	}

	var firstPoll models.Poll
	if err := db.DB.First(&firstPoll).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		CreatePolls()
	}

	if db.DB.Migrator().HasTable(&models.Question{}) == false {
		db.DB.AutoMigrate(&models.Question{})
	}

	var firstQuestion models.Question
	if err := db.DB.First(&firstQuestion).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		CreateQuestions()
	}

	if db.DB.Migrator().HasTable(&models.AnswerPossibleAnswer{}) == false {
		db.DB.AutoMigrate(&models.AnswerPossibleAnswer{})
	}

	if db.DB.Migrator().HasTable(&models.PossibleAnswer{}) == false {
		db.DB.AutoMigrate(&models.PossibleAnswer{})
	}

	var firstPossibleAnswer models.PossibleAnswer
	if err := db.DB.First(&firstPossibleAnswer).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		CreatePossibleAnswers()
	}

	if db.DB.Migrator().HasTable(&models.Answer{}) == false {
		db.DB.AutoMigrate(&models.Answer{})
	}

	var firstAnswer models.Answer
	var firstAnswerPossibleAnswer models.AnswerPossibleAnswer
	if err := db.DB.First(&firstAnswer).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if errapa := db.DB.First(&firstAnswerPossibleAnswer).Error; errors.Is(errapa, gorm.ErrRecordNotFound) {
			CreateAnswers()
		}
	}

	c.JSON(200, gin.H{
		"message": "Migrations applied successfully!",
	})
}

func createPolls() error {
	polls := []models.Poll{
		{Title: "Poll 1", URL: "poll-1-url"},
		{Title: "Poll 2", URL: "poll-2-url"},
		{Title: "Poll 3", URL: "poll-3-url"},
	}

	for _, poll := range polls {
		if err := db.DB.Create(&poll).Error; err != nil {
			return err
		}
	}

	return nil
}

func CreatePolls() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := createPolls(); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

func createQuestions() error {
	questions := []models.Question{
		{Text: "Question 1 Poll 1", Type: "single", PollID: 1},
		{Text: "Question 2 Poll 1", Type: "multiple", PollID: 1},
		{Text: "Question 3 Poll 1", Type: "single", PollID: 1},
		{Text: "Question 1 Poll 2", Type: "multiple", PollID: 2},
		{Text: "Question 2 Poll 2", Type: "single", PollID: 2},
		{Text: "Question 3 Poll 2", Type: "multiple", PollID: 2},
		{Text: "Question 1 Poll 3", Type: "single", PollID: 3},
		{Text: "Question 2 Poll 3", Type: "single", PollID: 3},
		{Text: "Question 3 Poll 3", Type: "multiple", PollID: 3},
	}

	for _, question := range questions {
		if err := db.DB.Create(&question).Error; err != nil {
			return err
		}
	}

	return nil
}

func CreateQuestions() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := createQuestions(); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

func createPossibleAnswers() error {
	possibleAnswers := []models.PossibleAnswer{
		{Text: "Answer 1 Question 1 Poll 1", QuestionID: 1},
		{Text: "Answer 2 Question 1 Poll 1", QuestionID: 1},
		{Text: "Answer 3 Question 1 Poll 1", QuestionID: 1},
		{Text: "Answer 1 Question 2 Poll 1", QuestionID: 2},
		{Text: "Answer 2 Question 2 Poll 1", QuestionID: 2},
		{Text: "Answer 3 Question 2 Poll 1", QuestionID: 2},
		{Text: "Answer 1 Question 3 Poll 1", QuestionID: 3},
		{Text: "Answer 2 Question 3 Poll 1", QuestionID: 3},
		{Text: "Answer 3 Question 3 Poll 1", QuestionID: 3},
		{Text: "Answer 1 Question 1 Poll 2", QuestionID: 4},
		{Text: "Answer 2 Question 1 Poll 2", QuestionID: 4},
		{Text: "Answer 3 Question 1 Poll 2", QuestionID: 4},
		{Text: "Answer 1 Question 2 Poll 2", QuestionID: 5},
		{Text: "Answer 2 Question 2 Poll 2", QuestionID: 5},
		{Text: "Answer 3 Question 2 Poll 2", QuestionID: 5},
		{Text: "Answer 1 Question 3 Poll 2", QuestionID: 6},
		{Text: "Answer 2 Question 3 Poll 2", QuestionID: 6},
		{Text: "Answer 3 Question 3 Poll 2", QuestionID: 6},
		{Text: "Answer 1 Question 1 Poll 3", QuestionID: 7},
		{Text: "Answer 2 Question 1 Poll 3", QuestionID: 7},
		{Text: "Answer 3 Question 1 Poll 3", QuestionID: 7},
		{Text: "Answer 1 Question 2 Poll 3", QuestionID: 8},
		{Text: "Answer 2 Question 2 Poll 3", QuestionID: 8},
		{Text: "Answer 3 Question 2 Poll 3", QuestionID: 8},
		{Text: "Answer 1 Question 3 Poll 3", QuestionID: 9},
		{Text: "Answer 2 Question 3 Poll 3", QuestionID: 9},
		{Text: "Answer 3 Question 3 Poll 3", QuestionID: 9},
	}

	for _, possibleAnswer := range possibleAnswers {
		if err := db.DB.Create(&possibleAnswer).Error; err != nil {
			return err
		}
	}

	return nil
}

func CreatePossibleAnswers() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := createPossibleAnswers(); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

func createAnswers() error {
	answers := []models.Answer{
		{UserID: 1, AccessToken: "user-1-token", QuestionID: 1, PollID: 1},
		{UserID: 2, AccessToken: "user-2-token", QuestionID: 2, PollID: 1},
		{UserID: 3, AccessToken: "user-3-token", QuestionID: 3, PollID: 1},
		{UserID: 4, AccessToken: "user-4-token", QuestionID: 4, PollID: 2},
		{UserID: 5, AccessToken: "user-5-token", QuestionID: 5, PollID: 2},
		{UserID: 6, AccessToken: "user-6-token", QuestionID: 6, PollID: 2},
		{UserID: 7, AccessToken: "user-7-token", QuestionID: 7, PollID: 3},
		{UserID: 8, AccessToken: "user-8-token", QuestionID: 8, PollID: 3},
		{UserID: 9, AccessToken: "user-9-token", QuestionID: 9, PollID: 3},
		{UserID: 10, AccessToken: "user-10-token", QuestionID: 1, PollID: 1},
		{UserID: 11, AccessToken: "user-11-token", QuestionID: 2, PollID: 1},
		{UserID: 12, AccessToken: "user-12-token", QuestionID: 3, PollID: 1},
		{UserID: 13, AccessToken: "user-13-token", QuestionID: 4, PollID: 2},
		{UserID: 14, AccessToken: "user-14-token", QuestionID: 5, PollID: 2},
		{UserID: 15, AccessToken: "user-15-token", QuestionID: 6, PollID: 2},
		{UserID: 16, AccessToken: "user-16-token", QuestionID: 7, PollID: 3},
		{UserID: 17, AccessToken: "user-17-token", QuestionID: 8, PollID: 3},
		{UserID: 18, AccessToken: "user-18-token", QuestionID: 9, PollID: 3},
	}

	for _, answer := range answers {
		if err := db.DB.Create(&answer).Error; err != nil {
			return err
		}

		// Получаем вопрос из базы данных
		var question models.Question
		if err := db.DB.First(&question, answer.QuestionID).Error; err != nil {
			return err
		}

		// Создаем связи с вариантами ответов через таблицу answer_possible_answers
		var possibleAnswers []models.PossibleAnswer
		if err := db.DB.Where("question_id = ?", answer.QuestionID).Find(&possibleAnswers).Error; err != nil {
			return err
		}

		// Получаем случайную перестановку индексов возможных ответов
		perm := rand.Perm(len(possibleAnswers))

		// Определяем количество случайных ответов в зависимости от типа вопроса
		var randomAnswersCount int
		if question.Type == "single" {
			randomAnswersCount = 1
		} else if question.Type == "multiple" {
			randomAnswersCount = rand.Intn(2) + 2
		}

		// Выбираем случайные ответы и создаем связи
		for i := 0; i < randomAnswersCount && i < len(possibleAnswers); i++ {
			randomAnswer := possibleAnswers[perm[i]]

			answerPossibleAnswer := models.AnswerPossibleAnswer{
				AnswerID:         answer.ID,
				PossibleAnswerID: randomAnswer.ID,
			}

			if err := db.DB.Create(&answerPossibleAnswer).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateAnswers() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := createAnswers(); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}
