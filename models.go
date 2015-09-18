package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

const ResultsPerRequest = 20

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

func FindAllYears() []int {
	years := []int{}
	selectString := "date_part('year', taken_at) as year"
	rows, _ := DB.Table("photos").Select(selectString).Group("year").Order("year").Rows()
	for rows.Next() {
		var year int
		rows.Scan(&year)
		years = append(years, year)
	}
	return years
}

func FindAllMonths(year int) []int {
	months := []int{}
	db := DB.Table("photos").Select("date_part('month', taken_at) as month")
	db = db.Group("month").Order("month")
	db = db.Where("date_part('year', taken_at) = ?", year)
	rows, _ := db.Rows()
	for rows.Next() {
		var month int
		rows.Scan(&month)
		months = append(months, month)
	}
	return months
}

func FindAllDays(year int, month int) []int {
	days := []int{}
	db := DB.Table("photos").Select("date_part('day', taken_at) as day")
	db = db.Group("day").Order("day")
	db = db.Where("date_part('year', taken_at) = ?", year)
	db = db.Where("date_part('month', taken_at) = ?", month)
	rows, _ := db.Rows()
	for rows.Next() {
		var day int
		rows.Scan(&day)
		days = append(days, day)
	}
	return days
}
