package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
)

func SetupDatabase(config Config) gorm.DB {
	// TODO: Catch the error here
	db, err := gorm.Open("postgres", config.DatabaseConnectionString)
	if err != nil {
		os.Exit(1)
	}

	db.AutoMigrate(&Photo{}, &Album{}, &Duplicate{})
	db.LogMode(config.LogDatabaseQueries)
	return db
}
