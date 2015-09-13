package main

import (
	"github.com/codegangsta/cli"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func SetupDatabase(c *cli.Context) gorm.DB {
	databasePath := c.String("database_url")

	db, err := gorm.Open("postgres", databasePath)
	if err != nil {
		panic("Unable to open database")
	}

	db.AutoMigrate(&Photo{}, &SimilarPhoto{})

	return db
}
