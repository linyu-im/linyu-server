package param

type ContactsRelDeleteParam struct {
	ContactsId string `json:"contactsId" binding:"required"`
}
