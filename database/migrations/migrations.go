package migrations

import (
	"github.com/sampalm/projectapi/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&models.API{}, &models.Project{})
}
