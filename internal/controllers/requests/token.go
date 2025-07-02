package requests

type GetTokenRequest struct {
	UserId string `json:"user_id" binding:"required" example:"e1e25658-3817-4051-8d0d-d13d575f08a4"`
}
