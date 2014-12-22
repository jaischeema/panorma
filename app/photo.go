package app

import (
	"time"
)

type Photo struct {
	Id         int64
	Path       string `sql:not null`
	UniqueHash string `sql:"not null;unique"`
	takenAt    time.Time
	Lat        float64
	Lng        float64
	Duplicates []Duplicate
}
