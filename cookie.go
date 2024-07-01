package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	CookieTag         = "cookie"
	DefaultSigningKey = "default-signing-key"
)

// Options contains the options for a cookie.
type Options struct {
	Path     string
	Domain   string
	Expires  time.Time
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
	Signed   bool
}

var (
	// SigningKey is the key used to sign cookies.
	SigningKey = []byte(DefaultSigningKey)

	// DefaultOptions are the default options for cookies.
	DefaultOptions = &Options{
		Path:     "/",
		Domain:   "",
		Expires:  time.Time{},
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteDefaultMode,
		Signed:   false,
	}
)

var (
	// ErrInvalidSignedFormat is returned when a signed cookie is not in the correct format.
	ErrInvalidSignedFormat = errors.New("cookie: invalid signed cookie format")

	// ErrInvalidSignature is returned when a signed cookie has an invalid signature.
	ErrInvalidSignature = errors.New("cookie: invalid cookie signature")
)

// ErrUnsupportedType is returned when a field type is not supported.
type ErrUnsupportedType struct {
	Type reflect.Type
}

// Error returns the error message.
func (e *ErrUnsupportedType) Error() string {
	return "cookie: unsupported type: " + e.Type.String()
}

// Set sets a cookie with the given name, value, and options.
func Set(w http.ResponseWriter, name, value string, options *Options) {
	mergedOptions := mergeOptions(options, DefaultOptions)
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     mergedOptions.Path,
		Domain:   mergedOptions.Domain,
		Expires:  mergedOptions.Expires,
		MaxAge:   mergedOptions.MaxAge,
		Secure:   mergedOptions.Secure,
		HttpOnly: mergedOptions.HttpOnly,
		SameSite: mergedOptions.SameSite,
	}

	if mergedOptions.Signed {
		signature := generateHMAC(value)
		cookie.Value = base64.URLEncoding.EncodeToString([]byte(value)) + "|" + signature
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

// SetSigned sets a signed cookie with the given name, value, and options.
func SetSigned(w http.ResponseWriter, name, value string, options *Options) {
	if options == nil {
		options = &Options{}
	}

	options.Signed = true
	Set(w, name, value, options)
}

// GetSigned retrieves the value of a signed cookie with the given name.
func GetSigned(r *http.Request, name string) (string, error) {
	signedValue, err := Get(r, name)
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(signedValue, "|", 2)
	if len(parts) != 2 {
		return "", ErrInvalidSignedFormat
	}

	value, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	signature, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, SigningKey)
	h.Write(value)
	expectedSignature := h.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return "", ErrInvalidSignature
	}

	return string(value), nil
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
		tagParts := strings.Split(tag, ",")

		if tagParts[0] == "" {
			continue
		}

		var cookie string
		var err error
		isSigned := DefaultOptions.Signed

		for _, part := range tagParts[1:] {
			if part == "signed" {
				isSigned = true
			} else if part == "unsigned" {
				isSigned = false
			}
		}

		if isSigned {
			cookie, err = GetSigned(r, tagParts[0])
		} else {
			cookie, err = Get(r, tagParts[0])
		}

		if err != nil {
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
				return &ErrUnsupportedType{fieldType.Type}
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
			return &ErrUnsupportedType{fieldType.Type}
		}
	}
	return nil
}

func generateHMAC(value string) string {
	h := hmac.New(sha256.New, SigningKey)
	h.Write([]byte(value))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func mergeOptions(provided, defaults *Options) *Options {
	if provided == nil {
		return defaults
	}

	merged := *defaults

	if provided.Path != "" {
		merged.Path = provided.Path
	}
	if provided.Domain != "" {
		merged.Domain = provided.Domain
	}
	if !provided.Expires.IsZero() {
		merged.Expires = provided.Expires
	}
	if provided.MaxAge != 0 {
		merged.MaxAge = provided.MaxAge
	}
	if provided.Secure {
		merged.Secure = provided.Secure
	}
	if provided.HttpOnly {
		merged.HttpOnly = provided.HttpOnly
	}
	if provided.SameSite != http.SameSiteDefaultMode {
		merged.SameSite = provided.SameSite
	}
	if provided.Signed {
		merged.Signed = provided.Signed
	}

	return &merged
}
