package controllers

import (
	"fmt"
	"github.com/albrow/learning/peeps-martini/models"
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
		r.JSON(500, map[string]interface{}{"Error": err.Error()})
	} else {
		r.JSON(200, p)
	}
}

func (Persons) Show(params martini.Params, data data.Data, r render.Render) {
	id := params["id"]
	if id == "" {
		data := map[string]interface{}{"Error": "Id cannot be empty"}
		r.JSON(500, data)
	}

	p := &models.Person{}
	if !data.KeyExists("include") {
		if err := zoom.ScanById(id, p); err != nil {
			r.JSON(500, map[string]interface{}{"Error": err.Error()})
		}
	} else {
		includes := data.GetStrings("include")
		persons := []*models.Person{}
		q := zoom.NewQuery("Person").Filter("Id =", id).Include(includes...)
		if err := q.Scan(&persons); err != nil {
			r.JSON(500, map[string]interface{}{"Error": err.Error()})
		}
		if len(persons) == 0 {
			msg := fmt.Sprintf("Could not find person with id %s", id)
			r.JSON(500, map[string]interface{}{"Error": msg})
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
			r.JSON(500, map[string]interface{}{"Error": err.Error()})
		}
	} else {
		includes := data.GetStrings("include")
		q := zoom.NewQuery("Person").Include(includes...)
		if err := q.Scan(&persons); err != nil {
			r.JSON(500, map[string]interface{}{"Error": err.Error()})
		}
	}
	r.JSON(200, persons)
}

func (Persons) Delete(params martini.Params, r render.Render) {
	id := params["id"]
	if id == "" {
		data := map[string]interface{}{"Error": "Id cannot be empty"}
		r.JSON(500, data)
	}

	if err := zoom.DeleteById("Person", id); err != nil {
		r.JSON(500, map[string]interface{}{"Error": err.Error()})
	} else {
		r.JSON(200, map[string]interface{}{"Message": "Ok"})
	}
}

func (Persons) Update(params martini.Params, data data.Data, r render.Render) {
	// Get the model by id
	id := params["id"]
	if id == "" {
		data := map[string]interface{}{"Error": "Id cannot be empty"}
		r.JSON(500, data)
	}
	p := &models.Person{}
	if err := zoom.ScanById(id, p); err != nil {
		r.JSON(500, map[string]interface{}{"Error": err.Error()})
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
		r.JSON(500, map[string]interface{}{"Error": err.Error()})
	} else {
		r.JSON(200, p)
	}
}
