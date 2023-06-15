package database

import (
	"cli_interactive/internal/model/entity"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func (repository *WikisRepositoryImpl) UpdateTopicDescriptionRepo(wikis *entity.Wikis) error {
	return repository.DB.Save(wikis).Error
}

func (repository *WikisRepositoryImpl) DeleteWikisRepo(id int) error {
	return repository.DB.Delete(&entity.Wikis{}, id).Error
}

func (repository *WikisRepositoryImpl) GetWikisWithEmptyDescriptionRepo() ([]*entity.Wikis, error) {
	var wikis []*entity.Wikis
	err := repository.DB.Where("description = '' OR description IS NULL").Find(&wikis).Error
	if err != nil {
		return nil, err
	}
	return wikis, nil
}

func (repository *WikisRepositoryImpl) UpdateDescriptionFromWikipedia(topic string) (string, error) {
	url := fmt.Sprintf("https://id.wikipedia.org/wiki/%s", topic)

	// Membuat HTTP GET request ke URL Wikipedia
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch Wikipedia page: %s", response.Status)
	}

	// Menggunakan goquery untuk memparsing response body
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Mengambil paragraf pertama dari halaman Wikipedia
	firstParagraph := doc.Find("div.mw-parser-output > p").First().Text()

	return firstParagraph, nil
}

func (repository *WikisRepositoryImpl) UpdateUpdatedAt(wikis *entity.Wikis) error {
	wikis.Updated_at = time.Now()
	return repository.DB.Save(wikis).Error
}
