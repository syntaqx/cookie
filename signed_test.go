package cookie

import (
	"crypto/hmac"
	"testing"
)

func TestSignVerify(t *testing.T) {
	data := []byte("example data")
	key := []byte("example key")

	expectedSignature := []byte{
		143, 44, 153, 63, 34, 126, 71, 71, 60, 146, 137, 245, 195, 249, 153, 4,
		171, 247, 130, 233, 162, 23, 163, 57, 160, 123, 76, 145, 124, 34, 222, 55,
	}

	signature := sign(data, key)

	if !hmac.Equal(signature, expectedSignature) {
		t.Error("sign failed to generate the expected signature")
	}

	if !verify(data, signature, key) {
		t.Error("verify failed to validate the signature")
	}
}
