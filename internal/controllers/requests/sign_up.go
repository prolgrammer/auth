package requests

type SignUp struct {
	Email    string `json:"email" binding:"required" example:"example@mail.ru"`
	Password string `json:"password" binding:"required" example:"123superPassword"`
}
