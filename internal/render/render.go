package render

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"reservation/internal/config"
	"reservation/internal/models"

	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var templatePath = "./templates"

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.SessionManager.PopString(r.Context(), "flash")
	td.Error = app.SessionManager.PopString(r.Context(), "error")
	td.Warning = app.SessionManager.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if ok == false {
		log.Println("Could not get templates from cache")
		return errors.New("Can't get templates from cache")
	}

	td = AddDefaultData(td, r)

	buffer := new(bytes.Buffer)
	err := t.Execute(buffer, td)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = buffer.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob(templatePath + "/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page)
		templateSet, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		matches, err := filepath.Glob(templatePath + "/*.layout.tmpl")
		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(templatePath + "/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[fileName] = templateSet
	}

	return templateCache, nil
}

func RenderTemplateWithoutCache(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles(templatePath+"/"+tmpl, templatePath+"/base.layout.tmpl")
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
