package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{Name: "cookieName", Value: "expectedValue"})

	value, err := Get(req, "cookieName")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedValue := "expectedValue"
	if value != expectedValue {
		t.Errorf("Expected value %s, but got %s", expectedValue, value)
	}
}

func TestGetSigned(t *testing.T) {
	DefaultManager = signedManager

	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	expectedValue := "myValue"

	cookieValue := signCookieValue(expectedValue, signedManager.signingKey)

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	value, err := GetSigned(r, cookieName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != expectedValue {
		t.Errorf("Expected value '%s', but got '%s'", expectedValue, value)
	}
}

func TestSet(t *testing.T) {
	_, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	err = Set(rr, "cookieName", "cookieValue")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cookie := rr.Result().Cookies()[0]
	if cookie.Name != "cookieName" {
		t.Errorf("Expected cookie name %s, but got %s", "cookieName", cookie.Name)
	}
	if cookie.Value != "cookieValue" {
		t.Errorf("Expected cookie value %s, but got %s", "cookieValue", cookie.Value)
	}
}

func TestSet_Signed(t *testing.T) {
	DefaultManager = signedManager

	w := httptest.NewRecorder()

	cookieName := "myCookie"
	expectedValue := "myValue"

	err := Set(w, cookieName, expectedValue, Options{Signed: true})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()

	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, but got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != cookieName {
		t.Errorf("Expected cookie name '%s', but got '%s'", cookieName, cookie.Name)
	}

	expectedCookieValue := signCookieValue(expectedValue, signedManager.signingKey)
	if cookie.Value != expectedCookieValue {
		t.Errorf("Expected cookie value '%s', but got '%s'", expectedCookieValue, cookie.Value)
	}
}

func TestSetSigned(t *testing.T) {
	DefaultManager = signedManager

	w := httptest.NewRecorder()

	cookieName := "myCookie"
	expectedValue := "myValue"

	err := SetSigned(w, cookieName, expectedValue, Options{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cookies := w.Result().Cookies()

	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, but got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != cookieName {
		t.Errorf("Expected cookie name '%s', but got '%s'", cookieName, cookie.Name)
	}

	expectedCookieValue := signCookieValue(expectedValue, signedManager.signingKey)
	if cookie.Value != expectedCookieValue {
		t.Errorf("Expected cookie value '%s', but got '%s'", expectedCookieValue, cookie.Value)
	}
}

func TestRemove(t *testing.T) {
	// Create a mock request
	_, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Call the Remove function
	err = Remove(rr, "cookieName")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the cookie was set in the response
	cookie := rr.Result().Cookies()[0]
	if cookie.Name != "cookieName" {
		t.Errorf("Expected cookie name %s, but got %s", "cookieName", cookie.Name)
	}
	if cookie.Value != "" {
		t.Errorf("Expected cookie value %s, but got %s", "", cookie.Value)
	}
}

func TestPopulateFromCookies(t *testing.T) {
	DefaultManager = signedManager
	value := "test"

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "cookie1", Value: value})

	type MyStruct struct {
		Default string `cookie:"cookie1"`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(req, dest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := &MyStruct{
		Default: value,
	}
	if dest.Default != expected.Default {
		t.Errorf("Expected value '%s', but got '%s'", expected.Default, dest.Default)
	}
}
