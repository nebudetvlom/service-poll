package handlers

import (
	"fmt"
	"net/http"
	"service-poll/pkg/db"
	"service-poll/pkg/models"

	"github.com/gin-gonic/gin"
)

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
	newPoll := models.Poll{
		Title: pollData.Title,
		URL:   pollData.URL,
	}

	// Сохраняем новый объект в базе данных
	if err := db.DB.Create(&newPoll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
		return
	}

	c.JSON(http.StatusOK, newPoll)
}

// GetPoll возвращает опрос с вопросами по указанному идентификатору.
func GetPoll(c *gin.Context) {
	// Получаем идентификатор опроса из параметра запроса
	pollID := c.Param("id")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll models.Poll
	if err := db.DB.Preload("Questions.PossibleAnswer").Preload("Answers.PossibleAnswers").First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Формируем ответ с данными об опросе
	var pollResponse struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		URL       string `json:"url"`
		Questions []struct {
			ID              uint   `json:"id"`
			Text            string `json:"text"`
			Type            string `json:"type"`
			PossibleAnswers []struct {
				ID   uint   `json:"id"`
				Text string `json:"text"`
			} `json:"possible_answers"`
		} `json:"questions"`
	}

	pollResponse.ID = existingPoll.ID
	pollResponse.Title = existingPoll.Title
	pollResponse.URL = existingPoll.URL

	for _, question := range existingPoll.Questions {
		var questionResponse struct {
			ID              uint   `json:"id"`
			Text            string `json:"text"`
			Type            string `json:"type"`
			PossibleAnswers []struct {
				ID   uint   `json:"id"`
				Text string `json:"text"`
			} `json:"possible_answers"`
		}

		questionResponse.ID = question.ID
		questionResponse.Text = question.Text
		questionResponse.Type = question.Type

		for _, possibleAnswer := range question.PossibleAnswer {
			questionResponse.PossibleAnswers = append(questionResponse.PossibleAnswers, struct {
				ID   uint   `json:"id"`
				Text string `json:"text"`
			}{
				ID:   possibleAnswer.ID,
				Text: possibleAnswer.Text,
			})
		}

		pollResponse.Questions = append(pollResponse.Questions, questionResponse)
	}

	c.JSON(http.StatusOK, pollResponse)
}

// UpdatePoll изменяет опрос по указанному идентификатору.
func UpdatePoll(c *gin.Context) {
	// Получаем идентификатор опроса из параметра запроса
	pollID := c.Param("id")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll models.Poll
	if err := db.DB.First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Извлекаем данные обновления опроса из тела запроса
	var updateData struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}

	// Проверяем и извлекаем данные из тела запроса
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем данные опроса, если они предоставлены
	if updateData.Title != "" {
		existingPoll.Title = updateData.Title
	}

	if updateData.URL != "" {
		existingPoll.URL = updateData.URL
	}

	// Сохраняем обновленный опрос в базе данных
	if err := db.DB.Save(&existingPoll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update poll"})
		return
	}

	c.JSON(http.StatusOK, existingPoll)
}

// DeletePoll удаляет опрос по указанному идентификатору.
func DeletePoll(c *gin.Context) {
	// Получаем идентификатор опроса из параметра запроса
	pollID := c.Param("id")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll models.Poll
	if err := db.DB.Preload("Questions.PossibleAnswer.Answers").Preload("Answers").First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Сохраняем название опроса перед удалением
	pollTitle := existingPoll.Title

	// Удаляем связанные ответы
	for _, answer := range existingPoll.Answers {
		if err := db.DB.Delete(&answer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete answers"})
			return
		}
	}

	// Удаляем связанные вопросы и их варианты ответов
	for _, question := range existingPoll.Questions {
		// Удаляем варианты ответов
		for _, possibleAnswer := range question.PossibleAnswer {
			if err := db.DB.Delete(&possibleAnswer).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete possible answers"})
				return
			}
		}

		// Удаляем вопрос
		if err := db.DB.Delete(&question).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete questions"})
			return
		}
	}

	// Удаляем опрос
	if err := db.DB.Delete(&existingPoll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete poll"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Poll '%s' deleted successfully", pollTitle)})
}

// GetPollResults возвращает результаты опроса по указанному идентификатору.
func GetPollResults(c *gin.Context) {
	// Получаем идентификатор опроса из параметра запроса
	pollID := c.Param("id")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll models.Poll
	if err := db.DB.Preload("Questions.PossibleAnswer.Answers.PossibleAnswers").Preload("Answers.PossibleAnswers").First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Структура для хранения результатов
	var pollResults struct {
		ID      uint `json:"id"`
		Results []struct {
			Question         string `json:"question"`
			Answer           string `json:"answer"`
			AnswerCnt        int    `json:"answer_cnt"`
			AnswerPercentage string `json:"answer_percentage"`
		} `json:"results"`
	}

	pollResults.ID = existingPoll.ID

	// Пройдемся по всем вопросам
	for _, question := range existingPoll.Questions {
		// Карта для подсчета количества ответов на каждый вариант ответа
		answerCounts := make(map[string]int)

		var answers []models.Answer
		if err := db.DB.Where("question_id = ?", question.ID).Find(&answers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch answers"})
			return
		}

		for _, answer := range answers {
			var answerPossibleAnswers []models.AnswerPossibleAnswer

			// Загружаем связанные записи из таблицы answer_possible_answers для данного ответа
			if err := db.DB.Where("answer_id = ?", answer.ID).Find(&answerPossibleAnswers).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch possible answers"})
				return
			}

			for _, apa := range answerPossibleAnswers {
				// Получаем связанный возможный ответ
				var possibleAnswer models.PossibleAnswer
				if err := db.DB.First(&possibleAnswer, apa.PossibleAnswerID).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch possible answer"})
					return
				}

				// Увеличиваем счетчик ответов для данного варианта ответа
				answerCounts[possibleAnswer.Text]++
			}
		}

		// Подсчитываем общее количество ответов на вопрос
		totalAnswers := len(answers)

		// Пройдемся по всем вариантам ответа и рассчитаем процентное соотношение
		for answerText, answerCount := range answerCounts {
			percentage := (float64(answerCount) / float64(totalAnswers)) * 100

			// Добавляем результат в структуру
			result := struct {
				Question         string `json:"question"`
				Answer           string `json:"answer"`
				AnswerCnt        int    `json:"answer_cnt"`
				AnswerPercentage string `json:"answer_percentage"`
			}{
				Question:         question.Text,
				Answer:           answerText,
				AnswerCnt:        answerCount,
				AnswerPercentage: fmt.Sprintf("%.3f", percentage),
			}

			pollResults.Results = append(pollResults.Results, result)
		}
	}

	c.JSON(http.StatusOK, pollResults)
}
