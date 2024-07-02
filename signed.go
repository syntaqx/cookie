package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// sign generates a HMAC signature for the given data using the provided key.
func sign(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// verify checks the HMAC signature of the given data using the provided key.
func verify(data, signature, key []byte) bool {
	expectedSignature := sign(data, key)
	return hmac.Equal(expectedSignature, signature)
}

// signCookieValue signs a cookie value using the provided key.
func signCookieValue(value string, key []byte) string {
	data := base64.URLEncoding.EncodeToString([]byte(value))
	signature := base64.URLEncoding.EncodeToString(sign([]byte(data), key))
	return data + "|" + signature
}
