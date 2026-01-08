package param

type SendMessageParam struct {
	ToUserId string `json:"toUserId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}
