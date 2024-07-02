package cookie

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestManager_PopulateFromCookies(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	manager := NewManager(WithSigningKey([]byte("super-secret-key")))

	value := "test"
	data := base64.URLEncoding.EncodeToString([]byte(value))
	signature := base64.URLEncoding.EncodeToString(sign([]byte(data), manager.signingKey))
	signedValue := data + "|" + signature

	cookies := []*http.Cookie{
		{Name: "cookie1", Value: value},
		{Name: "cookie2", Value: value},
		{Name: "cookie3", Value: string(signedValue)},
		{Name: "cookie4", Value: "true"},
		{Name: "cookie5", Value: "123"},
		{Name: "cookie6", Value: "123.45"},
		{Name: "cookie7", Value: "a,b,c"},
		{Name: "cookie8", Value: "1,2,3"},
		{Name: "cookie9", Value: "2021-01-02T15:04:05Z"},
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	type MyStruct struct {
		UntaggedField string
		Default       string    `cookie:"cookie1"`
		Unsigned      string    `cookie:"cookie2,unsigned"`
		Signed        string    `cookie:"cookie3,signed"`
		Boolean       bool      `cookie:"cookie4"`
		Integer       int       `cookie:"cookie5"`
		UInteger      uint      `cookie:"cookie5"`
		Float         float64   `cookie:"cookie6"`
		StringSlice   []string  `cookie:"cookie7"`
		IntSlice      []int     `cookie:"cookie8"`
		Timestamp     time.Time `cookie:"cookie9"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := &MyStruct{
		Default:     value,
		Unsigned:    value,
		Signed:      value,
		Boolean:     true,
		Integer:     123,
		UInteger:    123,
		Float:       123.45,
		StringSlice: []string{"a", "b", "c"},
		IntSlice:    []int{1, 2, 3},
		Timestamp:   time.Date(2021, 1, 2, 15, 4, 5, 0, time.UTC),
	}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Unexpected result. Got: %v, want: %v", dest, expected)
	}
}

func TestPopulateFromCookies_NonNilPointerRequired(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	manager := NewManager()

	var dest *struct{}
	err := manager.PopulateFromCookies(req, dest)
	if err != ErrNonNilPointerRequired {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestPopulateFromCookies_ErrNoCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	manager := NewManager()

	type MyStruct struct {
		Field string `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err != http.ErrNoCookie {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestPopulateFromCookies_ErrUnsupportedType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: "cookie", Value: "test"})

	manager := NewManager()

	type MyStruct struct {
		Field complex128 `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "cookie: unsupported type: complex128"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}

func TestPopulateFromCookies_InvalidBoolean(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: "cookie", Value: "invalid"})

	manager := NewManager()

	type MyStruct struct {
		Field bool `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "strconv.ParseBool: parsing \"invalid\": invalid syntax"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}

func TestPopulateFromCookies_InvalidInteger(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: "cookie", Value: "invalid"})

	manager := NewManager()

	type MyStruct struct {
		Field int `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "strconv.ParseInt: parsing \"invalid\": invalid syntax"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}

func TestPopulateFromCookies_InvalidUnsignedInteger(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: "cookie", Value: "-1"})

	manager := NewManager()

	type MyStruct struct {
		Field uint `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "strconv.ParseUint: parsing \"-1\": invalid syntax"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}

func TestPopulateFromCookies_InvalidFloat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: "cookie", Value: "invalid"})

	manager := NewManager()

	type MyStruct struct {
		Field float64 `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "strconv.ParseFloat: parsing \"invalid\": invalid syntax"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}

func TestPopulateFromCookies_InvalidIntSlice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: "cookie", Value: "invalid"})

	manager := NewManager()

	type MyStruct struct {
		Field []int `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "strconv.Atoi: parsing \"invalid\": invalid syntax"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}

func TestPopulateFromCookies_InvalidTimestamp(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{Name: "cookie", Value: "invalid"})

	manager := NewManager()

	type MyStruct struct {
		Field time.Time `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	expectedError := "parsing time \"invalid\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid\" as \"2006\""
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%v'", expectedError, err)
	}
}
