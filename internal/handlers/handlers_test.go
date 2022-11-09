package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reservation/internal/models"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	url                string
	method             string
	expectedStatusCode int
}{
	{
		url:                "/",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		url:                "/about",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		url:                "/generals-quarters",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		url:                "/majors-suite",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		url:                "/search-availability",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		url:                "/contact",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	// {
	// 	name:               "make reservation",
	// 	url:                "/make-reservation",
	// 	method:             "GET",
	// 	params:             []postData{},
	// 	expectedStatusCode: http.StatusOK,
	// },
	// {
	// 	name:               "search availability",
	// 	url:                "/search-availability",
	// 	method:             "POST",
	// 	params:             []postData{{key: "start", value: "2022-01-01"}, {key: "end", value: "2022-01-02"}},
	// 	expectedStatusCode: http.StatusOK,
	// },
	// {
	// 	name:               "search availability json",
	// 	url:                "/search-availability-json",
	// 	method:             "POST",
	// 	params:             []postData{{key: "start", value: "2022-01-01"}, {key: "end", value: "2022-01-02"}},
	// 	expectedStatusCode: http.StatusOK,
	// },
	// {
	// 	url:                "/make-reservation",
	// 	method:             "POST",
	// 	params:             []postData{{key: "first-name", value: "John"}, {key: "last-name", value: "Sam"}, {key: "email", value: "js@email.com"}, {key: "phone", value: "010-0000-0000"}},
	// 	expectedStatusCode: http.StatusOK,
	// },
}

func TestHandlers_GET(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range tests {
		response, err := testServer.Client().Get(testServer.URL + e.url)
		if err != nil {
			t.Error(err)
		}
		if response.StatusCode != e.expectedStatusCode {
			t.Errorf("expected %d but got %d for %s", e.expectedStatusCode, response.StatusCode, e.method+" "+e.url)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	sessionManager.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong status code: expected %d but got %d", http.StatusOK, rr.Code)
	}
}

func TestRepository_Reservation_NoSession(t *testing.T) {
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong status code: expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_Reservation_NoRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 99,
		Room: models.Room{
			ID:       99,
			RoomName: "Someone's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	sessionManager.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong status code: expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	requestBody := "start_date=2099-12-30"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=2099-12-31")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first-name=John")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last-name=Sam")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=js@domain.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=111-1111-1111")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	sessionManager.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong status code: expected %d but got %d", http.StatusSeeOther, rr.Code)
	}
}

func TestRepository_PostReservation_NoRequestBody(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("POST", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	sessionManager.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong status code: expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_PostReservation_InsertReservationError(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 99,
		Room: models.Room{
			ID:       99,
			RoomName: "General's Quarters",
		},
	}

	postedData := url.Values{}
	postedData.Add("start_date", "2099-12-30")
	postedData.Add("end_date", "2099-12-31")
	postedData.Add("first-name", "John")
	postedData.Add("last-name", "Sam")
	postedData.Add("email", "js@domain.com")
	postedData.Add("phone", "111-1111-1111")
	postedData.Add("room_id", "99")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	sessionManager.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong status code: expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_PostReservation_InsertRoomRestrictionError(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 999,
		Room: models.Room{
			ID:       999,
			RoomName: "General's Quarters",
		},
	}

	requestBody := "start_date=2099-12-30"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=2099-12-31")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first-name=John")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last-name=Sam")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=js@domain.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=111-1111-1111")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=999")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	sessionManager.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong status code: expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	requestBody := "start=2099-12-30"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end=2099-12-31")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(requestBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var jsonOutput jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &jsonOutput)
	if err != nil {
		t.Error("Failed to parse json")
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := sessionManager.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
