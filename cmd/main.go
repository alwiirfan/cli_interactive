package main

import (
	"cli_interactive/internal/config"
	"cli_interactive/internal/database"
	"cli_interactive/internal/service"
	"errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {

	// Set log level to display only INFO level
	logrus.SetLevel(logrus.InfoLevel)

	// Set log formatter to text format
	logrus.SetFormatter(&logrus.TextFormatter{})

	var cfg config.Config

	// initialize config
	envConfigErr := envconfig.Process("", &cfg)
	if envConfigErr != nil {
		logrus.Fatal(errors.New("initialize config is error"))
	}

	// initialize database connection
	db, err := database.ConnectDB(&cfg)
	if err != nil {
		return
	}

	// initialize migration
	err = database.Migrate(db)
	if err != nil {
		logrus.Fatalf("error running schema migration %v", err)
	}

	// initialize repository wikis
	wikisRepository := database.NewWikisRepository(db)

	// initialize handler wikis
	wikisHandler := service.NewHandlerImpl(wikisRepository)

	rootCmd := &cobra.Command{
		Use: "app",
	}

	// menambahkan topic
	addWikisCmd := &cobra.Command{
		Use:   "add-topic",
		Short: "Add wikis to the database",
		Run:   wikisHandler.AddTopicHandler,
	}

	rootCmd.AddCommand(addWikisCmd)

	//  Mengambil semua wiki dari database
	getAllWikisCmd := &cobra.Command{
		Use:   "get-all-wikis",
		Short: "Get all wikis from the database",
		Run:   wikisHandler.GetAllWikisHandler,
	}

	rootCmd.AddCommand(getAllWikisCmd)

	// Mengambil wiki berdasarkan ID dari database
	getWikiByIDCmd := &cobra.Command{
		Use:   "get-wiki-by-id",
		Short: "Get a wiki by ID from the database",
		Run:   wikisHandler.GetWikisByIDHandler,
	}

	rootCmd.AddCommand(getWikiByIDCmd)

	// Mengupdate topic dan deskripsi berdasarkan ID dari database
	updateTopicDescriptionHandler := &cobra.Command{
		Use:   "update-topic-description-wiki",
		Short: "Update a wiki in the database",
		Run:   wikisHandler.UpdateTopicDescriptionHandler,
	}

	rootCmd.AddCommand(updateTopicDescriptionHandler)

	//  Menghapus wiki dari database
	deleteWikisHandler := &cobra.Command{
		Use:   "delete-wiki",
		Short: "Delete a wiki from the database",
		Run:   wikisHandler.DeleteWikisHandler,
	}

	rootCmd.AddCommand(deleteWikisHandler)

	// Menjalankan worker
	runWorker := &cobra.Command{
		Use:   "worker",
		Short: "Run worker",
		Run:   wikisHandler.WorkerHandler,
	}

	rootCmd.AddCommand(runWorker)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal("application cannot running", err)
	}
}
