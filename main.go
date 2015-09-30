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

	configFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "config,c",
			EnvVar: "CONFIG_PATH",
			Value:  "./config.json",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "import",
			ShortName:   "i",
			Description: "import pictures",
			Action:      ImportImages,
			Flags:       configFlag,
		},
		{
			Name:      "server",
			ShortName: "s",
			Action:    RunServer,
			Flags:     configFlag,
		},
		{
			Name:      "thumbnails",
			ShortName: "t",
			Action:    ThumbnailImages,
			Flags:     configFlag,
		},
	}

	app.RunAndExitOnError()
}
