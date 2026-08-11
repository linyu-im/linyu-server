package param

type PwdLoginParam struct {
	Account     string `json:"account" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Platform    string `json:"platform" binding:"required"`
	VersionCode int    `json:"versionCode"`
}

type LdapLoginParam struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Platform    string `json:"platform" binding:"required"`
	VersionCode int    `json:"versionCode"`
}

type Oauth2LoginParam struct {
	Type        string `json:"type" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Platform    string `json:"platform" binding:"required"`
	VersionCode int    `json:"versionCode"`
}
