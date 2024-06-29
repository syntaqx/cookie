package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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

var (
	secretKey = "very-secret-key"
)

// generateHMAC generates an HMAC signature for a given value using the secret key.
func generateHMAC(value string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

// verifyHMAC verifies the HMAC signature of a given value.
func verifyHMAC(value, signature string) bool {
	expectedSignature := generateHMAC(value)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// Set sets a signed cookie with the given name, value, and options.
func Set(w http.ResponseWriter, name string, value interface{}, options *http.Cookie) error {
	if options == nil {
		options = &http.Cookie{}
	}

	valueStr := fmt.Sprintf("%v", value)
	signature := generateHMAC(valueStr)
	signedValue := fmt.Sprintf("%s|%s", valueStr, signature)

	cookie := &http.Cookie{
		Name:     name,
		Value:    signedValue,
		Path:     options.Path,
		Domain:   options.Domain,
		Expires:  options.Expires,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	}
	http.SetCookie(w, cookie)
	return nil
}

// Get retrieves and verifies the value of a signed cookie with the given name.
// It gracefully allows unsigned cookie values.
func Get(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(cookie.Value, "|", 2)
	if len(parts) == 1 {
		// Unsigned cookie value
		return parts[0], nil
	}

	value, signature := parts[0], parts[1]
	if !verifyHMAC(value, signature) {
		return "", errors.New("invalid cookie signature")
	}

	return value, nil
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

// PopulateFromCookies populates the fields of a struct based on signed cookie tags.
func PopulateFromCookies(r *http.Request, dest interface{}) error {
	val := reflect.ValueOf(dest).Elem()
	typ := val.Type()
	var errs []error

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get(CookieTag)

		if tag == "" {
			continue
		}

		cookieValue, err := Get(r, tag)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return ErrNoCookie
			}
			return err
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(cookieValue)
		case reflect.Int:
			intVal, err := strconv.Atoi(cookieValue)
			if err != nil {
				return err
			}
			field.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(cookieValue)
			if err != nil {
				return err
			}
			field.SetBool(boolVal)
		case reflect.Slice:
			switch fieldType.Type.Elem().Kind() {
			case reflect.String:
				field.Set(reflect.ValueOf(strings.Split(cookieValue, ",")))
			case reflect.Int:
				intStrings := strings.Split(cookieValue, ",")
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
				uid, err := uuid.FromString(cookieValue)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(uid))
			}
		case reflect.Struct:
			if fieldType.Type == reflect.TypeOf(time.Time{}) {
				timeVal, err := time.Parse(time.RFC3339, cookieValue)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(timeVal))
			}
		default:
			return &UnsupportedTypeError{fieldType.Type}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors occurred during cookie population: %v", errs)
	}
	return nil
}
