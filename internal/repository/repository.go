package repository

import "cli_interactive/internal/model/entity"

type WikisRepository interface {
	AddTopicWikisRepo(wikis *entity.Wikis) error
}
