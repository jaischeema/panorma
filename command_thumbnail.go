package main

import (
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	bimg "gopkg.in/h2non/bimg.v0"
)

func ThumbnailImages(c *cli.Context) {
	config, err := LoadConfig(c.String("config"))
	if err != nil {
		os.Exit(1)
	}

	db := SetupDatabase(config.DatabaseConnectionString)
	photos := PhotosNotThumbnailed(db)

	for _, photo := range photos {
		fullPath := path.Join(config.DestinationFolderPath, photo.Path)
		for name, size := range ThumbnailSizes {
			destinationFolder := path.Join(config.ThumbnailsFolderPath, PartitionIdAsPath(photo.Id))
			err := os.MkdirAll(destinationFolder, 0755)
			if err != nil {
				log.Fatalf(err.Error())
				panic(1)
			}
			thumbnailPath := path.Join(destinationFolder, name+ThumbnailExtension)
			err = resizeWithLib(fullPath, thumbnailPath, size.Width, size.Height)
			if err != nil {
				log.Printf("Unable to create thumbnail for %s", fullPath)
			} else {
				photo.Thumbnailed = true
				db.Save(&photo)
				log.Printf("Done: (%s) (%s)", name, fullPath)
			}
		}
	}
}

func resizeWithLib(sourceFile string, destinationFile string, width int, height int) error {
	options := bimg.Options{
		Width:     width,
		Height:    height,
		Crop:      true,
		Quality:   90,
		Interlace: true,
	}

	buffer, err := bimg.Read(sourceFile)
	if err != nil {
		return err
	}

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		return err
	}

	err = bimg.Write(destinationFile, newImage)
	if err != nil {
		return err
	}
	return nil
}
