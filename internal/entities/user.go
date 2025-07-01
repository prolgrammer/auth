package entities

import "time"

type User struct {
	Id               string
	Email            Email
	Password         Password
	RegistrationDate time.Time
}

func NewUser(email string, password string) User {
	var result User

	result.Email = Email(email)
	result.Password = Password(password)

	return result
}

func (a User) Validate() error {
	err := a.Email.Validate()
	if err != nil {
		return err
	}
	err = a.Password.Validate()
	if err != nil {
		return err
	}
	return nil
}
