package request

type RegisterUserReq struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type LoginUserReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
