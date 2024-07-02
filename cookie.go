package cookie

import (
	"encoding/base64"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// Options represent the options for an HTTP cookie as sent in the Set-Cookie
// header of an HTTP response or the Cookie header of an HTTP request.
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

// Manager handles cookie operations.
type Manager struct {
	signingKey     []byte
	customHandlers map[reflect.Type]CustomTypeHandler
}

// Option is a function type for configuring the Manager.
type Option func(*Manager)

// WithSigningKey sets the signing key for the Manager.
func WithSigningKey(key []byte) Option {
	return func(m *Manager) {
		m.signingKey = key
	}
}

// WithCustomHandler registers a custom type handler for the Manager.
func WithCustomHandler(typ reflect.Type, handler CustomTypeHandler) Option {
	return func(m *Manager) {
		m.customHandlers[typ] = handler
	}
}

// NewManager creates a new Manager with the given options.
func NewManager(opts ...Option) *Manager {
	m := &Manager{
		customHandlers: make(map[reflect.Type]CustomTypeHandler),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Get retrieves the value of an unsigned cookie.
func (m *Manager) Get(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetSigned retrieves the value of a signed cookie.
func (m *Manager) GetSigned(r *http.Request, name string) (string, error) {
	value, err := m.Get(r, name)
	if err != nil {
		return "", err
	}

	parts := strings.Split(value, "|")
	if len(parts) != 2 {
		return "", ErrInvalidSignedCookieFormat
	}

	data, signature := parts[0], parts[1]
	dataBytes, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	signatureBytes, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return "", err
	}

	if verify([]byte(data), signatureBytes, m.signingKey) {
		return string(dataBytes), nil
	}
	return "", ErrInvalidCookieSignature
}

// Set sets the value of a cookie.
func (m *Manager) Set(w http.ResponseWriter, name, value string, opts ...Options) error {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}

	if o.Signed && m.signingKey != nil {
		data := base64.URLEncoding.EncodeToString([]byte(value))
		signature := base64.URLEncoding.EncodeToString(sign([]byte(data), m.signingKey))
		value = data + "|" + signature
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     o.Path,
		Domain:   o.Domain,
		Expires:  o.Expires,
		MaxAge:   o.MaxAge,
		Secure:   o.Secure,
		HttpOnly: o.HttpOnly,
		SameSite: o.SameSite,
	}

	http.SetCookie(w, cookie)
	return nil
}

// SetSigned sets the value of a signed cookie.
func (m *Manager) SetSigned(w http.ResponseWriter, name, value string, opts ...Options) error {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}
	o.Signed = true
	return m.Set(w, name, value, o)
}

// Remove deletes a cookie.
func (m *Manager) Remove(w http.ResponseWriter, name string, opts ...Options) error {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     o.Path,
		Domain:   o.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   o.Secure,
		HttpOnly: o.HttpOnly,
		SameSite: o.SameSite,
	}
	http.SetCookie(w, cookie)
	return nil
}
