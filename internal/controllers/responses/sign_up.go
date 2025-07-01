package responses

type SignUp struct {
	Id      string  `json:"id"  binding:"required" example:"2"`
	Session Session `json:"session"  binding:"required"`
}

func NewSignUp(id string, session Session) SignUp {
	return SignUp{Id: id, Session: session}
}
