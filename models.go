package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type SimilarPhoto struct {
	Id             int64 `json:"id"`
	PhotoId        int64 `json:"photo_id"`
	SimilarPhotoId int64 `json:"similar_photo_id"`
}

type Photo struct {
	Id            int64          `json:"id"`
	Path          string         `json:"path";sql:"not null;unique"`
	HashValue     int64          `json:"hash_value";sql:"not null"`
	TakenAt       time.Time      `json:"taken_at"`
	Lat           float64        `json:"latitude"`
	Lng           float64        `json:"longitude"`
	Height        int            `json:"height"`
	Width         int            `json:"width"`
	Size          int64          `json:"size"`
	SimilarPhotos []SimilarPhoto `json:"similar_photos,omitempty"`
}

func (photo Photo) ExistsInDatabase(db gorm.DB) bool {
	var count int
	db.Model(Photo{}).Where("path = ?", photo.Path).Count(&count)
	return count >= 1
}
