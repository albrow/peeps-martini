package controllers

import (
	"fmt"
	"github.com/albrow/learning/peeps-martini/models"
	"github.com/albrow/zoom"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"strings"
)

type Persons struct{}

type PersonForm struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}

type PersonQuery struct {
	Include string `form:"include"`
}

func (Persons) Create(data PersonForm, r render.Render) {
	p := &models.Person{
		Name: data.Name,
		Age:  data.Age,
	}

	if err := zoom.Save(p); err != nil {
		r.JSON(500, map[string]interface{}{"Error": err.Error()})
	} else {
		r.JSON(200, p)
	}
}

func (Persons) Show(params martini.Params, data PersonQuery, r render.Render) {
	id := params["id"]
	if id == "" {
		data := map[string]interface{}{"Error": "Id cannot be empty"}
		r.JSON(500, data)
	}

	p := &models.Person{}
	if data.Include != "" {
		if err := zoom.ScanById(id, p); err != nil {
			r.JSON(500, map[string]interface{}{"Error": err.Error()})
		}
	} else {
		includes := strings.Split(data.Include, ",")
		q := zoom.NewQuery("Person").Filter("Id =", id).Include(includes...)
		persons := []*models.Person{}
		if err := q.Scan(persons); err != nil {
			r.JSON(500, map[string]interface{}{"Error": err.Error()})
		} else if len(persons) == 0 {
			msg := fmt.Sprintf("Could not find person with id %s", id)
			r.JSON(500, map[string]interface{}{"Error": msg})
		} else {
			p = persons[0]
		}
	}
	r.JSON(200, p)
}

func (Persons) Index(r render.Render) {
	persons := []*models.Person{}
	if err := zoom.NewQuery("Person").Scan(&persons); err != nil {
		r.JSON(500, map[string]interface{}{"Error": err.Error()})
	} else {
		r.JSON(200, persons)
	}
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

func (Persons) Update(params martini.Params, r render.Render, req *http.Request) {
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

	// Parse form data and update person model
	data, err := parseFormData(req)
	if err != nil {
		r.JSON(500, map[string]interface{}{"Error": err.Error()})
	}
	if name, found := data["name"]; found {
		p.Name = name
	}
	if ageStr, found := data["age"]; found {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			r.JSON(500, map[string]interface{}{"Error": err.Error()})
		}
		p.Age = age
	}

	// Save the model and render the result
	if err := zoom.Save(p); err != nil {
		r.JSON(500, map[string]interface{}{"Error": err.Error()})
	} else {
		r.JSON(200, p)
	}
}

// Converts multipart or url encoded form data to map[string]string
// Useful for solving the "update problem." The problem is that normally,
// martini-contrib/binding converts fields which aren't included to the
// zero value, which causes those fields in the model to be overwritten
// and "blanked out" The desired behavior is that if a field is not included,
// we should retain the original value for that field.
func parseFormData(req *http.Request) (map[string]string, error) {
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(2048); err != nil {
			return nil, err
		}
		// convert from map[string][]string to map[string]string
		// the first value for each key is considered the right value
		// I can't even think of a case where there will be more than
		// one item per key.
		result := map[string]string{}
		for key, val := range req.MultipartForm.Value {
			result[key] = val[0]
		}
		return result, nil
	} else if strings.Contains(contentType, "form-urlencoded") {
		if err := req.ParseForm(); err != nil {
			return nil, err
		}
		// convert from map[string][]string to map[string]string
		// the first value for each key is considered the right value
		// I can't even think of a case where there will be more than
		// one item per key.
		result := map[string]string{}
		for key, val := range req.PostForm {
			result[key] = val[0]
		}
		return result, nil
	}
	return nil, fmt.Errorf("Unrecognized content-type: %s. Only multipart/form-data or form-urlencoded are accepted", contentType)
}
