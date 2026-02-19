package validate

import "github.com/go-playground/validator/v10"

var validatorInstance *validator.Validate

// Validator returns a singleton validator instance.
func Validator() *validator.Validate {
	if validatorInstance == nil {
		v := validator.New()
		validatorInstance = v
	}
	return validatorInstance
}

