package helpers

import (
	"fmt"
	"net/http"
	"reservation/internal/config"
	"runtime/debug"
	"time"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func FormatStringToTime(s string) (time.Time, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, s)
	return date, err
}

func FormatTimeToString(d time.Time) string {
	return d.Format("2006-01-02")
}
