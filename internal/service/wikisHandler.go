package service

import (
	"fmt"
	"log"
	"time"

	"cli_interactive/internal/database"
	"cli_interactive/internal/model/entity"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type HandlerImpl struct {
	wikisRepository *database.WikisRepositoryImpl
}

func NewHandlerImpl(repository *database.WikisRepositoryImpl) *HandlerImpl {
	return &HandlerImpl{
		wikisRepository: repository,
	}
}

func (handler *HandlerImpl) StartInteractiveCLI(cmd *cobra.Command, args []string) {
	prompt := promptui.Prompt{
		Label: "Enter a topic",
	}

	for {
		topic, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		wiki := &entity.Wikis{
			ID:          0,
			Topic:       topic,
			Description: "",
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}

		err = handler.wikisRepository.AddTopicWikisRepo(wiki)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Topic saved to the database.")
	}
}
