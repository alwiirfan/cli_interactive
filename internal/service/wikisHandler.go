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

// untuk mendapatkan semua data wiki dari repository
func (handler *HandlerImpl) GetAllWikisHandler(cmd *cobra.Command, args []string) {
	wikis, err := handler.wikisRepository.GetAllWikisRepo()
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
	} else {
		id = args[0]
	}

	wikiID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal("Invalid ID format:", err)
	}

	wiki, err := handler.wikisRepository.GetWikisByIDRepo(wikiID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %d, Topic: %s, Description: %s\n", wiki.ID, wiki.Topic, wiki.Description)
}

// Handler untuk mengambil input dari pengguna dan memperbarui topik dan deskripsi wiki
func (handler *HandlerImpl) UpdateTopicDescriptionHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Please provide a single argument: ID")
	}

	id := args[0]
	wikiID := dto.ConvertToInt(id)

	wiki, err := handler.wikisRepository.GetWikisByIDRepo(wikiID)
	if err != nil {
		log.Fatal(err)
	}

	prompt := promptui.Prompt{
		Label:    "Enter a new topic",
		Default:  wiki.Topic,
		Validate: dto.ValidateNonEmptyInput,
	}

	newTopic, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Mengambil deskripsi baru dari Wikipedia berdasarkan topik baru
	newDescription, err := handler.wikisRepository.UpdateDescriptionFromWikipedia(newTopic)
	if err != nil {
		log.Fatal(err)
	}

	wiki.Topic = newTopic
	wiki.Description = newDescription
	wiki.Updated_at = time.Now()

	err = handler.wikisRepository.UpdateTopicDescriptionRepo(wiki)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Wiki updated successfully.")
}

// untuk menghapus data wiki berdasarkan ID yang diberikan sebagai argumen.
func (handler *HandlerImpl) DeleteWikisHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Please provide a single argument: ID")
	}

	id := args[0]
	wikiID := dto.ConvertToInt(id)

	err := handler.wikisRepository.DeleteWikisRepo(wikiID)
	if err != nil {
		log.Fatal(err)
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
				log.Println(err)
				continue
			}

			// Mengupdate deskripsi wiki dengan nilai yang baru.
			wiki.Description = description

			// untuk mengupdate waktu terakhir diubah (updatedAt) dari wiki.
			err = handler.wikisRepository.UpdateUpdatedAt(wiki)
			if err != nil {
				log.Println(err)
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
