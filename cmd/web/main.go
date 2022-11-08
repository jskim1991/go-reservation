package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"reservation/internal/config"
	"reservation/internal/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: chiRoutes(&app),
	}
	srv.ListenAndServe()
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

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

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=")
	if err != nil {
		log.Fatal("Cannot connect to DB")
		return nil, err
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)


	return db, nil
}