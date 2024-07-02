package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
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
