package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type SimilarPhoto struct {
	Id             int64
	PhotoId        int64
	SimilarPhotoId int64
}

type Photo struct {
	Id            int64
	Path          string `sql:"not null;unique"`
	HashValue     string `sql:"not null"`
	TakenAt       time.Time
	Lat           float64
	Lng           float64
	Height        int
	Width         int
	Size          int64
	SimilarPhotos []SimilarPhoto
}

func (photo Photo) ExistsInDatabase(db gorm.DB) bool {
	var count int
	db.Model(Photo{}).Where("path = ?", photo.Path).Count(&count)
	return count >= 1
}
