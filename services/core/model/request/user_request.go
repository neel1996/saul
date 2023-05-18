package request

type UserLoginRequest struct {
	UserId string `json:"userId" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
}
