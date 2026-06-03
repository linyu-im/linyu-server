package param

type ApplyAddFriendParam struct {
	PeerId      string `json:"peerId" binding:"required"`
	ApplySource string `json:"applySource" binding:"required"`
	Describe    string `json:"describe"`
}

type ApplyAgreeFriendParam struct {
	ApplyId string `json:"applyId" binding:"required"`
}

type ApplyRejectParam struct {
	ApplyId string `json:"applyId" binding:"required"`
}

type ApplyCancelParam struct {
	ApplyId string `json:"applyId" binding:"required"`
}
