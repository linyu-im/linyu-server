package param

type AvCallFriendInviteParam struct {
	UserId   string `json:"userId" binding:"required"`
	CallType string `json:"callType" binding:"required"` // audio / video
}

// AvCallGroupInviteParam 群聊通话邀请
type AvCallGroupInviteParam struct {
	GroupId  string   `json:"groupId" binding:"required"`
	UserIds  []string `json:"userIds" binding:"required,min=1"`
	CallType string   `json:"callType" binding:"required"` // audio / video
}

// AvCallUserHangupParam 单聊通话挂断
type AvCallUserHangupParam struct {
	UserId string `json:"userId" binding:"required"`
}

// AvCallGroupHangupParam 群聊通话挂断
type AvCallGroupHangupParam struct {
	GroupId string   `json:"groupId" binding:"required"`
	UserIds []string `json:"userIds" binding:"required,min=1"`
}
