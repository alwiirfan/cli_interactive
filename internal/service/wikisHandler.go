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

func (handler *HandlerImpl) GetAllWikisHandler(cmd *cobra.Command, args []string) {
	wikis, err := handler.wikisRepository.GetAllWikisRepo()
	if err != nil {
		log.Fatal(err)
	}

	for _, wiki := range wikis {
		fmt.Printf("ID: %d, Topic: %s, Description: %s\n", wiki.ID, wiki.Topic, wiki.Description)
	}
}

func (handler *HandlerImpl) GetWikisByIDHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Please provide a single argument: ID")
	}

	id := args[0]
	wikiID := convertToInt(id)

	wiki, err := handler.wikisRepository.GetWikisByIDRepo(wikiID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %d, Topic: %s, Description: %s\n", wiki.ID, wiki.Topic, wiki.Description)
}

func (handler *HandlerImpl) UpdateWikisHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Please provide a single argument: ID")
	}

	id := args[0]
	wikiID := convertToInt(id)

	wiki, err := handler.wikisRepository.GetWikisByIDRepo(wikiID)
	if err != nil {
		log.Fatal(err)
	}

	prompt := promptui.Prompt{
		Label:    "Enter a new topic",
		Default:  wiki.Topic,
		Validate: validateNonEmptyInput,
	}

	newTopic, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	wiki.Topic = newTopic
	wiki.Updated_at = time.Now()

	err = handler.wikisRepository.UpdateWikisRepo(wiki)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Wiki updated successfully.")
}

func (handler *HandlerImpl) DeleteWikisHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Please provide a single argument: ID")
	}

	id := args[0]
	wikiID := convertToInt(id)

	err := handler.wikisRepository.DeleteWikisRepo(wikiID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Wiki deleted successfully.")
}

func (handler *HandlerImpl) AddTopicHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Please provide a single argument: Topic")
	}

	topic := args[0]

	wiki := &entity.Wikis{
		ID:          0,
		Topic:       topic,
		Description: "",
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	err := handler.wikisRepository.AddTopicWikisRepo(wiki)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Topic saved to the database.")
}

// Helper function to convert string to int.
func convertToInt(str string) int {
	var val int
	_, err := fmt.Sscanf(str, "%d", &val)
	if err != nil {
		log.Fatal("Invalid ID provided")
	}

	return val
}

// Helper function to validate non-empty input.
func validateNonEmptyInput(input string) error {
	if input == "" {
		return fmt.Errorf("input must not be empty")
	}
	return nil
}
