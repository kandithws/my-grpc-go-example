package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(model interface{}) error
}

var vObj *vObject

func NewValidator() Validator {
	// Singleton
	if vObj == nil {
		vObj = &vObject{validate: validator.New()}
	}

	return vObj
}

type vObject struct {
	validate *validator.Validate
}

func (v *vObject) Validate(model interface{}) error {
	err := v.validate.Struct(model)
	if err != nil {
		// TODO custom validations
		return err
	}
	return nil
}
