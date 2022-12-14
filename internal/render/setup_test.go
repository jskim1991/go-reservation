package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"reservation/internal/config"
	"reservation/internal/models"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	testApp.IsProduction = false // change this to true for production
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	sessionManager = scs.New()
	sessionManager.Lifetime = 30 * time.Minute
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = testApp.IsProduction
	testApp.SessionManager = sessionManager

	app = &testApp

	os.Exit(m.Run())
}

type dummyWriter struct {
}

func (dw *dummyWriter) Header() http.Header {
	var h http.Header
	return h
}

func (dw *dummyWriter) WriteHeader(statusCode int) {

}

func (dw *dummyWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
