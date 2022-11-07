package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	dummyHandler := dummyHandler{}
	handler := NoSurf(&dummyHandler)

	switch v := handler.(type) {
	case http.Handler:
	default:
		t.Error(fmt.Printf("expected type http.Handler but found %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	dummyHandler := dummyHandler{}
	handler := SessionLoad(&dummyHandler)

	switch v := handler.(type) {
	case http.Handler:
	default:
		t.Error(fmt.Printf("expected type http.Handler but found %T", v))
	}
}
