package param

type ChatCreateParam struct {
	PeerId   string `json:"peerId" binding:"required"`
	ChatType string `json:"chatType" binding:"required"`
}

type ChatSetTopParam struct {
	ChatId string `json:"chatId" binding:"required"`
	IsTop  bool   `json:"isTop"`
}
