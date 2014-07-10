package main

import (
	"github.com/albrow/learning/peeps-martini/controllers"
	"github.com/albrow/learning/peeps-martini/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"log"
)

func main() {
	models.Init()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(func(c martini.Context, log *log.Logger) {
		log.Println("before a request")

		c.Map("poop")
		c.Next()

		log.Println("after a request")
	})

	personsController := controllers.Persons{}
	m.Post("/persons", binding.Bind(controllers.PersonForm{}), personsController.Create)
	m.Get("/persons/:id", binding.Bind(controllers.PersonQuery{}), personsController.Show)
	m.Get("/persons", binding.Bind(controllers.PersonQuery{}), personsController.Index)
	m.Delete("/persons/:id", personsController.Delete)
	m.Patch("/persons/:id", binding.Bind(controllers.PersonForm{}), personsController.Update)

	m.Get("/", func(str string) string {
		return str
	})

	m.Run()
}
