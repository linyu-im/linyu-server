package param

type ResetPasswordEmailCodeParam struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordParam struct {
	Email           string `json:"email" binding:"required,email"`
	Code            string `json:"code" binding:"required"`
	Password        string `json:"password" binding:"required,password"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}
