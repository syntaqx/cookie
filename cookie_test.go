package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	_, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	name := "myCookie"
	value := "myValue"

	options := &http.Cookie{
		Path:     "/",
		Domain:   "example.com",
		Expires:  time.Now().Add(24 * time.Hour).UTC(),
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	Set(w, name, value, options)

	// Get the response cookies
	cookies := w.Result().Cookies()

	// Check if the cookie was set correctly
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != name {
		t.Errorf("Expected cookie name %s, got %s", name, cookie.Name)
	}
	if cookie.Value != value {
		t.Errorf("Expected cookie value %s, got %s", value, cookie.Value)
	}
	if cookie.Path != options.Path {
		t.Errorf("Expected cookie path %s, got %s", options.Path, cookie.Path)
	}
	if cookie.Domain != options.Domain {
		t.Errorf("Expected cookie domain %s, got %s", options.Domain, cookie.Domain)
	}
	if cookie.MaxAge != options.MaxAge {
		t.Errorf("Expected cookie max age %d, got %d", options.MaxAge, cookie.MaxAge)
	}
	if cookie.Secure != options.Secure {
		t.Errorf("Expected cookie secure %t, got %t", options.Secure, cookie.Secure)
	}
	if cookie.HttpOnly != options.HttpOnly {
		t.Errorf("Expected cookie HttpOnly %t, got %t", options.HttpOnly, cookie.HttpOnly)
	}
	if cookie.SameSite != options.SameSite {
		t.Errorf("Expected cookie SameSite %d, got %d", options.SameSite, cookie.SameSite)
	}
}
