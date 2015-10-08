package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
)

var (
	Render    = render.New()
	DB        gorm.DB
	AppConfig Config
)

func RunServer(c *cli.Context) {
	router := mux.NewRouter()

	var err error
	AppConfig, err = LoadConfig(c.String("config"))
	if err != nil {
		os.Exit(1)
	}

	DB = SetupDatabase(AppConfig.DatabaseConnectionString)

	router.HandleFunc("/photos", PhotosHandler)
	router.HandleFunc("/photos/{id}", PhotoHandler)
	router.HandleFunc("/similar", SimilarPhotosHandler)
	router.HandleFunc("/all_dates", AllDates)

	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":3001")
}

type JSONResponse map[string]interface{}

func AllDates(response http.ResponseWriter, request *http.Request) {
	Render.JSON(response, http.StatusOK, JSONResponse{
		"dates": AllDistinctDates(DB),
	})
}

func PhotosHandler(response http.ResponseWriter, request *http.Request) {
	page := parseIntValueForParamWithDefault(request, "page", 1)
	year := parseIntValueForParamWithDefaultZero(request, "year")
	month := parseIntValueForParamWithDefaultZero(request, "month")
	day := parseIntValueForParamWithDefaultZero(request, "day")

	photos := FindAllPhotos(DB, page, year, month, day)
	Render.JSON(response, http.StatusOK, photos)
}

func PhotoHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var photo Photo
	if DB.Preload("SimilarPhotos").First(&photo, id).RecordNotFound() {
		errorResponse := JSONResponse{"error": "No photo found for ID " + id}
		Render.JSON(response, http.StatusNotFound, errorResponse)
	} else {
		Render.JSON(response, http.StatusOK, photo)
	}
}

func SimilarPhotosHandler(response http.ResponseWriter, request *http.Request) {
	var similiarPhotos []SimilarPhoto
	DB.Find(&similiarPhotos)
	Render.JSON(response, http.StatusOK, similiarPhotos)
}

func parseIntValueForParamWithDefault(request *http.Request, param string, defaultValue int) int {
	paramString := request.URL.Query().Get(param)
	paramValue := defaultValue
	if paramString != "" {
		paramValue, _ = strconv.Atoi(paramString)
	}
	return paramValue
}

func parseIntValueForParamWithDefaultZero(request *http.Request, param string) int {
	return parseIntValueForParamWithDefault(request, param, 0)
}
