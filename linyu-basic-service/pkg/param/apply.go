package param

type ApplyAddContactsParam struct {
	PeerId   string `json:"peerId" binding:"required"`
	Describe string `json:"describe"`
}

type ApplyAgreeContactsParam struct {
	ApplyId string `json:"applyId" binding:"required"`
}

type ApplyRejectParam struct {
	ApplyId string `json:"applyId" binding:"required"`
}

type ApplyCancelParam struct {
	ApplyId string `json:"applyId" binding:"required"`
}
