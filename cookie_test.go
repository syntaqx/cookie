package cookie

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
)

func TestSet(t *testing.T) {
	_, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	name := "myCookie"
	value := "myValue"

	options := &Options{
		Path:     "/",
		Domain:   "example.com",
		Expires:  time.Now().Add(24 * time.Hour),
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

func TestSet_WithoutOptions(t *testing.T) {
	_, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	name := "myCookie"
	value := "myValue"

	Set(w, name, value, nil)

	// Get the response cookies
	cookies := w.Result().Cookies()

	// Check if the cookie was set correctly
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
}

func TestSetSigned(t *testing.T) {
	_, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	name := "myCookie"
	value := "myValue"

	SetSigned(w, name, value, nil)

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

	// Check if the cookie value is signed
	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		t.Errorf("Expected signed cookie value, got %s", cookie.Value)
	}
}

func TestGet(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookieName := "myCookie"
	cookieValue := "myValue"
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: cookieValue,
	}
	r.AddCookie(cookie)

	value, err := Get(r, cookieName)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value != cookieValue {
		t.Errorf("Expected cookie value %s, got %s", cookieValue, value)
	}
}

func TestGetNonexistentCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookieName := "nonexistentCookie"

	_, err := Get(r, cookieName)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetSigned(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookieName := "myCookie"
	cookieValue := "myValue"
	signature := generateHMAC(cookieValue)
	signedValue := base64.URLEncoding.EncodeToString([]byte(cookieValue)) + "|" + signature

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: signedValue,
	}

	r.AddCookie(cookie)

	value, err := GetSigned(r, cookieName)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value != cookieValue {
		t.Errorf("Expected cookie value %s, got %s", cookieValue, value)
	}
}

func TestGetSignedNonexistentCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookieName := "nonexistentCookie"

	_, err := GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetSignedInvalidValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookieName := "myCookie"
	cookieValue := "myValue"
	signedValue := cookieValue + "|invalid"

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: signedValue,
	}

	r.AddCookie(cookie)

	_, err := GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetSignedInvalidSignature(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookieName := "myCookie"
	cookieValue := "myValue"
	signedValue := base64.URLEncoding.EncodeToString([]byte(cookieValue)) + "|invalid"

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: signedValue,
	}

	r.AddCookie(cookie)

	_, err := GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetSignedInvalidHMAC(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookieName := "myCookie"
	cookieValue := "myValue"
	signature := generateHMAC("invalid")
	signedValue := base64.URLEncoding.EncodeToString([]byte(cookieValue)) + "|" + signature

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: signedValue,
	}

	r.AddCookie(cookie)

	_, err := GetSigned(r, cookieName)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRemove(t *testing.T) {
	w := httptest.NewRecorder()
	name := "myCookie"
	Remove(w, name)
	// Get the response cookies
	cookies := w.Result().Cookies()
	// Check if the cookie was removed correctly
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != name {
		t.Errorf("Expected cookie name %s, got %s", name, cookie.Name)
	}
	if cookie.Value != "" {
		t.Errorf("Expected cookie value to be empty, got %s", cookie.Value)
	}
	if cookie.Path != "/" {
		t.Errorf("Expected cookie path /, got %s", cookie.Path)
	}
	if cookie.MaxAge != -1 {
		t.Errorf("Expected cookie max age -1, got %d", cookie.MaxAge)
	}
}

func TestPopulateFromCookies(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	cookies := map[string]string{
		"myCookie":         "myValue",
		"myIntCookie":      "123",
		"myBoolCookie":     "true",
		"mySliceCookie":    "val1,val2,val3",
		"myIntSliceCookie": "1,2,3",
		"myUUIDCookie":     uuid.Must(uuid.NewV4()).String(),
		"myTimeCookie":     time.Now().Format(time.RFC3339),
	}
	for name, value := range cookies {
		r.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	r.AddCookie(&http.Cookie{
		Name:  "signedCookie",
		Value: base64.URLEncoding.EncodeToString([]byte("signedValue")) + "|" + generateHMAC("signedValue"),
	})

	type MyStruct struct {
		StringField  string    `cookie:"myCookie"`
		IntField     int       `cookie:"myIntCookie"`
		BoolField    bool      `cookie:"myBoolCookie"`
		StringSlice  []string  `cookie:"mySliceCookie"`
		IntSlice     []int     `cookie:"myIntSliceCookie"`
		UUIDField    uuid.UUID `cookie:"myUUIDCookie"`
		TimeField    time.Time `cookie:"myTimeCookie"`
		SignedCookie string    `cookie:"signedCookie,signed"`
		Unsupported  complex64 `cookie:""`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(r, dest)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if dest.StringField != cookies["myCookie"] {
		t.Errorf("Expected StringField %s, got %s", cookies["myCookie"], dest.StringField)
	}

	expectedInt, _ := strconv.Atoi(cookies["myIntCookie"])
	if dest.IntField != expectedInt {
		t.Errorf("Expected IntField %d, got %d", expectedInt, dest.IntField)
	}

	expectedBool, _ := strconv.ParseBool(cookies["myBoolCookie"])
	if dest.BoolField != expectedBool {
		t.Errorf("Expected BoolField %t, got %t", expectedBool, dest.BoolField)
	}

	expectedStringSlice := strings.Split(cookies["mySliceCookie"], ",")
	if !reflect.DeepEqual(dest.StringSlice, expectedStringSlice) {
		t.Errorf("Expected StringSlice %v, got %v", expectedStringSlice, dest.StringSlice)
	}

	intStrings := strings.Split(cookies["myIntSliceCookie"], ",")
	expectedIntSlice := make([]int, len(intStrings))
	for i, s := range intStrings {
		expectedIntSlice[i], _ = strconv.Atoi(s)
	}
	if !reflect.DeepEqual(dest.IntSlice, expectedIntSlice) {
		t.Errorf("Expected IntSlice %v, got %v", expectedIntSlice, dest.IntSlice)
	}

	expectedUUID, _ := uuid.FromString(cookies["myUUIDCookie"])
	if dest.UUIDField != expectedUUID {
		t.Errorf("Expected UUIDField %s, got %s", expectedUUID, dest.UUIDField)
	}

	expectedTime, _ := time.Parse(time.RFC3339, cookies["myTimeCookie"])
	if !dest.TimeField.Equal(expectedTime) {
		t.Errorf("Expected TimeField %v, got %v", expectedTime, dest.TimeField)
	}

	expectedSignedValue := "signedValue"
	if dest.SignedCookie != expectedSignedValue {
		t.Errorf("Expected SignedCookie %s, got %s", expectedSignedValue, dest.SignedCookie)
	}
}

func TestPopulateFromCookies_InvalidIntValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "myIntCookie",
		Value: "invalid",
	})

	type MyStruct struct {
		IntField int `cookie:"myIntCookie"`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(r, dest)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestPopulateFromCookies_InvalidBoolValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "myBoolCookie",
		Value: "invalid",
	})

	type MyStruct struct {
		BoolField bool `cookie:"myBoolCookie"`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(r, dest)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestPopulateFromCookies_InvalidIntSliceValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "myIntSliceCookie",
		Value: "1,2,invalid",
	})

	type MyStruct struct {
		IntSlice []int `cookie:"myIntSliceCookie"`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(r, dest)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestPopulateFromCookies_InvalidUUIDValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "myUUIDCookie",
		Value: "invalid",
	})

	type MyStruct struct {
		UUIDField uuid.UUID `cookie:"myUUIDCookie"`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(r, dest)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestPopulateFromCookies_InvalidTimeValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "myTimeCookie",
		Value: "invalid",
	})

	type MyStruct struct {
		TimeField time.Time `cookie:"myTimeCookie"`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(r, dest)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestPopulateFromCookies_NotFound(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	type MyStruct struct {
		StringField string `cookie:"myCookie"`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(r, dest)
	if err != ErrNoCookie {
		t.Errorf("Expected error ErrNoCookie, got %v", err)
	}
}

func TestPopulateFromCookies_UnsupportedType(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "myCookie",
		Value: "myValue",
	})

	type MyStruct struct {
		Unsupported complex64 `cookie:"myCookie"`
	}

	dest := &MyStruct{}
	err := PopulateFromCookies(r, dest)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*UnsupportedTypeError); !ok {
		t.Errorf("Expected error of type UnsupportedTypeError, got %T", err)
	}

	expected := "cookie: unsupported type: complex64"
	if err.Error() != expected {
		t.Errorf("Expected error message %s, got %s", expected, err.Error())
	}
}
