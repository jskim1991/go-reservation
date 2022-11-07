package handlers

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reservation/internal/config"
	"reservation/internal/models"
	"reservation/internal/render"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var sessionManager *scs.SessionManager
var templatePath = "./../../templates"

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	app.IsProduction = false // change this to true for production

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	sessionManager = scs.New()
	sessionManager.Lifetime = 30 * time.Minute
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.IsProduction
	app.SessionManager = sessionManager

	tc, err := testCreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/reservation-summary", Repo.ReservationSummary)

	staticFileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))
	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.IsProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}

func testCreateTemplateCache() (map[string]*template.Template, error) {
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
