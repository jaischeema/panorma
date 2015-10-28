package main

import (
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/jaischeema/panorma/bktree"
	"github.com/jinzhu/gorm"
	"github.com/rwcarlsen/goexif/exif"
)

type Result struct {
	Id        int64
	HashValue int64
}

var (
	validImageExts = []string{".jpg", ".jpeg", ".tiff", ".tif", ".gif", ".png"}
	validMovieExts = []string{".mov", ".m4v", ".mp4", ".mov"}
)

func ImportMedia(c *cli.Context) {
	config, err := LoadConfig(c.String("config"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	db := SetupDatabase(config.DatabaseConnectionString)

	log.WithFields(log.Fields{
		"source":      config.SourceFolderPath,
		"destination": config.DestinationFolderPath,
	}).Info("Starting image import.")

	processItemsFromSource(db, config.SourceFolderPath, config.DestinationFolderPath)
}

func createTreeFromDatabase(db gorm.DB) bktree.Node {
	log.Info("Initializing BKTree")

	var media []Media
	db.Select("id, hash_value").Find(&media)
	if len(media) > 0 {
		log.WithFields(log.Fields{
			"count": len(media),
		}).Info("Media found in database")

		firstItem := media[0]
		tree := bktree.New(uint64(firstItem.HashValue), firstItem.Id)
		for _, item := range media[1:] {
			tree.Insert(uint64(item.HashValue), item.Id)
		}
		return tree
	} else {
		log.Info("No media, creating empty tree")
		return bktree.Node{}
	}
}

const allowedHammingDistance = 10

func processItemsFromSource(db gorm.DB, sourcePath string, destinationPath string) {
	tree := createTreeFromDatabase(db)
	treeNeedsRootNode := (tree.HashValue == 0)

	walkFunc := func(itemPath string, info os.FileInfo, err error) error {
		if err != nil {
			log.WithFields(log.Fields{
				"path": itemPath,
				"err":  err,
			}).Warn("Unable to process item")

			return err
		}

		fileExtension := filepath.Ext(itemPath)
		isMovie := strIn(fileExtension, validMovieExts)
		isImage := strIn(fileExtension, validImageExts)

		if !info.IsDir() && (isMovie || isImage) {
			mediaItem, err := processMediaItem(itemPath, info)
			if err != nil {
				log.WithFields(log.Fields{
					"path": itemPath,
					"err":  err,
				}).Warn("Unable to open item")
				return err
			}

			log.WithFields(log.Fields{
				"height":   mediaItem.Height,
				"width":    mediaItem.Width,
				"taken_at": mediaItem.TakenAt,
				"path":     mediaItem.Path,
			}).Info("Processed attributes")

			if mediaItem.ExistsInDatabase(db) {
				// TODO: Move the duplicates to duplicate folder
				log.WithFields(log.Fields{
					"path": itemPath,
				}).Info("Already in database")
			} else {
				tx := db.Begin()

				mediaItem.Size = info.Size()
				mediaItem.Ext = strings.ToLower(path.Ext(itemPath))
				mediaItem.Name = strings.TrimSuffix(path.Base(itemPath), mediaItem.Ext)
				mediaItem.IsVideo = isMovie

				var mediaItemHashValue uint64
				if isImage {
					mediaItemHashValue = bktree.PHashValueForImage(itemPath)
					mediaItem.HashValue = int64(mediaItemHashValue)
				}

				db.Save(&mediaItem)

				if isImage {
					if treeNeedsRootNode {
						tree = bktree.New(mediaItemHashValue, mediaItem.Id)
						treeNeedsRootNode = false
					} else {
						tree.Insert(mediaItemHashValue, mediaItem.Id)
					}

					duplicateIds := tree.Find(mediaItemHashValue, allowedHammingDistance)
					for _, duplicateId := range duplicateIds {
						if duplicateId == mediaItem.Id {
							continue
						}
						var resemblance Resemblance
						db.Where(Resemblance{
							MediaId:           mediaItem.Id,
							ResemblingMediaId: duplicateId.(int64),
						}).FirstOrCreate(&resemblance)
					}
				}

				err = moveFileInTransaction(itemPath, destinationPath, mediaItem.Path)

				if err != nil {
					tx.Rollback()
					return err
				}

				log.WithFields(log.Fields{
					"path":        itemPath,
					"destination": mediaItem.Path,
				}).Info("Moved to archive")

				tx.Commit()
			}
		}
		return nil
	}

	err := filepath.Walk(sourcePath, walkFunc)
	if err != nil {
		fmt.Println(err)
	}
}

func moveFileInTransaction(itemPath string, destinationRoot string, destinationPath string) error {
	sourceinfo, err := os.Stat(destinationRoot)
	if err != nil {
		log.WithFields(log.Fields{
			"destination": destinationRoot,
			"err":         err,
		}).Error("Destination Directory cannot be accessed.")

		return err
	}

	fullPath := path.Join(destinationRoot, destinationPath)
	basePath, _ := path.Split(fullPath)

	err1 := os.MkdirAll(basePath, sourceinfo.Mode())

	if err1 != nil {
		log.WithFields(log.Fields{
			"destination": basePath,
			"err":         err1,
		}).Error("Unable to create the directory structure.")

		return err1
	}

	err2 := os.Rename(itemPath, fullPath)

	if err2 != nil {
		log.WithFields(log.Fields{
			"item":     itemPath,
			"fullPath": fullPath,
			"err":      err2,
		}).Error("Unable to move the file.")

		return err2
	}
	return nil
}

func processMediaItem(path string, info os.FileInfo) (mediaItem Media, err error) {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Processing item")

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := exif.Decode(file)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Unable to decode EXIF data")

		return extractItemWithDimensions(path, info), nil
	}

	tm, err := data.DateTime()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Unable to find EXIF datetime")

		return extractItemWithDimensions(path, info), nil
	}

	lat, long, _ := data.LatLong()
	widthTag, err := data.Get("PixelXDimension")
	heightTag, err := data.Get("PixelYDimension")

	if widthTag == nil || heightTag == nil {
		log.Warn("WidthTag OR HeightTag is not available")
		return extractItemWithDimensions(path, info), nil
	}

	width, _ := widthTag.Int(0)
	height, _ := heightTag.Int(0)

	mediaItem = Media{
		Path:    calculatePathWithDate(path, tm),
		TakenAt: tm,
		Lat:     lat,
		Lng:     long,
		Height:  height,
		Width:   width,
	}
	return
}

func calculatePathWithDate(filepath string, takenAt time.Time) string {
	timeFormat := takenAt.Format("2006/01-January/02")

	_, file := path.Split(filepath)
	return path.Join(timeFormat, file)
}

func strIn(a string, list []string) bool {
	for _, b := range list {
		if strings.EqualFold(a, b) {
			return true
		}
	}
	return false
}

func extractItemWithDimensions(path string, info os.FileInfo) Media {
	width, height := dimensionsFromFile(path)

	return Media{
		Path:    calculatePathWithDate(path, info.ModTime()),
		TakenAt: info.ModTime(),
		Height:  height,
		Width:   width,
	}
}

func dimensionsFromFile(imagePath string) (int, int) {
	file, _ := os.Open(imagePath)
	defer file.Close()
	image, _, _ := image.DecodeConfig(file)
	return image.Width, image.Height
}
