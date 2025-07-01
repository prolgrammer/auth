package responses

type SignIn struct {
	Id      string  `json:"id" example:"2"`
	Session Session `json:"session"`
}

func NewSignIn(id string, session Session) SignIn {
	return SignIn{Id: id, Session: session}
}
