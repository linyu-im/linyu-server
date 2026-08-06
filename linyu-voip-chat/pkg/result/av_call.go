package result

// AvCallWsContent 音视频通话 WS 推送内容
type AvCallWsContent struct {
	Action    string   `json:"action"` // invite / hangup / change
	SessionId string   `json:"sessionId"`
	FromId    string   `json:"fromId"`
	ToUserIds []string `json:"toUserIds"` // 被邀请/通知的用户
	CallType  string   `json:"callType"`  // audio / video
	SceneType string   `json:"sceneType"`
}

// AvCallInviteResult 通话邀请结果
type AvCallInviteResult struct {
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId,omitempty"`  // 单聊对端用户 id
	GroupId   string `json:"groupId,omitempty"` // 群聊 id
	CallType  string `json:"callType"`
	SceneType string `json:"sceneType"`
}
