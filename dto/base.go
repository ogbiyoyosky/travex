package dto

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Property string
	Tag      string
	Value    string
	Message  string
}

var Validator = validator.New()
