package param

type PwdLoginParam struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LdapLoginParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Oauth2LoginParam struct {
	Type string `json:"type" binding:"required"`
	Code string `json:"code" binding:"required"`
}
