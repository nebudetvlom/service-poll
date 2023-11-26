package handlers

import (
	"fmt"
	"net/http"
	"service-poll/pkg/db"
	"service-poll/pkg/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterAnswer Регистрирует ответ на вопрос к опросу по указанному идентификатору.
func RegisterAnswer(c *gin.Context) {
	// Получаем идентификатор опроса из параметра запроса
	pollID := c.Param("id")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll models.Poll
	if err := db.DB.First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// AnswerResponse структура для формирования ответа на регистрацию ответов
	type AnswerResponse struct {
		AccessToken string `json:"access_token"`
		PollID      uint   `json:"poll_id"`
		UserID      uint   `json:"user_id"`
		Results     []struct {
			QuestionID uint   `json:"question_id"`
			Answers    []uint `json:"answers"`
		} `json:"results"`
	}

	// AnswerRequest структура для разбора запроса
	type AnswerRequest struct {
		Attributes struct {
			AccessToken string `json:"access_token" binding:"required"`
			UserID      uint   `json:"user_id" binding:"required"`
			Username    string `json:"username"`
			Email       string `json:"email"`
			UserData    struct {
				Reason string `json:"reason"`
				State  string `json:"state"`
			} `json:"user_data"`
			Results []struct {
				Question string `json:"question" binding:"required"`
				Answer   string `json:"answer" binding:"required"`
			} `json:"results" binding:"required"`
		} `json:"attributes"`
	}

	var answerRequest AnswerRequest

	if err := c.ShouldBindJSON(&answerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var answerResponse AnswerResponse
	answerResponse.AccessToken = answerRequest.Attributes.AccessToken
	answerResponse.PollID = existingPoll.ID
	answerResponse.UserID = answerRequest.Attributes.UserID

	for _, result := range answerRequest.Attributes.Results {
		// Проверяем, существует ли вопрос с указанным текстом в опросе
		var existingQuestion models.Question
		if err := db.DB.Where("text = ? AND poll_id = ?", result.Question, existingPoll.ID).First(&existingQuestion).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Question not found (%s)", result.Question)})
			return
		}

		// Проверяем тип вопроса
		if existingQuestion.Type == "single" && strings.Count(result.Answer, ",") > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "For single type question, only one answer is allowed"})
			return
		} else if existingQuestion.Type == "multiple" && strings.Count(result.Answer, ",") == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "For multiple type question, at least two answers are required"})
			return
		}

		// Разделяем варианты ответов, если их несколько
		answers := strings.Split(result.Answer, ",")

		// Собираем все возможные ответы для данного вопроса
		var possibleAnswers []models.PossibleAnswer
		for _, answerText := range answers {
			answerText = strings.TrimSpace(answerText)
			// Проверяем, существует ли вариант ответа для данного вопроса
			var existingPossibleAnswer models.PossibleAnswer
			if err := db.DB.Where("text = ? AND question_id = ?", answerText, existingQuestion.ID).First(&existingPossibleAnswer).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Possible Answer (%s) not found for the specified question (%s)", answerText, result.Question)})
				return
			}
			possibleAnswers = append(possibleAnswers, existingPossibleAnswer)
		}

		// Создаем новый объект Answer и связываем с вариантами ответов
		newAnswer := models.Answer{
			UserID:          answerRequest.Attributes.UserID,
			AccessToken:     answerRequest.Attributes.AccessToken,
			QuestionID:      existingQuestion.ID,
			PollID:          existingPoll.ID,
			PossibleAnswers: possibleAnswers,
		}

		if err := db.DB.Create(&newAnswer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register answer"})
			return
		}

		// Формируем результаты ответа для структуры ответа
		var resultAnswers []uint
		for _, possibleAnswer := range possibleAnswers {
			resultAnswers = append(resultAnswers, possibleAnswer.ID)
		}

		// Добавляем результаты в структуру ответа
		answerResponse.Results = append(answerResponse.Results, struct {
			QuestionID uint   `json:"question_id"`
			Answers    []uint `json:"answers"`
		}{
			QuestionID: existingQuestion.ID,
			Answers:    resultAnswers,
		})
	}

	c.JSON(http.StatusOK, answerResponse)
}
