package main

import (
	"log"
	"net/http"
	"reservation/pkg/config"
	"reservation/pkg/handlers"
	"reservation/pkg/render"
	"time"

	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var sessionManager *scs.SessionManager

func main() {
	app.IsProduction = false
	sessionManager = scs.New()
	sessionManager.Lifetime = 30 * time.Minute
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.IsProduction
	app.SessionManager = sessionManager

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
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
}
