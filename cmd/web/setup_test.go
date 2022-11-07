package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type dummyHandler struct {
}

func (mh *dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
