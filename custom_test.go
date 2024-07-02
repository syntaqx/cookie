package cookie

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type CustomType struct {
}

func CustomTypeFromString(value string) (CustomType, error) {
	return CustomType{}, nil
}

func CustomTypeErrorMaker(value string) (CustomType, error) {
	return CustomType{}, errors.New("just a big ol fail")
}

func TestWithCustomHandler(t *testing.T) {
	manager := NewManager(
		WithCustomHandler(reflect.TypeOf(CustomType{}), func(value string) (interface{}, error) {
			return CustomTypeFromString(value)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "cookie", Value: "test"})

	type MyStruct struct {
		Field CustomType `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := CustomType{}
	if dest.Field != expected {
		t.Errorf("Expected value '%s', but got '%s'", expected, dest.Field)
	}
}

func TestWithCustomHandler_HandlerErr(t *testing.T) {
	manager := NewManager(
		WithCustomHandler(reflect.TypeOf(CustomType{}), func(value string) (interface{}, error) {
			return CustomTypeErrorMaker(value)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "cookie", Value: "test"})

	type MyStruct struct {
		Field CustomType `cookie:"cookie"`
	}

	dest := &MyStruct{}
	err := manager.PopulateFromCookies(req, dest)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}
