package main

import (
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
	"panorma/app"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	db, _ := gorm.Open("postgres", "user=jais dbname=panorma_dev sslmode=disable")
	db.AutoMigrate(&app.Photo{}, &app.Album{}, &app.Duplicate{})
	db.LogMode(true)

	m.Get("/api", func(r render.Render) {
		r.HTML(200, "hello", "jeremy")
	})

	m.Get("/api/albums", func(r render.Render) {
		r.JSON(200, app.RootAlbums(db))
	})

	m.Run()
}
