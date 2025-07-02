package requests

type RefreshSession struct {
	AccessToken  string `json:"accessToken" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOiIxMjM0NTY3ODkwIiwidXNlcklkIjoiMSIsImlzcyI6IlRPRE8ifQ.K-6Tzcaoae1Cj7jbIMalrtsLXZFrAlg_F_XLegWGo7o"`
	RefreshToken string `json:"refreshToken" binding:"required" example:"$2a$10$9UKV92GI6504S7RpPPZApe1Llp3fyOS7TH4tQC9ty6OQLxcjIT8uC"`
}

type CloseSessions struct {
	AccessToken string `json:"accessToken" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOiIxMjM0NTY3ODkwIiwidXNlcklkIjoiMSIsImlzcyI6IlRPRE8ifQ.K-6Tzcaoae1Cj7jbIMalrtsLXZFrAlg_F_XLegWGo7o"`
}
