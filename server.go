package main

import (
	"github.com/albrow/learning/peeps-martini/controllers"
	"github.com/albrow/learning/peeps-martini/models"
	data "github.com/albrow/martini-data"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func main() {
	models.Init()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(data.Parser())

	personsController := controllers.Persons{}
	m.Post("/persons", personsController.Create)
	m.Get("/persons/:id", personsController.Show)
	m.Get("/persons", personsController.Index)
	m.Delete("/persons/:id", personsController.Delete)
	m.Patch("/persons/:id", personsController.Update)

	m.Run()
}
