package postrepository

import (
	"SociLinkApi/models"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

func CreatePost(post *models.Post, db *gorm.DB) error {
	result := db.Clauses(clause.Returning{}).Create(&post)

	return result.Error
}
