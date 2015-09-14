package main

import (
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"github.com/unrolled/render"
	"net/http"
)

var Render = render.New()

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	Render.JSON(w, http.StatusOK, map[string]string{"hello": "world"})
}

func RunServer(c *cli.Context) {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)

	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":3000")
}
