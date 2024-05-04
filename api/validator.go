package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/izsal/simple_bank/util"
)

// This code snippet is defining a custom validation function named `validCurrency` using the Go
// Playground Validator library. The function takes a `validator.FieldLevel` parameter and returns a
// boolean value.
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
