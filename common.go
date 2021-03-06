package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const ThumbnailExtension = ".jpg"

type ThumbnailSize struct {
	Width  int
	Height int
}

var ThumbnailSizes = map[string]ThumbnailSize{
	"small":  ThumbnailSize{100, 100},
	"medium": ThumbnailSize{300, 200},
	"large":  ThumbnailSize{1200, 666},
}

type Config struct {
	SourceFolderPath         string `json:"source_folder_path"`
	DestinationFolderPath    string `json:"destination_folder_path"`
	ThumbnailsFolderPath     string `json:"thumbnails_folder_path"`
	DatabaseConnectionString string `json:"database_connection_string"`
}

func LoadConfig(configPath string) (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	return config, err
}

func SetupDatabase(connectionString string) gorm.DB {
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic("Unable to open database")
	}

	db.AutoMigrate(&Media{}, &Resemblance{})
	return db
}

func PartitionIdAsPath(input int64) string {
	inputRunes := []rune(fmt.Sprintf("%09d", input))
	return path.Join(
		string(inputRunes[0:3]),
		string(inputRunes[3:6]),
		string(inputRunes[6:]),
	)
}
