package param

type ContactsRelDeleteParam struct {
	ContactsId string `json:"contactsId" binding:"required"`
}

type ContactsIsFriendParam struct {
	UserId string `json:"userId" binding:"required"`
}
