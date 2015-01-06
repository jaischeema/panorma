package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"panorma/app"
	"path"
	"path/filepath"
	"time"
)

var validExts = []string{".jpg", ".jpeg", ".tiff", ".tif", ".gif", ".png", ".JPG", ".mov", ".m4v", ".3gp"}

func main() {
	config := app.Config{
		DatabaseConnectionString: "user=jais dbname=panorma_dev sslmode=disable",
		LogDatabaseQueries:       false,
	}

	db := app.SetupDatabase(config)
	processPhotos(db, "/Users/jais/Desktop/Images/", "/Users/jais/Archive/Sorted/", "/Users/jais/Archive/Duplicates/")
}

func processPhotos(db gorm.DB, path string, archivePath string, duplicatesPath string) {
	walkFunc := func(itemPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strIn(filepath.Ext(itemPath), validExts) {
			photo, err := processPhoto(itemPath, info)
			if err != nil {
				fmt.Println("Error opening file: %v", itemPath)
				return nil
			}
			photo.Size = info.Size()
			photo.UniqueHash = app.ChecksumFile(itemPath)

			tx := db.Begin()
			if photo.ExistsInDatabase(db) {
				fmt.Println("Skipping: ", itemPath)
			} else if photo.IsDuplicate(db) {
				fmt.Println("Duplicate: ", itemPath)
				duplicatePath := calculatePhotoTimedPath(itemPath, photo.TakenAt)
				dbPhoto := app.PhotoForPathAndUniqueHash(db, photo.Path, photo.UniqueHash)
				var duplicate app.Duplicate
				db.Where(app.Duplicate{PhotoId: dbPhoto.Id, Path: duplicatePath}).FirstOrInit(&duplicate)
				if db.NewRecord(duplicate) {
					db.Save(&duplicate)
					err = moveFileInTransaction(itemPath, duplicatesPath, duplicate.Path)
					if err != nil {
						tx.Rollback()
					}
					fmt.Println("Duplicate saved")
				} else {
					fmt.Println("Skipping, duplicate already exists")
				}
			} else {
				db.Save(&photo)
				err = moveFileInTransaction(itemPath, archivePath, photo.Path)
				if err != nil {
					tx.Rollback()
				}
				fmt.Println("Move to Archive directory")
			}
			tx.Commit()
		}
		return nil
	}

	err := filepath.Walk(path, walkFunc)
	if err != nil {
		fmt.Println(err)
	}
}

func moveFileInTransaction(filePath string, destinationRoot string, destinationPath string) error {
	sourceinfo, err := os.Stat(destinationRoot)
	if err != nil {
		fmt.Println("Destination Directory cannot be accessed, rolling back.")
		return err
	}

	fullPath := path.Join(destinationRoot, destinationPath)
	basePath, _ := path.Split(fullPath)

	err1 := os.MkdirAll(basePath, sourceinfo.Mode())

	if err1 != nil {
		fmt.Println("Unable to create the directory structure, rolling back.")
		return err1
	}

	err2 := os.Rename(filePath, fullPath)

	if err2 != nil {
		fmt.Println("Unable to move the file, rolling back.")
		return err2
	}
	return nil
}

func processPhoto(path string, info os.FileInfo) (photo app.Photo, err error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file: ", path)
		return
	}
	defer file.Close()

	data, err := exif.Decode(file)
	if err != nil {
		return extractImageWithDimensions(path, info), nil
	}

	tm, _ := data.DateTime()
	lat, long, _ := data.LatLong()
	widthTag, err := data.Get("PixelXDimension")
	heightTag, err := data.Get("PixelYDimension")

	if widthTag == nil || heightTag == nil {
		return extractImageWithDimensions(path, info), nil
	}

	width, _ := widthTag.Int(0)
	height, _ := heightTag.Int(0)

	photo = app.Photo{
		Path:    calculatePhotoTimedPath(path, tm),
		TakenAt: tm,
		Lat:     lat,
		Lng:     long,
		Height:  height,
		Width:   width,
	}
	return
}

func extractImageWithDimensions(path string, info os.FileInfo) app.Photo {
	width, height := imageDimensions(path)

	return app.Photo{
		Path:    calculatePhotoTimedPath(path, info.ModTime()),
		TakenAt: info.ModTime(),
		Height:  height,
		Width:   width,
	}
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

func imageDimensions(imagePath string) (int, int) {
	file, _ := os.Open(imagePath)
	defer file.Close()
	image, _, _ := image.DecodeConfig(file)
	return image.Width, image.Height
}
