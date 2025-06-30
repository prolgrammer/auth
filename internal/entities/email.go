package entities

import (
	"errors"
	"fmt"
)

type Email string

const (
	minEmailLen = 6
	maxEmailLen = 254
)

func (e Email) Validate() error {
	if !e.validateLength() {
		return errors.New(fmt.Sprintf("email length can't be less than %d OR more %d", minEmailLen, maxEmailLen))
	}
	if !e.validateFormat() {
		return errors.New("wrong email format")
	}
	return nil
}

func (e Email) validateLength() bool {
	return stringLengthInRange(string(e), minEmailLen, maxEmailLen)
}

func (e Email) validateFormat() bool {
	return isEmail(string(e))
}
