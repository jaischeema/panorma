package main

import (
	"github.com/codegangsta/cli"
	// "github.com/go-martini/martini"
	// _ "github.com/jinzhu/gorm"
	// "github.com/martini-contrib/render"
)

func RunServer(c *cli.Context) {
	// m := martini.Classic()
	// m.Use(render.Renderer())
	//
	// config := app.Config{DatabaseConnectionString: "user=jais dbname=panorma_dev sslmode=disable", LogDatabaseQueries: true}
	//
	// db = app.SetupDatabase(config)
	//
	// m.Get("/api/albums", func(r render.Render) {
	// 	r.JSON(200, app.RootAlbums(db))
	// })

	//api/photos?year=(&month=(&day=))
	// GET /api/albums/:album_id -> { albums, photos }
	// PST /api/albums
	// GET /api/photos/:photo_id -> { photo, duplicates }
	// PUT /api/albums/:album_id
	// GET /api/albums/:album_id/add/:photo_id
	// GET /api/albums/:album_id/remove/:photo_id
	// GET /api/duplicates
	// GET /api/duplicates/:duplicate_id
	// DEL /api/duplicates/:duplicate_id

	// m.Run()
}
