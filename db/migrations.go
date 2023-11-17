package db

import (
	"github.com/ahmad1702/ultrawide-compat/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&models.Todo{})
}
