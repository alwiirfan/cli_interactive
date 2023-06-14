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

func (repository *WikisRepositoryImpl) GetAllWikisRepo() ([]*entity.Wikis, error) {
	var wikis []*entity.Wikis
	err := repository.DB.Find(&wikis).Error
	if err != nil {
		return nil, err
	}
	return wikis, nil
}

func (repository *WikisRepositoryImpl) GetWikisByIDRepo(id int) (*entity.Wikis, error) {
	wikis := &entity.Wikis{}
	err := repository.DB.First(wikis, id).Error
	if err != nil {
		return nil, err
	}
	return wikis, nil
}

func (repository *WikisRepositoryImpl) UpdateWikisRepo(wikis *entity.Wikis) error {
	return repository.DB.Save(wikis).Error
}

func (repository *WikisRepositoryImpl) DeleteWikisRepo(id int) error {
	return repository.DB.Delete(&entity.Wikis{}, id).Error
}
