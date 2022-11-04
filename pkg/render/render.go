package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"reservation/models"
	"reservation/pkg/config"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if ok == false {
		log.Fatal("Could not get templates from cache")
	}

	td = AddDefaultData(td)

	buffer := new(bytes.Buffer)
	err := t.Execute(buffer, td)
	if err != nil {
		log.Println(err)
	}

	_, err = buffer.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page)
		templateSet, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[fileName] = templateSet
	}

	return templateCache, nil
}

func RenderTemplateWithoutCache(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
	parsedTemplate.Execute(w, nil)
}

/* 1. Cache
var templateCache = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, t string) {
	_, inMap := templateCache[t]
	if inMap == false {
		log.Println(t, "template not found in cache")
		err := createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	}

	tmpl := templateCache[t]
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache(t string) error {
	templates := []string{
		"./templates/" + t,
		"./templates/base.layout.tmpl",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	templateCache[t] = tmpl
	return nil
}
*/
