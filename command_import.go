package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jinzhu/gorm"
	"panorma/bktree"
	// "github.com/rwcarlsen/goexif/exif"
	// "image"
	// _ "image/jpeg"
	// _ "image/png"
	// "os"
	// "path"
	// "path/filepath"
	// "time"
)

type Result struct {
	Id        int64
	HashValue uint64
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
		tree := bktree.New(firstPhoto.HashValue, firstPhoto.Id)
		for _, photo := range photos[1:] {
			tree.Insert(photo.HashValue, photo.Id)
		}
		return tree
	} else {
		return bktree.Node{}
	}
}

func processPhotos(db gorm.DB, sourcePath string, destinationPath string) {
	tree := createTreeFromDatabase(db)
	fmt.Println(tree.HashValue)
	fmt.Println(tree.Object.(int64))
	//
	// 	walkFunc := func(itemPath string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	//
	// 		if !info.IsDir() && strIn(filepath.Ext(itemPath), validExts) {
	// 			photo, err := processPhoto(itemPath, info)
	// 			if err != nil {
	// 				fmt.Println("Error opening file: %v", itemPath)
	// 				return nil
	// 			}
	// 			photo.Size = info.Size()
	// 			photo.UniqueHash = app.ChecksumFile(itemPath)
	//
	// 			tx := db.Begin()
	// 			if photo.ExistsInDatabase(db) {
	// 				fmt.Println("Skipping: ", itemPath)
	// 			} else if photo.IsDuplicate(db) {
	// 				fmt.Println("Duplicate: ", itemPath)
	// 				duplicatePath := calculatePhotoTimedPath(itemPath, photo.TakenAt)
	// 				dbPhoto := app.PhotoForPathAndUniqueHash(db, photo.Path, photo.UniqueHash)
	// 				var duplicate app.Duplicate
	// 				db.Where(app.Duplicate{PhotoId: dbPhoto.Id, Path: duplicatePath}).FirstOrInit(&duplicate)
	// 				if db.NewRecord(duplicate) {
	// 					db.Save(&duplicate)
	// 					err = moveFileInTransaction(itemPath, duplicatesPath, duplicate.Path)
	// 					if err != nil {
	// 						tx.Rollback()
	// 					}
	// 					fmt.Println("Duplicate saved")
	// 				} else {
	// 					fmt.Println("Skipping, duplicate already exists")
	// 				}
	// 			} else {
	// 				db.Save(&photo)
	// 				err = moveFileInTransaction(itemPath, archivePath, photo.Path)
	// 				if err != nil {
	// 					tx.Rollback()
	// 				}
	// 				fmt.Println("Move to Archive directory")
	// 			}
	// 			tx.Commit()
	// 		}
	// 		return nil
	// 	}
	//
	// 	err := filepath.Walk(path, walkFunc)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
}

//
// func moveFileInTransaction(filePath string, destinationRoot string, destinationPath string) error {
// 	sourceinfo, err := os.Stat(destinationRoot)
// 	if err != nil {
// 		fmt.Println("Destination Directory cannot be accessed, rolling back.")
// 		return err
// 	}
//
// 	fullPath := path.Join(destinationRoot, destinationPath)
// 	basePath, _ := path.Split(fullPath)
//
// 	err1 := os.MkdirAll(basePath, sourceinfo.Mode())
//
// 	if err1 != nil {
// 		fmt.Println("Unable to create the directory structure, rolling back.")
// 		return err1
// 	}
//
// 	err2 := os.Rename(filePath, fullPath)
//
// 	if err2 != nil {
// 		fmt.Println("Unable to move the file, rolling back.")
// 		return err2
// 	}
// 	return nil
// }
//
// func processPhoto(path string, info os.FileInfo) (photo app.Photo, err error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		fmt.Println("Error opening file: ", path)
// 		return
// 	}
// 	defer file.Close()
//
// 	data, err := exif.Decode(file)
// 	if err != nil {
// 		return extractImageWithDimensions(path, info), nil
// 	}
//
// 	tm, _ := data.DateTime()
// 	lat, long, _ := data.LatLong()
// 	widthTag, err := data.Get("PixelXDimension")
// 	heightTag, err := data.Get("PixelYDimension")
//
// 	if widthTag == nil || heightTag == nil {
// 		return extractImageWithDimensions(path, info), nil
// 	}
//
// 	width, _ := widthTag.Int(0)
// 	height, _ := heightTag.Int(0)
//
// 	photo = app.Photo{
// 		Path:    calculatePhotoTimedPath(path, tm),
// 		TakenAt: tm,
// 		Lat:     lat,
// 		Lng:     long,
// 		Height:  height,
// 		Width:   width,
// 	}
// 	return
// }
//
// func extractImageWithDimensions(path string, info os.FileInfo) app.Photo {
// 	width, height := imageDimensions(path)
//
// 	return app.Photo{
// 		Path:    calculatePhotoTimedPath(path, info.ModTime()),
// 		TakenAt: info.ModTime(),
// 		Height:  height,
// 		Width:   width,
// 	}
// }
//
// func calculatePhotoTimedPath(filepath string, takenAt time.Time) string {
// 	timeFormat := takenAt.Format("2006/01-January/02")
//
// 	_, file := path.Split(filepath)
// 	return path.Join(timeFormat, file)
// }
//
// func strIn(a string, list []string) bool {
// 	for _, b := range list {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// func imageDimensions(imagePath string) (int, int) {
// 	file, _ := os.Open(imagePath)
// 	defer file.Close()
// 	image, _, _ := image.DecodeConfig(file)
// 	return image.Width, image.Height
// }
