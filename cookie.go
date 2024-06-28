package cookie

import (
	"net/http"
	"reflect"

	"github.com/gofrs/uuid/v5"
)

// PopulateFromCookies populates the fields of a struct based on cookie tags.
func PopulateFromCookies(r *http.Request, dest interface{}) error {
	val := reflect.ValueOf(dest).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("cookie")

		if tag == "" {
			continue
		}

		cookie, err := r.Cookie(tag)
		if err != nil {
			return err
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(cookie.Value)
		case reflect.Array:
			if fieldType.Type == reflect.TypeOf(uuid.UUID{}) {
				uid, err := uuid.FromString(cookie.Value)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(uid))
			} else {
				return &UnsupportedTypeError{fieldType.Type}
			}
		default:
			return &UnsupportedTypeError{fieldType.Type}
		}
	}
	return nil
}

// UnsupportedTypeError is returned when a field type is not supported by PopulateFromCookies.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "cookie: unsupported type: " + e.Type.String()
}

// Set sets a cookie with the given name, value, and options.
func Set(w http.ResponseWriter, name, value string, options *http.Cookie) {
	if options == nil {
		options = &http.Cookie{}
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		Expires:  options.Expires,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	}
	http.SetCookie(w, cookie)
}

// Get retrieves the value of a cookie with the given name.
func Get(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Remove removes a cookie by setting its MaxAge to -1.
func Remove(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
