package param

import basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"

type ContactsSearchParam struct {
	Keyword string `json:"keyword" binding:"required"`
}

type ContactsSearchResult struct {
	Friends []*basicModel.Contacts `json:"friends"`
	Groups  []*basicModel.Contacts `json:"groups"`
}

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
