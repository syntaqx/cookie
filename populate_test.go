package cookie

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestPopulateFromCookies(t *testing.T) {
	// Create a mock HTTP request with cookies
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create an instance of the Manager
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

	// Create a struct to populate with cookie values
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

	// Call the PopulateFromCookies function
	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the populated values
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
	// Create a mock HTTP request with cookies
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create an instance of the Manager
	manager := NewManager()

	// Call the PopulateFromCookies function with a nil pointer
	var dest *struct{}
	err := manager.PopulateFromCookies(req, dest)
	if err != ErrNonNilPointerRequired {
		t.Errorf("Unexpected error: %v", err)
	}
}
