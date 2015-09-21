package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func ThumbnailImages(c *cli.Context) {
	config, err := LoadConfig(c.String("config"))
	if err != nil {
		os.Exit(1)
	}

	db := SetupDatabase(config.DatabaseConnectionString)

}
