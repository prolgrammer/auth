package responses

import "time"

type User struct {
	Id               string
	Email            string
	RegistrationDate time.Time
}
