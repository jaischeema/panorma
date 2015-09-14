package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jinzhu/gorm"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"os"
	"panorma/bktree"
	"path"
	"path/filepath"
	"time"
)

type Result struct {
	Id        int64
	HashValue int64
}

var validExts = []string{".jpg", ".jpeg", ".tiff", ".tif", ".gif", ".png", ".JPG", ".mov", ".m4v", ".3gp"}

func ImportImages(c *cli.Context) {
	db := SetupDatabase(c)
	sourcePath := c.String("source_path")
	destinationPath := c.String("destination_path")

	processPhotos(db, sourcePath, destinationPath)
}

func createTreeFromDatabase(db gorm.DB) bktree.Node {
	var photos []Photo
	db.Select("id, hash_value").Find(&photos)
	if len(photos) > 0 {
		firstPhoto := photos[0]
		tree := bktree.New(uint64(firstPhoto.HashValue), firstPhoto.Id)
		for _, photo := range photos[1:] {
			tree.Insert(uint64(photo.HashValue), photo.Id)
		}
		return tree
	} else {
		return bktree.Node{}
	}
}

const allowedHammingDistance = 10

func processPhotos(db gorm.DB, sourcePath string, destinationPath string) {
	tree := createTreeFromDatabase(db)
	treeNeedsRootNode := (tree.HashValue == 0)

	walkFunc := func(itemPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strIn(filepath.Ext(itemPath), validExts) {
			photo, err := processPhoto(itemPath, info)
			if err != nil {
				fmt.Println("Error opening file: ", itemPath)
				return err
			}

			if photo.ExistsInDatabase(db) {
				fmt.Println("Skipping: ", itemPath)
			} else {
				tx := db.Begin()

				photo.Size = info.Size()
				photoHashValue := bktree.PHashValueForImage(itemPath)
				photo.HashValue = int64(photoHashValue)

				db.Save(&photo)

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

				err = moveFileInTransaction(itemPath, destinationPath, photo.Path)

				if err != nil {
					tx.Rollback()
					return err
				}
				fmt.Println("Move to Archive directory :", itemPath)
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
		fmt.Println("Destination Directory cannot be accessed.")
		return err
	}

	fullPath := path.Join(destinationRoot, destinationPath)
	basePath, _ := path.Split(fullPath)

	err1 := os.MkdirAll(basePath, sourceinfo.Mode())

	if err1 != nil {
		fmt.Println("Unable to create the directory structure.")
		return err1
	}

	err2 := os.Rename(itemPath, fullPath)

	if err2 != nil {
		fmt.Println("Unable to move the file.")
		return err2
	}
	return nil
}

func processPhoto(path string, info os.FileInfo) (photo Photo, err error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file: ", path)
		return
	}
	defer file.Close()

	data, err := exif.Decode(file)

	if err != nil {
		fmt.Println(err)
		return
	}
	//TODO: Replace this logic again
	// if err != nil {
	// 	return extractImageWithDimensions(path, info), nil
	// }

	tm, _ := data.DateTime()
	lat, long, _ := data.LatLong()
	widthTag, err := data.Get("PixelXDimension")
	heightTag, err := data.Get("PixelYDimension")

	var width int
	var height int

	if widthTag == nil || heightTag == nil {
		width, height = getDimensionsForImageFile(file)
	} else {
		width, _ = widthTag.Int(0)
		height, _ = heightTag.Int(0)
	}

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
		if b == a {
			return true
		}
	}
	return false
}

func getDimensionsForImageFile(file *os.File) (int, int) {
	image, _, _ := image.DecodeConfig(file)
	return image.Width, image.Height
}
