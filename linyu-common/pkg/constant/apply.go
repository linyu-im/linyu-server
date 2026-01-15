package constant

type applyType struct {
	AddContacts string // 添加联系人申请
	JoinGroup   string // 进群申请
}

// ApplyType 申请类型
var ApplyType = applyType{
	AddContacts: "addContacts",
	JoinGroup:   "joinGroup",
}

func (c applyType) Validate(v string) bool {
	switch v {
	case c.AddContacts, c.JoinGroup:
		return true
	default:
		return false
	}
}

type applyStatus struct {
	Wait   string //等待
	Agree  string //同意
	Cancel string //取消
	Reject string //拒绝
}

// ApplyStatus 申请状态
var ApplyStatus = applyStatus{
	Wait:   "wait",
	Agree:  "agree",
	Cancel: "cancel",
	Reject: "reject",
}
