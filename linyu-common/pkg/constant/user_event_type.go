package constant

type userEventType struct {
	Login string //登录事件
}

var UserEventType = userEventType{
	Login: "login",
}
