package controllers

import (
	"../models"
	"errors"
	"fmt"
	data "github.com/albrow/martini-data"
	"github.com/albrow/zoom"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type Persons struct{}

func (Persons) Create(data data.Data, r render.Render) {
	p := &models.Person{
		Name: data.Get("name"),
		Age:  data.GetInt("age"),
	}

	if err := zoom.Save(p); err != nil {
		panic(err)
	} else {
		r.JSON(200, p)
	}
}

func (Persons) Show(params martini.Params, data data.Data, r render.Render) {
	id := params["id"]
	if id == "" {
		err := errors.New("Id cannot be empty.")
		r.JSON(400, newJSONError("invalidParameters", err))
	}

	p := &models.Person{}
	if !data.KeyExists("include") {
		if err := zoom.ScanById(id, p); err != nil {
			if _, ok := err.(*zoom.KeyNotFoundError); ok {
				err := fmt.Errorf("Could not find person with id %s", id)
				r.JSON(400, newJSONError("invalidParameters", err))
				return
			} else {
				panic(err)
			}
		}
	} else {
		includes := data.GetStrings("include")
		persons := []*models.Person{}
		q := zoom.NewQuery("Person").Filter("Id =", id).Include(includes...)
		if err := q.Scan(&persons); err != nil {
			panic(err)
		}
		if len(persons) == 0 {
			err := fmt.Errorf("Could not find person with id %s", id)
			r.JSON(400, newJSONError("invalidParameters", err))
		} else {
			p = persons[0]
		}
	}
	r.JSON(200, p)
}

func (Persons) Index(data data.Data, r render.Render) {
	persons := []*models.Person{}
	if !data.KeyExists("include") {
		if err := zoom.NewQuery("Person").Scan(&persons); err != nil {
			panic(err)
		}
	} else {
		includes := data.GetStrings("include")
		q := zoom.NewQuery("Person").Include(includes...)
		if err := q.Scan(&persons); err != nil {
			panic(err)
		}
	}
	r.JSON(200, persons)
}

func (Persons) Delete(params martini.Params, r render.Render) {
	id := params["id"]
	if id == "" {
		err := errors.New("Id cannot be empty.")
		r.JSON(400, newJSONError("invalidParameters", err))
	}

	if err := zoom.DeleteById("Person", id); err != nil {
		panic(err)
	} else {
		r.JSON(200, newJSONOk())
	}
}

func (Persons) Update(params martini.Params, data data.Data, r render.Render) {
	// Get the model by id
	id := params["id"]
	if id == "" {
		err := errors.New("Id cannot be empty.")
		r.JSON(400, newJSONError("invalidParameters", err))
	}
	p := &models.Person{}
	if err := zoom.ScanById(id, p); err != nil {
		panic(err)
	}

	// Update person model
	if data.KeyExists("name") {
		p.Name = data.Get("name")
	}
	if data.KeyExists("age") {
		p.Age = data.GetInt("age")
	}

	// Save the model and render the result
	if err := zoom.Save(p); err != nil {
		panic(err)
	} else {
		r.JSON(200, p)
	}
}
