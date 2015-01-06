package main

import (
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"panorma/app"
)

var db gorm.DB

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	config := app.Config{DatabaseConnectionString: "user=jais dbname=panorma_dev sslmode=disable", LogDatabaseQueries: true}

	db = app.SetupDatabase(config)

	m.Get("/api", func(r render.Render) {
		r.HTML(200, "hello", "jeremy")
	})

	m.Get("/api/albums", func(r render.Render) {
		r.JSON(200, app.RootAlbums(db))
	})

	m.Run()
}
