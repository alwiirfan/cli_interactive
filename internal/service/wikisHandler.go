package service

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"cli_interactive/internal/database"
	"cli_interactive/internal/model/dto"
	"cli_interactive/internal/model/entity"

	"github.com/go-co-op/gocron"
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

// menambahkan topik
func (handler *HandlerImpl) AddTopicHandler(cmd *cobra.Command, args []string) {
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
			log.Println("topic cannot be empty", err)
		}

		fmt.Println("Topic saved to the database.")
	}
}

// untuk mendapatkan semua data wiki dari repository
func (handler *HandlerImpl) GetAllWikisHandler(cmd *cobra.Command, args []string) {
	wikis, err := handler.wikisRepository.GetAllWikisRepo()
	if err != nil {
		log.Println("failed to get all wikis:", err)
	}

	for _, wiki := range wikis {
		fmt.Printf("ID: %d, Topic: %s, Description: %s\n", wiki.ID, wiki.Topic, wiki.Description)
	}
}

// Menggunakan prompt untuk meminta pengguna memasukkan ID
func (handler *HandlerImpl) GetWikisByIDHandler(cmd *cobra.Command, args []string) {
	var id string
	if len(args) != 1 {
		prompt := promptui.Prompt{
			Label: "Enter an ID",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("ID cannot be empty")
				}
				return nil
			},
		}

		var err error
		id, err = prompt.Run()
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		id = args[0]
	}

	wikiID, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Invalid ID format:", err)
	}

	wiki, err := handler.wikisRepository.GetWikisByIDRepo(wikiID)
	if err != nil {
		log.Println("failed to get wiki by id:", err)
	}

	fmt.Printf("ID: %d, Topic: %s, Description: %s\n", wiki.ID, wiki.Topic, wiki.Description)
}

// Menggunakan prompt untuk meminta pengguna memasukkan ID untuk update topik dan description
func (handler *HandlerImpl) UpdateTopicDescriptionHandler(cmd *cobra.Command, args []string) {
	var id string
	var topic string

	if len(args) == 2 {
		id = args[0]
		topic = args[1]
	} else {

		// Menggunakan prompt untuk meminta ID
		idPrompt := promptui.Prompt{
			Label:    "Enter the ID",
			Validate: dto.ValidateNonEmptyInput,
		}

		var err error

		id, err = idPrompt.Run()
		if err != nil {
			log.Println(err.Error())
		}

		// Menggunakan prompt untuk meminta topik
		topicPrompt := promptui.Prompt{
			Label:    "Enter a new topic",
			Validate: dto.ValidateNonEmptyInput,
		}

		topic, err = topicPrompt.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}

	wikiID := dto.ConvertToInt(id)

	wiki, err := handler.wikisRepository.GetWikisByIDRepo(wikiID)
	if err != nil {
		log.Println("failed to get wiki by id:", err)
	}

	// Mengambil deskripsi baru dari Wikipedia berdasarkan topik baru
	newDescription, err := handler.wikisRepository.UpdateDescriptionFromWikipedia(topic)
	if err != nil {
		log.Println("failed to update description from wikipedia:", err)
	}

	wiki.Topic = topic
	wiki.Description = newDescription
	wiki.Updated_at = time.Now()

	err = handler.wikisRepository.UpdateTopicDescriptionRepo(wiki)
	if err != nil {
		log.Println("failed to update topic and description wiki:", err)
	}

	fmt.Println("Wiki updated successfully.")
}

// Menghapus data wiki berdasarkan ID yang diberikan sebagai argumen.
func (handler *HandlerImpl) DeleteWikisHandler(cmd *cobra.Command, args []string) {
	var id string

	if len(args) == 1 {
		id = args[0]
	} else {
		// Menggunakan prompt untuk meminta ID
		idPrompt := promptui.Prompt{
			Label:    "Enter the ID",
			Validate: dto.ValidateNonEmptyInput,
		}

		var err error
		id, err = idPrompt.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}

	wikiID := dto.ConvertToInt(id)

	err := handler.wikisRepository.DeleteWikisRepo(wikiID)
	if err != nil {
		log.Println("failed to delete wiki:", err)
	}

	fmt.Println("Wiki deleted successfully.")
}

func (handler *HandlerImpl) WorkerHandler(cmd *cobra.Command, args []string) {
	// Membuat sebuah scheduler baru
	scheduler := gocron.NewScheduler(time.UTC)

	// Mengatur jadwal untuk menjalankan fungsi yang diberikan setiap 1 menit.
	_, err := scheduler.Every(1).Minute().Do(func() {

		// untuk mendapatkan daftar wikis yang memiliki deskripsi kosong.
		wikis, err := handler.wikisRepository.GetWikisWithEmptyDescriptionRepo()
		if err != nil {
			log.Println("failed to get wikis:", err)
			return
		}

		// Melakukan iterasi pada setiap wiki dalam daftar wikis.
		for _, wiki := range wikis {

			// Mencetak informasi ID, topik, dan deskripsi dari wiki saat ini menggunakan log.Printf().
			log.Printf("ID: %d, Topic: %s, Description: %s\n", wiki.ID, wiki.Topic, wiki.Description)

			// untuk mengupdate deskripsi wiki dari sumber Wikipedia.
			description, err := handler.wikisRepository.UpdateDescriptionFromWikipedia(wiki.Topic)
			if err != nil {
				log.Println("failed to update description from wikipedia:", err)
				continue
			}

			// Mengupdate deskripsi wiki dengan nilai yang baru.
			wiki.Description = description

			// untuk mengupdate waktu terakhir diubah (updatedAt) dari wiki.
			err = handler.wikisRepository.UpdateUpdatedAt(wiki)
			if err != nil {
				log.Println("failed to update updated_at wiki:", err)
				continue
			}

			log.Println("Description updated successfully.")
		}
	})

	if err != nil {
		log.Fatal("failed to schedule worker:", err)
	}

	// Memulai scheduler dan menjalankan pekerjaan secara terjadwal.
	scheduler.StartBlocking()
}
