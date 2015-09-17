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

func PhotosHandler(response http.ResponseWriter, request *http.Request) {
	pageString := request.URL.Query().Get("page")
	page := 1
	if pageString != "" {
		page, _ = strconv.Atoi(pageString)
	}
	var photos []Photo
	offset := (page - 1) * 20
	DB.Offset(offset).Limit(20).Find(&photos)
	Render.JSON(response, http.StatusOK, photos)
}

func PhotoHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var photo Photo
	if DB.Preload("SimilarPhotos").First(&photo, id).RecordNotFound() {
		errorResponse := map[string]string{"error": "No photo found for ID " + id}
		Render.JSON(response, http.StatusNotFound, errorResponse)
	} else {
		Render.JSON(response, http.StatusOK, photo)
	}
}

func SimilarPhotosHandler(response http.ResponseWriter, request *http.Request) {
	Render.JSON(response, http.StatusOK, map[string]string{"hello": "world"})
}

func RunServer(c *cli.Context) {
	router := mux.NewRouter()
	DB = SetupDatabase(c)

	router.HandleFunc("/photos", PhotosHandler)
	router.HandleFunc("/photos/{id}", PhotoHandler)
	router.HandleFunc("/similar", SimilarPhotosHandler)

	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":3001")
}
