package param

type SendMessageToUserParam struct {
	ToUserId string `json:"toUserId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type SendMessageToGroupParam struct {
	ToGroupId string `json:"toGroupId" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
