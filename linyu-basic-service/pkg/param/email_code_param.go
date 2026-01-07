package param

type EmailCodeParam struct {
	Email string `json:"email" binding:"required,email"`
}
