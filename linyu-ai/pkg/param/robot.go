package param

type RobotAnswersParam struct {
	PeerId    string `json:"peerId" binding:"required"`
	SceneType string `json:"sceneType" binding:"required"`
	RobotId   string `json:"robotId"  binding:"required"`
	Question  string `json:"question" binding:"required"`
}

type GetRobotAvatarParam struct {
	RobotId string `json:"robotId" binding:"required"`
}
