package main

import (
	"github.com/albrow/learning/peeps-martini/controllers"
	"github.com/albrow/learning/peeps-martini/models"
	data "github.com/albrow/martini-data"
	recovery "github.com/albrow/martini-json-recovery"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type CustomMartini struct {
	*martini.Martini
	martini.Router
}

func main() {
	models.Init()

	m := CustomMartini{
		Martini: martini.New(),
	}
	r := martini.NewRouter()
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	m.Router = r
	m.Use(martini.Logger())
	m.Use(render.Renderer())
	m.Use(data.Parser())
	m.Use(recovery.JSONRecovery())

	personsController := controllers.Persons{}
	m.Post("/persons", personsController.Create)
	m.Get("/persons/:id", personsController.Show)
	m.Get("/persons", personsController.Index)
	m.Delete("/persons/:id", personsController.Delete)
	m.Patch("/persons/:id", personsController.Update)

	m.Run()
}
