package cookie

import "net/http"

// DefaultManager is the default cookie manager exposed by this package.
var DefaultManager = NewManager()

// Get retrieves an unsigned cooke value.
func Get(r *http.Request, name string) (string, error) {
	return DefaultManager.Get(r, name)
}

// GetSigned retrieves a signed cookie value.
func GetSigned(r *http.Request, name string) (string, error) {
	return DefaultManager.GetSigned(r, name)
}

// Set sets the value of a cookie.
func Set(w http.ResponseWriter, name, value string, opts ...Options) error {
	return DefaultManager.Set(w, name, value, opts...)
}

// SetSigned sets a signed value of a cookie.
func SetSigned(w http.ResponseWriter, name, value string, opts ...Options) error {
	return DefaultManager.SetSigned(w, name, value, opts...)
}

// Remove removes a cookie from the response.
func Remove(w http.ResponseWriter, name string) error {
	return DefaultManager.Remove(w, name)
}

// PopulateFromCookies populates a struct with cookie values.
func PopulateFromCookies(r *http.Request, dest interface{}) error {
	return DefaultManager.PopulateFromCookies(r, dest)
}
