package main

import (
	"net/http"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
)

var Render = render.New()
var DB gorm.DB

func RunServer(c *cli.Context) {
	router := mux.NewRouter()
	DB = SetupDatabase(c)

	router.HandleFunc("/photos", PhotosHandler)
	router.HandleFunc("/photos/{id}", PhotoHandler)
	router.HandleFunc("/similar", SimilarPhotosHandler)
	router.HandleFunc("/intervals", Intervals)

	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":3001")
}

type JSONResponse map[string]interface{}

func Intervals(response http.ResponseWriter, request *http.Request) {
	year := parseIntValueForParamWithDefaultZero(request, "year")
	month := parseIntValueForParamWithDefaultZero(request, "month")

	if year == 0 {
		Render.JSON(response, http.StatusOK, JSONResponse{"years": FindAllYears()})
	} else {
		if month == 0 {
			Render.JSON(response, http.StatusOK, JSONResponse{
				"months": FindAllMonths(year),
				"year":   year,
			})
		} else {
			Render.JSON(response, http.StatusOK, JSONResponse{
				"days":  FindAllDays(year, month),
				"month": month,
				"year":  year,
			})
		}
	}
}

func PhotosHandler(response http.ResponseWriter, request *http.Request) {
	page := parseIntValueForParamWithDefault(request, "page", 1)
	year := parseIntValueForParamWithDefaultZero(request, "year")
	month := parseIntValueForParamWithDefaultZero(request, "month")
	day := parseIntValueForParamWithDefaultZero(request, "day")

	photos := FindAllPhotos(page, year, month, day)
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
