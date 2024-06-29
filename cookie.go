package cookie

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	CookieTag = "cookie"
)

// ErrNoCookie is returned when a cookie is not found.
var ErrNoCookie = http.ErrNoCookie

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

// PopulateFromCookies populates the fields of a struct based on cookie tags.
func PopulateFromCookies(r *http.Request, dest interface{}) error {
	val := reflect.ValueOf(dest).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get(CookieTag)

		if tag == "" {
			continue
		}

		cookie, err := Get(r, tag)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return ErrNoCookie
			}
			return err
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(cookie)
		case reflect.Int:
			intVal, err := strconv.Atoi(cookie)
			if err != nil {
				return err
			}
			field.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(cookie)
			if err != nil {
				return err
			}
			field.SetBool(boolVal)
		case reflect.Slice:
			switch fieldType.Type.Elem().Kind() {
			case reflect.String:
				field.Set(reflect.ValueOf(strings.Split(cookie, ",")))
			case reflect.Int:
				intStrings := strings.Split(cookie, ",")
				intSlice := make([]int, len(intStrings))
				for i, s := range intStrings {
					intVal, err := strconv.Atoi(s)
					if err != nil {
						return err
					}
					intSlice[i] = intVal
				}
				field.Set(reflect.ValueOf(intSlice))
			default:
				return &UnsupportedTypeError{fieldType.Type}
			}
		case reflect.Array:
			if fieldType.Type == reflect.TypeOf(uuid.UUID{}) {
				uid, err := uuid.FromString(cookie)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(uid))
			}
		case reflect.Struct:
			if fieldType.Type == reflect.TypeOf(time.Time{}) {
				timeVal, err := time.Parse(time.RFC3339, cookie)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(timeVal))
			}
		default:
			return &UnsupportedTypeError{fieldType.Type}
		}
	}
	return nil
}
