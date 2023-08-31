package api

import (
	"github.com/go-playground/validator/v10"
	"tutorial.sqlc.dev/app/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	//FieldLevel is an interface that contains all info and helper functions to validate a field.
	if currency, ok := fieldLevel.Field().Interface().(string); ok { //get the value of a field.
		return util.IsSupportedCurrency(currency)
	}
	//the field isn't a string.
	return false

}
