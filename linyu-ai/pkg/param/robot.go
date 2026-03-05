package param

type RobotAnswersParam struct {
	SessionId string `json:"sessionId"`
	RobotId   string `json:"robotId"  binding:"required"`
	Question  string `json:"question" binding:"required"`
}
