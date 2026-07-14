package param

type ContactsFriendDeleteParam struct {
	UserId string `json:"userId" binding:"required"`
}

type ContactsIsFriendParam struct {
	UserId string `json:"userId" binding:"required"`
}

type ContactsUpdateRemarkParam struct {
	PeerId string `json:"peerId" binding:"required"`
	Remark string `json:"remark"`
}

type ContactsUpdateTagParam struct {
	PeerId string `json:"peerId" binding:"required"`
	Tag    string `json:"tag"`
}
