package param

type AccountCheckParam struct {
	Account string `json:"account" binding:"required,account"`
}

type EmailCodeParam struct {
	Email string `json:"email" binding:"required,email"`
}

type EmailRegisterParam struct {
	Email           string `json:"email" binding:"required,email"`
	Code            string `json:"code" binding:"required"`
	Account         string `json:"account" binding:"required,account"`
	Password        string `json:"password" binding:"required,password"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}
