package cookie

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManager_Get(t *testing.T) {
	m := NewManager()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	cookieValue := "myValue"

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	value, err := m.Get(r, cookieName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != cookieValue {
		t.Errorf("Expected value '%s', but got '%s'", cookieValue, value)
	}
}

func TestManager_GetError(t *testing.T) {
	m := NewManager()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"

	_, err := m.Get(r, cookieName)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func TestManager_GetSigned(t *testing.T) {
	m := NewManager(
		WithSigningKey([]byte("super-secret-key")),
	)
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	expectedValue := "myValue"

	data := base64.URLEncoding.EncodeToString([]byte(expectedValue))
	signature := base64.URLEncoding.EncodeToString(sign([]byte(data), m.signingKey))
	cookieValue := data + "|" + signature

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	value, err := m.GetSigned(r, cookieName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != expectedValue {
		t.Errorf("Expected value '%s', but got '%s'", expectedValue, value)
	}
}

func TestManager_GetSignedError(t *testing.T) {
	m := NewManager(
		WithSigningKey([]byte("super-secret-key")),
	)
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"

	_, err := m.GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func TestManager_Set(t *testing.T) {
	m := NewManager()
	w := httptest.NewRecorder()

	cookieName := "myCookie"
	cookieValue := "myValue"

	err := m.Set(w, cookieName, cookieValue)
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

	if cookie.Value != cookieValue {
		t.Errorf("Expected cookie value '%s', but got '%s'", cookieValue, cookie.Value)
	}
}

func TestManager_SetSigned(t *testing.T) {
	m := NewManager(
		WithSigningKey([]byte("super-secret-key")),
	)
	w := httptest.NewRecorder()

	cookieName := "myCookie"
	expectedValue := "myValue"

	err := m.Set(w, cookieName, expectedValue, Options{Signed: true})
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

	data := base64.URLEncoding.EncodeToString([]byte(expectedValue))
	signature := base64.URLEncoding.EncodeToString(sign([]byte(data), m.signingKey))
	expectedCookieValue := data + "|" + signature

	if cookie.Value != expectedCookieValue {
		t.Errorf("Expected cookie value '%s', but got '%s'", expectedCookieValue, cookie.Value)
	}
}

func TestManager_Remove(t *testing.T) {
	m := NewManager()
	w := httptest.NewRecorder()

	cookieName := "myCookie"

	err := m.Remove(w, cookieName)
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

	if cookie.Value != "" {
		t.Errorf("Expected empty cookie value, but got '%s'", cookie.Value)
	}

	if cookie.Expires.Unix() != 0 {
		t.Errorf("Expected cookie to be expired, but it expires at %v", cookie.Expires)
	}

	if cookie.MaxAge != -1 {
		t.Errorf("Expected cookie to be expired, but it has MaxAge %d", cookie.MaxAge)
	}
}

func TestManager_RemoveWithOptions(t *testing.T) {
	m := NewManager()
	w := httptest.NewRecorder()

	cookieName := "myCookie"
	path := "/path"

	err := m.Remove(w, cookieName, Options{
		Path: path,
	})

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

	if cookie.Path != path {
		t.Errorf("Expected cookie path '%s', but got '%s'", path, cookie.Path)
	}
}
