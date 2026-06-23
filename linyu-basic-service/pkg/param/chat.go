package param

type ChatCreateParam struct {
	PeerId    string `json:"peerId" binding:"required"`
	SceneType string `json:"sceneType" binding:"required"`
}

type ChatSetTopParam struct {
	ChatId string `json:"chatId" binding:"required"`
	IsTop  bool   `json:"isTop"`
}

type ChatMuteParam struct {
	ChatId string `json:"chatId" binding:"required"`
	IsMute bool   `json:"isMute"`
}

type ChatDeleteParam struct {
	ChatId string `json:"chatId" binding:"required"`
}

type ChatMarkReadParam struct {
	ChatId string `json:"chatId" binding:"required"`
}
