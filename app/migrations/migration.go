package migrations

import (
	"service-poll/pkg/db"
	"service-poll/pkg/models"

	"github.com/gin-gonic/gin"
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
