package structValidator

import (
	"fmt"
	"github.com/go-playground/validator"
)

var Validator = validator.New()

func ValidateStruct(s interface{}) error {
	err := Validator.Struct(s)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			for _, currentField := range errors {
				return fmt.Errorf("field '%s' validation failed on the '%s' tag", currentField.Field(), currentField.Tag())
			}
		}
		return err
	}
	return nil
}
