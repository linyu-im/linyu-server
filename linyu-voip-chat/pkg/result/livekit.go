package result

// LivekitInfo LiveKit 服务信息
type LivekitInfo struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
}

// LivekitRoomUser 房间内参与者信息
type LivekitRoomUser struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	State    int32  `json:"state"` // LiveKit ParticipantInfo.State
}
