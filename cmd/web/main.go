package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"reservation/internal/config"
	"reservation/internal/handlers"
	"reservation/internal/helpers"
	"reservation/internal/models"
	"reservation/internal/render"
	"time"

	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var sessionManager *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	gob.Register(models.Reservation{})

	app.IsProduction = false // change this to true for production

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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

	helpers.NewHelpers(&app)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: chiRoutes(&app),
	}
	srv.ListenAndServe()

	return nil
}
