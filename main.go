package main

import (
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Panorma"
	app.Author = "Jais Cheema"
	app.Version = "0.0.1"
	app.Email = "jaischeema@gmail.com"
	app.Usage = "panorma server"

	commonFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "source_path,s",
			EnvVar: "SOURCE_PATH",
			Value:  "/Users/jais/Pictures/Originals",
		},
		cli.StringFlag{
			Name:   "destination_path,e",
			EnvVar: "DESTINATION_PATH",
			Value:  "/Users/jais/Pictures/Archive",
		},
		cli.StringFlag{
			Name:   "thumbnails_path,t",
			EnvVar: "THUMBNAILS_PATH",
			Value:  "/Users/jais/Pictures/Thumbnails",
		},
		cli.StringFlag{
			Name:   "database_url,d",
			EnvVar: "DATABASE_URL",
			Value:  "dbname=panorma_dev sslmode=disable",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "import",
			ShortName:   "i",
			Description: "import pictures",
			Action:      ImportImages,
			Flags:       commonFlags,
		},
		{
			Name:      "server",
			ShortName: "s",
			Action:    RunServer,
			Flags:     commonFlags,
		},
	}

	app.RunAndExitOnError()
}
