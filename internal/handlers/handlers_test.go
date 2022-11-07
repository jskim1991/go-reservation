package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{
		name:               "home",
		url:                "/",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "about",
		url:                "/about",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "general's quarters",
		url:                "/generals-quarters",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "major's suite",
		url:                "/majors-suite",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "search availability",
		url:                "/search-availability",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "contact",
		url:                "/contact",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "make reservation",
		url:                "/make-reservation",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "search availability",
		url:                "/search-availability",
		method:             "POST",
		params:             []postData{{key: "start", value: "2022-01-01"}, {key: "end", value: "2022-01-02"}},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "search availability json",
		url:                "/search-availability-json",
		method:             "POST",
		params:             []postData{{key: "start", value: "2022-01-01"}, {key: "end", value: "2022-01-02"}},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "make reservation",
		url:                "/make-reservation",
		method:             "POST",
		params:             []postData{{key: "first-name", value: "John"}, {key: "last-name", value: "Sam"}, {key: "email", value: "js@email.com"}, {key: "phone", value: "010-0000-0000"}},
		expectedStatusCode: http.StatusOK,
	},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range tests {
		if e.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Error(err)
			}
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("expected %d but got %d for %s", e.expectedStatusCode, response.StatusCode, e.method+" "+e.url)
			}
		} else {
			values := url.Values{}
			for _, param := range e.params {
				values.Add(param.key, param.value)
			}
			response, err := testServer.Client().PostForm(testServer.URL+e.url, values)
			if err != nil {
				t.Error(err)
			}
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("expected %d but got %d for %s", e.expectedStatusCode, response.StatusCode, e.method+" "+e.url)
			}
		}
	}
}
