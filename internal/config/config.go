package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache       bool
	TemplateCache  map[string]*template.Template
	IsProduction   bool
	SessionManager *scs.SessionManager
}
