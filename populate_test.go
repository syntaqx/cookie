package cookie

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Create a struct to populate with cookie values
	type MyStruct struct {
		UntaggedField string
		Default       string `cookie:"cookie1"`
		Unsigned      string `cookie:"cookie2,unsigned"`
		Signed        string `cookie:"cookie3,signed"`
	}

	// Call the PopulateFromCookies function
	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the populated values
	expected := &MyStruct{
		Default:  value,
		Unsigned: value,
		Signed:   value,
	}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Unexpected result. Got: %v, want: %v", dest, expected)
	}
}
