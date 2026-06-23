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
