package responses

type Session struct {
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOiIxMjM0NTY3ODkwIiwidXNlcklkIjoiMSIsImlzcyI6IlRPRE8ifQ.K-6Tzcaoae1Cj7jbIMalrtsLXZFrAlg_F_XLegWGo7o"`
	RefreshToken string `json:"refreshToken" example:"$2a$10$9UKV92GI6504S7RpPPZApe1Llp3fyOS7TH4tQC9ty6OQLxcjIT8uC"`
	ExpiresAt    int64  `json:"expiresAt" example:"1592572800"`
}

func NewSession(
	accessToken string,
	refreshToken string,
	expiresAt int64) Session {
	return Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
}
