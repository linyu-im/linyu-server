package service

import (
	"errors"
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

var ChatService = newChatService()

func newChatService() *chatService {
	return &chatService{}
}

type chatService struct{}

func (s *chatService) ChatList(userId string) ([]*basicModel.Chat, error) {
	list, err := basicDao.ChatDao.ChatList(db.MysqlDB, userId)
	return list, err
}

func (s *chatService) UpdateUserChat(userId, peerId, chatType string, message *basicModel.Message) error {
	chat, err := basicDao.ChatDao.GetChatByUserAndPeer(db.MysqlDB, userId, peerId)
	if err != nil {
		return err
	}
	if chat == nil {
		return basicDao.ChatDao.Create(db.MysqlDB, &basicModel.Chat{
			ID:             utils.GenerateSfIDString(),
			UserID:         userId,
			PeerID:         peerId,
			LastMsgContent: message,
			Type:           chatType,
		})
	}
	chat.LastMsgContent = message
	chat.UnreadNum = chat.UnreadNum + 1
	return basicDao.ChatDao.Update(db.MysqlDB, chat)
}

func (s *chatService) ChatCreate(userId string, param *basicParam.ChatCreateParam) (*basicModel.Chat, error) {
	if !constant.ChatType.Validate(param.ChatType) {
		return nil, errors.New("param.type-not-exist")
	}
	chat, _ := basicDao.ChatDao.GetChatByUserAndPeer(db.MysqlDB, userId, param.PeerId)
	if chat != nil {
		return chat, nil
	}
	chat = &basicModel.Chat{
		ID:     utils.GenerateSfIDString(),
		UserID: userId,
		PeerID: param.PeerId,
		Type:   param.ChatType,
	}
	err := basicDao.ChatDao.Create(db.MysqlDB, chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
