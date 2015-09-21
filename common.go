package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const Thumbnails = map[string][]int{
	"small": []int{100, 100},
	"large": []int{500, 500},
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

	db.AutoMigrate(&Photo{}, &SimilarPhoto{})
	return db
}
