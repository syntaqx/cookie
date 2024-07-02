package cookie

import (
	"errors"
	"reflect"
)

// ErrInvalidSignedCookieFormat is returned when the format of a signed cookie is invalid.
var ErrInvalidSignedCookieFormat = errors.New("invalid signed cookie format")

// ErrInvalidCookieSignature is returned when the signature of a signed cookie is invalid.
var ErrInvalidCookieSignature = errors.New("invalid cookie signature")

// ErrNonNilPointerRequired is returned when the destination parameter must be a non-nil pointer.
var ErrNonNilPointerRequired = errors.New("dest must be a non-nil pointer")

// ErrUnsupportedType is returned when a field type is not supported.
type ErrUnsupportedType struct {
	Type reflect.Type
}

// Error returns the error message.
func (e *ErrUnsupportedType) Error() string {
	return "cookie: unsupported type: " + e.Type.String()
}
