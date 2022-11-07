package forms

import (
	"net/url"
	"strings"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	form := New(url.Values{})

	isValid := form.IsValid()
	if !isValid {
		t.Error("should be valid")
	}
}

func TestForm_Required(t *testing.T) {
	form := New(url.Values{})

	form.Required("a", "b", "c")
	if form.IsValid() {
		t.Error("should be invalid")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.IsValid() {
		t.Error("should be valid")
	}
}

func TestForm_Has(t *testing.T) {
	form := New(url.Values{})

	has := form.Has("a")
	if has {
		t.Error("should be false")
	}

	postData := url.Values{}
	postData.Add("a", "some great value")
	form = New(postData)

	has = form.Has("a")
	if !has {
		t.Error("should be true")
	}
}

func TestForm_MinLength(t *testing.T) {
	invalidPostData := url.Values{}
	invalidPostData.Add("a", "")
	form := New(invalidPostData)

	valid := form.MinLength("a", 1)
	if valid {
		t.Error("should be invalid")
	}

	errMsg := form.Errors.Get("a")
	if strings.TrimSpace(errMsg) == "" {
		t.Error("should have an error message")
	}

	validPostData := url.Values{}
	validPostData.Add("a", "a")
	form = New(validPostData)

	valid = form.MinLength("a", 1)
	if !valid {
		t.Error("should be valid")
	}
}

func TestForm_IsEmail(t *testing.T) {
	form := New(url.Values{})

	form.IsEmail("x")
	if form.IsValid() {
		t.Error("should be invalid")
	}

	postData := url.Values{}
	postData.Add("x", "me")
	form = New(postData)
	form.IsEmail("x")
	if form.IsValid() {
		t.Error("should be invalid")
	}

	postData = url.Values{}
	postData.Add("x", "me@email.com")
	form = New(postData)
	form.IsEmail("x")
	if form.IsValid() == false {
		t.Error("should be valid")
	}
}
