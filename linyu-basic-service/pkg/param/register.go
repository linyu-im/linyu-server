package param

type EmailCodeParam struct {
	Email string `json:"email" binding:"required,email"`
}

type EmailRegisterParam struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}
