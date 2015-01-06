package app

import (
	"crypto/sha1"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"math"
	"os"
	"time"
)

type Photo struct {
	Id         int64
	Path       string `sql:"not null;unique"`
	UniqueHash string `sql:"not null;unique"`
	TakenAt    time.Time
	Lat        float64
	Lng        float64
	Height     int
	Width      int
	Size       int64
	Duplicates []Duplicate
}

func (photo Photo) AddToYearlyAlbum() {
	// CreateOrFindAlbumForYear(photo.takenAt.year)
	// CreateOrFindAlbumForYearAndMonth(photo.takenAt.month)
	//
}

func (photo Photo) ExistsInDatabase(db gorm.DB) bool {
	var count int
	db.Model(Photo{}).Where("unique_hash = ? AND path = ?", photo.UniqueHash, photo.Path).Count(&count)
	return count >= 1
}

func (photo Photo) IsDuplicate(db gorm.DB) bool {
	var count int
	db.Model(Photo{}).Where("unique_hash = ? AND path != ?", photo.UniqueHash, photo.Path).Count(&count)
	return count >= 1
}

func PhotoForPathAndUniqueHash(db gorm.DB, path string, uniqueHash string) Photo {
	var photo Photo
	db.Where("unique_hash = ? AND path != ?", uniqueHash, path).First(&photo)
	return photo
}

const fileChunk = 8192

func ChecksumFile(path string) string {
	fh, err := os.Open(path)

	if err != nil {
		panic(err.Error())
	}

	defer fh.Close()
	stat, _ := fh.Stat()
	size := stat.Size()
	chunks := uint64(math.Ceil(float64(size) / float64(fileChunk)))
	h := sha1.New()

	for i := uint64(0); i < chunks; i++ {
		csize := int(math.Min(fileChunk, float64(size-int64(i*fileChunk))))
		buf := make([]byte, csize)

		fh.Read(buf)
		io.WriteString(h, string(buf))
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
