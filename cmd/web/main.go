package main

import (
	"log"
	"net/http"
	"reservation/pkg/config"
	"reservation/pkg/handlers"
	"reservation/pkg/render"
)

func main() {
	var config config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	config.TemplateCache = tc
	config.UseCache = false

	repo := handlers.NewRepo(&config)
	handlers.NewHandlers(repo)

	render.NewTemplates(&config)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	// http.ListenAndServe(":8080", nil)

	/* Routing using pat */
	srv := &http.Server{
		Addr:    ":8080",
		Handler: chiRoutes(&config),
	}
	srv.ListenAndServe()
}
