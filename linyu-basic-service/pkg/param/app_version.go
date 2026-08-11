package param

type AppVersionCheckParam struct {
	Platform string `json:"platform" binding:"required"`
}
