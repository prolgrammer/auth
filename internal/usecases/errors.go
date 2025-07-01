package usecases

import "errors"

var ErrInvalidEntity = errors.New("validation error")
var ErrEntityNotFound = errors.New("entity not found")
var ErrEntityAlreadyExists = errors.New("entity already exists")

var ErrWrongPassword = errors.New("wrong password")

var ErrAccessTokenExpired = errors.New("access token is expired")
var ErrRefreshTokenExpired = errors.New("refresh token is expired")
var ErrNotAValidAccessToken = errors.New("invalid access token")
var ErrNotAValidRefreshToken = errors.New("invalid refresh token")
