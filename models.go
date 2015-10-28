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

type Resemblance struct {
	Id                int64 `json:"id"`
	MediaId           int64 `json:"media_id"`
	ResemblingMediaId int64 `json:"resembling_media_id"`
}

type Media struct {
	Id           int64         `json:"id"`
	Path         string        `json:"path";sql:"not null;unique"`
	Name         string        `json:"name"`
	Ext          string        `json:"extension"`
	IsVideo      bool          `json:"is_video";sql:"default:'false'"`
	HashValue    int64         `json:"-";sql:"not null"`
	TakenAt      time.Time     `json:"taken_at"`
	Lat          float64       `json:"latitude"`
	Lng          float64       `json:"longitude"`
	Height       int           `json:"height"`
	Width        int           `json:"width"`
	Size         int64         `json:"size"`
	Resemblances []Resemblance `json:"resemblances,omitempty"`
	Thumbnailed  bool          `json:"-";sql:"default:'false'"`
}

func (media Media) ExistsInDatabase(db gorm.DB) bool {
	var count int
	db.Model(Media{}).Where("path = ?", media.Path).Count(&count)
	return count >= 1
}

func AllMediaForDate(db gorm.DB, page, year, month, day int) []Media {
	var media []Media
	offset := (page - 1) * ResultsPerRequest
	scope := db.Offset(offset).Limit(ResultsPerRequest)

	if year > 0 {
		scope = scope.Where("date_part('year', taken_at) = ?", year)

		if month > 0 {
			scope = scope.Where("date_part('month', taken_at) = ?", month)

			if day > 0 {
				scope = scope.Where("date_part('day', taken_at) = ?", day)
			}
		}
	}
	scope.Find(&media)
	return media
}

func AllDistinctDates(db gorm.DB) []Date {
	var results []Date
	scope := db.Table("media").Select("date_part('day', taken_at) as day, date_part('month', taken_at) as month, date_part('year', taken_at) as year")
	scope = scope.Group("day, month, year").Order("year, month, day")
	scope.Scan(&results)
	return results
}

func ImagesNotThumbnailed(db gorm.DB) []Media {
	var media []Media
	db.Where("thumbnailed = false AND is_video = false").Find(&media)
	return media
}
