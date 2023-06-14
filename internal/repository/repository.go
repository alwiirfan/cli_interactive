package repository

import "cli_interactive/internal/model/entity"

type WikisRepository interface {
	AddTopicWikisRepo(wikis *entity.Wikis) error
	GetAllWikisRepo() ([]*entity.Wikis, error)
	GetWikisByIDRepo(id int) (*entity.Wikis, error)
	UpdateWikisRepo(wikis *entity.Wikis) error
	DeleteWikisRepo(id int) error
}
