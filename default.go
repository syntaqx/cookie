package cookie

import "net/http"

var DefaultManager = NewManager()

func Get(r *http.Request, name string) (string, error) {
	return DefaultManager.Get(r, name)
}

func GetSigned(r *http.Request, name string) (string, error) {
	return DefaultManager.GetSigned(r, name)
}

func Set(w http.ResponseWriter, name, value string, opts ...Options) error {
	return DefaultManager.Set(w, name, value, opts...)
}

func SetSigned(w http.ResponseWriter, name, value string, opts ...Options) error {
	return DefaultManager.SetSigned(w, name, value, opts...)
}

func Remove(w http.ResponseWriter, name string) error {
	return DefaultManager.Remove(w, name)
}

func PopulateFromCookies(r *http.Request, dest interface{}) error {
	return DefaultManager.PopulateFromCookies(r, dest)
}
