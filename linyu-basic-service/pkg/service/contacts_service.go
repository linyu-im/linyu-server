package service

var ContactsService = newContactsService()

func newContactsService() *contactsService {
	return &contactsService{}
}

type contactsService struct{}
