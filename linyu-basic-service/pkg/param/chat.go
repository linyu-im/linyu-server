package param

type ChatCreateParam struct {
	PeerId   string `json:"peerId" binding:"required"`
	ChatType string `json:"chatType" binding:"required"`
}
