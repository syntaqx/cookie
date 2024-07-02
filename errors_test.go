package cookie

import (
	"reflect"
	"testing"
)

func TestErrUnsupportedType_Error(t *testing.T) {
	err := &ErrUnsupportedType{Type: reflect.TypeOf(0)}
	expected := "cookie: unsupported type: int"

	if err.Error() != expected {
		t.Errorf("Expected error message '%s', but got '%s'", expected, err.Error())
	}
}
