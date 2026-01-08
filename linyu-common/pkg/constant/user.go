package constant

type userStatus struct {
	Active string
	Banned string
}

var UserStatus = userStatus{
	Active: "active",
	Banned: "banned",
}
