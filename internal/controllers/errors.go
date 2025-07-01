package controllers

import "errors"

var (
	ErrDataBindError = errors.New("wrong data format")
	ErrAuthRequired  = errors.New("auth is required")
)
