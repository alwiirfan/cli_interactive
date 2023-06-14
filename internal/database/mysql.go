package database

import (
	"cli_interactive/internal/config"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigration "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {

	fmt.Printf("%+v\n ", cfg)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})

	if err != nil {
		panic("cannot connect to database")
	}

	logrus.Info("connect to database successfully")
	return db, err
}

// Migrate
func Migrate(db *gorm.DB) error {
	logrus.Info("running database migration")

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := mysqlMigration.WithInstance(sqlDB, &mysqlMigration.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"mysql", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err == migrate.ErrNoChange {
		logrus.Info("no schema changes to apply")
		return nil
	}

	return err
}
