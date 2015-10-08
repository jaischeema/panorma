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
	validPhotoExts = []string{".jpg", ".jpeg", ".tiff", ".tif", ".gif", ".png"}
	validMovieExts = []string{".mov", ".m4v", ".mp4", ".mov"}
)

func ImportImages(c *cli.Context) {
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

	processPhotos(db, config.SourceFolderPath, config.DestinationFolderPath)
}

func createTreeFromDatabase(db gorm.DB) bktree.Node {
	log.Info("Initializing BKTree")

	var photos []Photo
	db.Select("id, hash_value").Find(&photos)
	if len(photos) > 0 {
		log.WithFields(log.Fields{
			"count": len(photos),
		}).Info("Photos found in database")

		firstPhoto := photos[0]
		tree := bktree.New(uint64(firstPhoto.HashValue), firstPhoto.Id)
		for _, photo := range photos[1:] {
			tree.Insert(uint64(photo.HashValue), photo.Id)
		}
		return tree
	} else {
		log.Info("No photos, creating empty tree")
		return bktree.Node{}
	}
}

const allowedHammingDistance = 10

func processPhotos(db gorm.DB, sourcePath string, destinationPath string) {
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
		isPhoto := strIn(fileExtension, validPhotoExts)

		if !info.IsDir() && (isMovie || isPhoto) {
			photo, err := processPhoto(itemPath, info)
			if err != nil {
				log.WithFields(log.Fields{
					"path": itemPath,
					"err":  err,
				}).Warn("Unable to open item")
				return err
			}

			log.WithFields(log.Fields{
				"height":   photo.Height,
				"width":    photo.Width,
				"taken_at": photo.TakenAt,
				"path":     photo.Path,
			}).Info("Processed attributes")

			if photo.ExistsInDatabase(db) {
				// TODO: Move the duplicates to duplicate folder
				log.WithFields(log.Fields{
					"path": itemPath,
				}).Info("Already in database")
			} else {
				tx := db.Begin()

				photo.Size = info.Size()
				photo.Ext = path.Ext(itemPath)
				photo.Name = strings.TrimSuffix(path.Base(itemPath), photo.Ext)
				photo.IsVideo = isMovie

				var photoHashValue uint64
				if isPhoto {
					photoHashValue = bktree.PHashValueForImage(itemPath)
					photo.HashValue = int64(photoHashValue)
				}

				db.Save(&photo)

				if isPhoto {
					if treeNeedsRootNode {
						tree = bktree.New(photoHashValue, photo.Id)
						treeNeedsRootNode = false
					} else {
						tree.Insert(photoHashValue, photo.Id)
					}

					duplicateIds := tree.Find(photoHashValue, allowedHammingDistance)
					for _, duplicateId := range duplicateIds {
						if duplicateId == photo.Id {
							continue
						}
						var similarPhoto SimilarPhoto
						db.Where(SimilarPhoto{
							PhotoId:        photo.Id,
							SimilarPhotoId: duplicateId.(int64),
						}).FirstOrCreate(&similarPhoto)
					}
				}

				err = moveFileInTransaction(itemPath, destinationPath, photo.Path)

				if err != nil {
					tx.Rollback()
					return err
				}

				log.WithFields(log.Fields{
					"path":        itemPath,
					"destination": photo.Path,
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

func processPhoto(path string, info os.FileInfo) (photo Photo, err error) {
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

		return extractImageWithDimensions(path, info), nil
	}

	tm, err := data.DateTime()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Unable to find EXIF datetime")

		return extractImageWithDimensions(path, info), nil
	}

	lat, long, _ := data.LatLong()
	widthTag, err := data.Get("PixelXDimension")
	heightTag, err := data.Get("PixelYDimension")

	if widthTag == nil || heightTag == nil {
		log.Warn("WidthTag OR HeightTag is not available")
		return extractImageWithDimensions(path, info), nil
	}

	width, _ := widthTag.Int(0)
	height, _ := heightTag.Int(0)

	photo = Photo{
		Path:    calculatePhotoTimedPath(path, tm),
		TakenAt: tm,
		Lat:     lat,
		Lng:     long,
		Height:  height,
		Width:   width,
	}
	return
}

func calculatePhotoTimedPath(filepath string, takenAt time.Time) string {
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

func extractImageWithDimensions(path string, info os.FileInfo) Photo {
	width, height := imageDimensions(path)

	return Photo{
		Path:    calculatePhotoTimedPath(path, info.ModTime()),
		TakenAt: info.ModTime(),
		Height:  height,
		Width:   width,
	}
}

func imageDimensions(imagePath string) (int, int) {
	file, _ := os.Open(imagePath)
	defer file.Close()
	image, _, _ := image.DecodeConfig(file)
	return image.Width, image.Height
}
