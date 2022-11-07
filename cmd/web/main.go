package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"reservation/internal/config"
	"reservation/internal/handlers"
	"reservation/internal/models"
	"reservation/internal/render"
	"time"

	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var sessionManager *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	gob.Register(models.Reservation{})

	app.IsProduction = false // change this to true for production
	sessionManager = scs.New()
	sessionManager.Lifetime = 30 * time.Minute
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.IsProduction
	app.SessionManager = sessionManager

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: chiRoutes(&app),
	}
	srv.ListenAndServe()

	return nil
}