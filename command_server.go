package main

import (
	// "fmt"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
	"net/http"
)

var Render = render.New()
var DB gorm.DB

func PhotosHandler(response http.ResponseWriter, request *http.Request) {
	Render.JSON(response, http.StatusOK, map[string]string{"hello": "world"})
}

func PhotoHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	var photo Photo
	DB.Preload("SimilarPhotos").First(&photo, id)
	Render.JSON(response, http.StatusOK, photo)
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
	server.Run(":3000")
}
