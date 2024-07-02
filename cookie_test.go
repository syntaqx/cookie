package cookie

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

var unsignedManager = NewManager()
var signedManager = NewManager(WithSigningKey([]byte("super-secret-key")))

func TestManager_Get(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	cookieValue := "myValue"

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	value, err := unsignedManager.Get(r, cookieName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != cookieValue {
		t.Errorf("Expected value '%s', but got '%s'", cookieValue, value)
	}
}

func TestManager_GetError(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"

	_, err := unsignedManager.Get(r, cookieName)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func TestManager_GetSigned(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	expectedValue := "myValue"

	cookieValue := signCookieValue(expectedValue, signedManager.signingKey)

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	value, err := signedManager.GetSigned(r, cookieName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != expectedValue {
		t.Errorf("Expected value '%s', but got '%s'", expectedValue, value)
	}
}

func TestManager_GetSignedError(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"

	_, err := signedManager.GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func TestManager_GetSignedInvalidFormat(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	cookieValue := "invalidFormat"

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	_, err := signedManager.GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	if err != ErrInvalidSignedCookieFormat {
		t.Errorf("Expected error '%v', but got '%v'", ErrInvalidSignedCookieFormat, err)
	}
}

func TestManager_GetSignedInvalidSignature(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	expectedValue := "myValue"

	data := base64.URLEncoding.EncodeToString([]byte(expectedValue))
	signature := base64.URLEncoding.EncodeToString(sign([]byte("invalidData"), signedManager.signingKey))
	cookieValue := data + "|" + signature

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	_, err := signedManager.GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	if err != ErrInvalidCookieSignature {
		t.Errorf("Expected error '%v', but got '%v'", ErrInvalidCookieSignature, err)
	}
}

func TestManager_GetSigned_Base64DataDecodeError(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	cookieValue := "invalidBase64|invalidBase64"

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	_, err := signedManager.GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "illegal base64 data at input byte 12"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}

func TestManager_GetSigned_Base64SignatureDecodeError(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	cookieName := "myCookie"
	cookieValue := "ZXhhbXBsZQ==|invalidBase64"

	r.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})

	_, err := signedManager.GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "illegal base64 data at input byte 12"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}

func TestManager_Set(t *testing.T) {
	w := httptest.NewRecorder()

	cookieName := "myCookie"
	cookieValue := "myValue"

	err := unsignedManager.Set(w, cookieName, cookieValue)
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

func TestManager_Set_Signed(t *testing.T) {
	w := httptest.NewRecorder()

	cookieName := "myCookie"
	expectedValue := "myValue"

	err := signedManager.Set(w, cookieName, expectedValue, Options{Signed: true})
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

func TestManager_SetSigned(t *testing.T) {
	w := httptest.NewRecorder()

	cookieName := "myCookie"
	expectedValue := "myValue"

	err := signedManager.SetSigned(w, cookieName, expectedValue, Options{})
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

func TestManager_Remove(t *testing.T) {
	w := httptest.NewRecorder()

	cookieName := "myCookie"

	err := unsignedManager.Remove(w, cookieName)
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
	w := httptest.NewRecorder()

	cookieName := "myCookie"
	path := "/path"

	err := unsignedManager.Remove(w, cookieName, Options{
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
