package main

import (
	"net/http"
	"os"
	"path"
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

	router.HandleFunc("/media", MediaIndexHandler)
	router.HandleFunc("/media/{id}", MediaHandler)
	router.HandleFunc("/media/{id}/{thumb}", ThumbnailHandler)
	router.HandleFunc("/similar", ResemblancesHandler)
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

func MediaIndexHandler(response http.ResponseWriter, request *http.Request) {
	page := parseIntValueForParamWithDefault(request, "page", 1)
	year := parseIntValueForParamWithDefaultZero(request, "year")
	month := parseIntValueForParamWithDefaultZero(request, "month")
	day := parseIntValueForParamWithDefaultZero(request, "day")

	media := AllMediaForDate(DB, page, year, month, day)
	Render.JSON(response, http.StatusOK, media)
}

func MediaHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var media Media
	if DB.Preload("Resemblance").First(&media, id).RecordNotFound() {
		errorResponse := JSONResponse{"error": "No media found for ID " + id}
		Render.JSON(response, http.StatusNotFound, errorResponse)
	} else {
		Render.JSON(response, http.StatusOK, media)
	}
}

func ThumbnailHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	thumb := vars["thumb"]

	var media Media
	if DB.First(&media, id).RecordNotFound() {
		errorResponse := JSONResponse{"error": "No media found for ID " + id}
		Render.JSON(response, http.StatusNotFound, errorResponse)
	} else {
		thumbFile := path.Join(
			AppConfig.ThumbnailsFolderPath,
			PartitionIdAsPath(media.Id),
			thumb+ThumbnailExtension,
		)
		http.ServeFile(response, request, thumbFile)
	}
}

func ResemblancesHandler(response http.ResponseWriter, request *http.Request) {
	var resemblances []Resemblance
	DB.Find(&resemblances)
	Render.JSON(response, http.StatusOK, resemblances)
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
