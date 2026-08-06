package param

type LivekitTokenParam struct {
	SessionId string `json:"sessionId" binding:"required"`
}

type LivekitUserTokenParam struct {
	SessionId string `json:"sessionId" binding:"required"`
}
