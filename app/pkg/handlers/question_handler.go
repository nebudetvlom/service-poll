package handlers

import (
	"net/http"

	"service-poll/pkg/db"
	"service-poll/pkg/models"

	"github.com/gin-gonic/gin"
)

// AddQuestion Добавляет вопрос к опросу по указанному идентификатору.
func AddQuestion(c *gin.Context) {
	// Получаем идентификатор опроса из параметра запроса
	pollID := c.Param("id")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll models.Poll
	if err := db.DB.First(&existingPoll, pollID).Error; err != nil {
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
	newQuestion := models.Question{
		Text:   questionData.Text,
		Type:   questionData.Type,
		PollID: existingPoll.ID,
	}

	// Сохраняем новый вопрос в базе данных
	if err := db.DB.Create(&newQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	// Сохраняем варианты ответов
	var createdAnswers []models.PossibleAnswer
	for _, answerText := range questionData.Answers {
		possibleAnswer := models.PossibleAnswer{
			Text:       answerText,
			QuestionID: newQuestion.ID,
		}

		if err := db.DB.Create(&possibleAnswer).Error; err != nil {
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
	// Получаем идентификатор опроса и вопроса из параметра запроса
	pollID := c.Param("id")
	questionID := c.Param("qid")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll models.Poll
	if err := db.DB.First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Проверяем, существует ли вопрос с указанным идентификатором в рамках данного опроса
	var existingQuestion models.Question
	if err := db.DB.Where("id = ? AND poll_id = ?", questionID, pollID).First(&existingQuestion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Теперь можешь получить связанные данные (например, PossibleAnswers) при необходимости:
	var possibleAnswers []models.PossibleAnswer
	if err := db.DB.Where("question_id = ?", existingQuestion.ID).Find(&possibleAnswers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch possible answers"})
		return
	}

	// Привязываем данные из запроса к структуре Question
	var updatedQuestionData struct {
		Text    string   `json:"text" binding:"required"`
		Type    string   `json:"type" binding:"required,oneof=single multiple"`
		Answers []string `json:"answers" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&updatedQuestionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка значения Type
	if updatedQuestionData.Type != "single" && updatedQuestionData.Type != "multiple" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be either 'single' or 'multiple'"})
		return
	}

	// Обновляем данные вопроса
	existingQuestion.Text = updatedQuestionData.Text
	existingQuestion.Type = updatedQuestionData.Type

	// Сохраняем изменения в базе данных
	if err := db.DB.Save(&existingQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
		return
	}

	// Обновляем варианты ответов
	for _, possibleAnswer := range possibleAnswers {
		// Проверяем, есть ли вариант ответа в списке обновленных
		found := false
		for _, updatedAnswer := range updatedQuestionData.Answers {
			if possibleAnswer.Text == updatedAnswer {
				found = true
				break
			}
		}

		// Если вариант ответа не найден, удаляем его
		if !found {
			if err := db.DB.Delete(&possibleAnswer).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete possible answer"})
				return
			}
			if err := db.DB.Where("possible_answer_id = ?", possibleAnswer.ID).Delete(&models.AnswerPossibleAnswer{}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exact answer"})
				return
			}
		}
	}

	// Добавляем новые варианты ответов
	for _, updatedAnswer := range updatedQuestionData.Answers {
		found := false
		for _, possibleAnswer := range possibleAnswers {
			if updatedAnswer == possibleAnswer.Text {
				found = true
				break
			}
		}

		// Если вариант ответа не найден, добавляем его
		if !found {
			newPossibleAnswer := models.PossibleAnswer{
				Text:       updatedAnswer,
				QuestionID: existingQuestion.ID,
			}

			if err := db.DB.Create(&newPossibleAnswer).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create possible answer"})
				return
			}
		}
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, existingQuestion)
}

// DeleteQuestion удаляет вопрос к опросу по указанному идентификатору.
func DeleteQuestion(c *gin.Context) {
	// Получаем идентификатор опроса и вопроса из параметра запроса
	pollID := c.Param("id")
	questionID := c.Param("qid")

	// Проверяем, существует ли опрос с указанным идентификатором
	var existingPoll models.Poll
	if err := db.DB.First(&existingPoll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Проверяем, существует ли вопрос с указанным идентификатором в рамках данного опроса
	var existingQuestion models.Question
	if err := db.DB.Where("id = ? AND poll_id = ?", questionID, pollID).First(&existingQuestion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Сохраняем название удаляемого вопроса
	deletedQuestionText := existingQuestion.Text

	// Собираем все идентификаторы возможных вариантов ответа для данного вопроса
	var possibleAnswerIDs []uint
	if err := db.DB.Model(&models.PossibleAnswer{}).Where("question_id = ?", existingQuestion.ID).Pluck("id", &possibleAnswerIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch possible answer IDs"})
		return
	}

	// Удаляем связанные записи из таблицы answer_possible_answers
	if err := db.DB.Where("possible_answer_id IN ?", possibleAnswerIDs).Delete(&models.AnswerPossibleAnswer{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete answer_possible_answers"})
		return
	}

	// Удаляем варианты ответов
	if err := db.DB.Where("question_id = ?", existingQuestion.ID).Delete(&models.PossibleAnswer{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete possible answers"})
		return
	}

	// Удаляем вопрос
	if err := db.DB.Delete(&existingQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	// Возвращаем успешный ответ с названием удаляемого вопроса
	c.JSON(http.StatusOK, gin.H{"message": "Question '" + deletedQuestionText + "' deleted successfully"})
}
