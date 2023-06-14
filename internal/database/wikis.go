package database

import (
	"cli_interactive/internal/model/entity"

	"gorm.io/gorm"
)

type WikisRepositoryImpl struct {
	DB *gorm.DB
}

func NewWikisRepository(DB *gorm.DB) *WikisRepositoryImpl {
	return &WikisRepositoryImpl{
		DB: DB,
	}
}

func (repository *WikisRepositoryImpl) AddTopicWikisRepo(wikis *entity.Wikis) error {
	return repository.DB.Create(wikis).Error
}
