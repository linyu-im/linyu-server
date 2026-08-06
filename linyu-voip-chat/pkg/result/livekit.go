package result

// LivekitRoomUser 房间内参与者信息
type LivekitRoomUser struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	State    int32  `json:"state"` // LiveKit ParticipantInfo.State
}
