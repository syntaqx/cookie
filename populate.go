package cookie

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// PopulateFromCookies populates a struct with cookie values.
func (m *Manager) PopulateFromCookies(r *http.Request, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrNonNilPointerRequired
	}
	v = v.Elem()

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("cookie")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]
		signed := false
		unsigned := false
		omitempty := false

		for _, part := range parts[1:] {
			if part == "signed" {
				signed = true
			} else if part == "unsigned" {
				unsigned = true
			} else if part == "omitempty" {
				omitempty = true
			}
		}

		var value string
		var err error
		if signed && !unsigned {
			value, err = m.GetSigned(r, name)
		} else {
			value, err = m.Get(r, name)
		}
		if err != nil {
			if err == http.ErrNoCookie && omitempty {
				continue
			}
			return err
		}

		fieldVal := v.Field(i)

		// TODO: Is this necessary? How can I test it?
		// if !fieldVal.CanSet() {
		// 	continue
		// }

		err = m.setFieldValue(fieldVal, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// setFieldValue sets the value of a struct field based on its type.
func (m *Manager) setFieldValue(fieldVal reflect.Value, value string) error {
	if handler, ok := m.customHandlers[fieldVal.Type()]; ok {
		customValue, err := handler(value)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(customValue))
		return nil
	}

	switch fieldVal.Kind() {
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		fieldVal.SetBool(boolVal)
	case reflect.String:
		fieldVal.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, fieldVal.Type().Bits())
		if err != nil {
			return err
		}
		fieldVal.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, fieldVal.Type().Bits())
		if err != nil {
			return err
		}
		fieldVal.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, fieldVal.Type().Bits())
		if err != nil {
			return err
		}
		fieldVal.SetFloat(floatVal)
	case reflect.Slice:
		switch fieldVal.Type().Elem().Kind() {
		case reflect.String:
			fieldVal.Set(reflect.ValueOf(strings.Split(value, ",")))
		case reflect.Int:
			strSlice := strings.Split(value, ",")
			intSlice := make([]int, len(strSlice))
			for i, str := range strSlice {
				intVal, err := strconv.Atoi(str)
				if err != nil {
					return err
				}
				intSlice[i] = intVal
			}
			fieldVal.Set(reflect.ValueOf(intSlice))
		}
	case reflect.Struct:
		if fieldVal.Type() == reflect.TypeOf(time.Time{}) {
			timeVal, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			fieldVal.Set(reflect.ValueOf(timeVal))
		}
	default:
		return &ErrUnsupportedType{Type: fieldVal.Type()}
	}
	return nil
}
