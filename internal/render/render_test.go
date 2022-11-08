package render

import (
	"net/http"
	"reservation/internal/models"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	sessionManager.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result == nil {
		t.Error("failed")
	}
	if result.Flash != "123" {
		t.Error("flash value not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	templatePath = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	ww := dummyWriter{}

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error(err)
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	templatePath = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/dummy", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = sessionManager.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}
