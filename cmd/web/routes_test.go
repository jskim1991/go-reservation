package main

import (
	"fmt"
	"reservation/internal/config"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var config config.AppConfig

	mux := chiRoutes(&config)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Error(fmt.Printf("expected type *chi.Mux but found %T", v))
	}
}
