package cookie

import "errors"

// ErrInvalidSignedCookieFormat is returned when the format of a signed cookie is invalid.
var ErrInvalidSignedCookieFormat = errors.New("invalid signed cookie format")

// ErrInvalidCookieSignature is returned when the signature of a signed cookie is invalid.
var ErrInvalidCookieSignature = errors.New("invalid cookie signature")

// ErrUnsupportedFieldType is returned when a field type is not supported.
var ErrUnsupportedFieldType = errors.New("unsupported field type")

// ErrNonNilPointerRequired is returned when the destination parameter must be a non-nil pointer.
var ErrNonNilPointerRequired = errors.New("dest must be a non-nil pointer")
