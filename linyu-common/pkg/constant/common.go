package constant

// ------------------公共常量------------------
type vectorCollection struct {
	LongTermMemory string
}

var VectorCollection = &vectorCollection{
	LongTermMemory: "ltm",
}

type sceneType struct {
	User  string //用户
	Group string //群
}

// SceneType 消息场景类型
var SceneType = sceneType{
	User:  "user",
	Group: "group",
}

func (c sceneType) Validate(v string) bool {
	switch v {
	case c.User, c.Group:
		return true
	default:
		return false
	}
}

type avCallType struct {
	Audio string // 音频通话
	Video string // 视频通话
}

// AvCallType 音视频通话类型
var AvCallType = avCallType{
	Audio: "audio",
	Video: "video",
}

func (c avCallType) Validate(v string) bool {
	switch v {
	case c.Audio, c.Video:
		return true
	default:
		return false
	}
}

type avCallAction struct {
	Invite string // 邀请
	Hangup string // 挂断
	Change string // 成员变更
}

// AvCallAction 音视频通话动作
var AvCallAction = avCallAction{
	Invite: "invite",
	Hangup: "hangup",
	Change: "change",
}

func (c avCallAction) Validate(v string) bool {
	switch v {
	case c.Invite, c.Hangup, c.Change:
		return true
	default:
		return false
	}
}

type avCallStatus struct {
	Ended    string // 结束（已接通）
	Missed   string // 未接通
	Rejected string // 拒绝
	Canceled string // 取消
	Calling  string // 通话中
}

// AvCallStatus 通话记录状态
var AvCallStatus = avCallStatus{
	Ended:    "ended",
	Missed:   "missed",
	Rejected: "rejected",
	Canceled: "canceled",
	Calling:  "calling",
}

func (c avCallStatus) Validate(v string) bool {
	switch v {
	case c.Ended, c.Missed, c.Rejected, c.Canceled, c.Calling:
		return true
	default:
		return false
	}
}
