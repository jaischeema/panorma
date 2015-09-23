package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

const ResultsPerRequest = 20

type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

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

func FindAllPhotos(page, year, month, day int) []Photo {
	var photos []Photo
	offset := (page - 1) * ResultsPerRequest
	db := DB.Offset(offset).Limit(ResultsPerRequest)

	if year > 0 {
		db = db.Where("date_part('year', taken_at) = ?", year)

		if month > 0 {
			db = db.Where("date_part('month', taken_at) = ?", month)

			if day > 0 {
				db = db.Where("date_part('day', taken_at) = ?", day)
			}
		}
	}
	db.Find(&photos)
	return photos
}

func AllDistinctDates() []Date {
	var results []Date
	db := DB.Table("photos").Select("date_part('day', taken_at) as day, date_part('month', taken_at) as month, date_part('year', taken_at) as year")
	db = db.Group("day, month, year").Order("year, month, day")
	db.Scan(&results)
	return results
}
