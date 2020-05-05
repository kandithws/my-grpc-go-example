package utils

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrRequestBindingError = errors.New("Binding Error")
)

func BindJSON(in interface{}, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return ErrRequestBindingError
	}

	if reflect.ValueOf(out).IsNil() {
		return ErrRequestBindingError
	}

	byt, _ := json.Marshal(in)
	if err := json.Unmarshal(byt, out); err != nil {
		return err
	}

	return nil
}
