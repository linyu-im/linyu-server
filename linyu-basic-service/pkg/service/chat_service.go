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
	list, err := basicDao.ChatDao.ContactsChatList(db.RDB, userId)
	return list, err
}

func (s *chatService) SaveOrUpdateUserIncUnreadNum(userId, sessionId string, message *basicModel.Message) ([]string, error) {
	userIds := utils.Split1v1SessionID(sessionId)
	var toUserId string
	for _, id := range userIds {
		if id != userId {
			toUserId = id
			break
		}
	}
	err := ChatService.SaveOrUpdateIncUnreadNum(toUserId, userId, sessionId, message)
	if err != nil {
		return nil, err
	}
	//更新自己的会话
	err = ChatService.SaveOrUpdate(userId, toUserId, sessionId, message)
	if err != nil {
		return nil, err
	}
	return userIds, nil

}

func (s *chatService) SaveOrUpdateGroupIncUnreadNum(userId, sessionId string, message *basicModel.Message) ([]string, error) {
	memberIds := GroupService.GetMemberUserIdsByGroupId(sessionId)
	for _, id := range memberIds {
		if id != userId {
			_ = ChatService.SaveOrUpdateIncUnreadNum(id, sessionId, sessionId, message)
		}
	}
	err := ChatService.SaveOrUpdate(userId, sessionId, sessionId, message)
	if err != nil {
		return nil, err
	}
	return memberIds, nil
}

func (s *chatService) SaveOrUpdateIncUnreadNum(userId, peerId, sessionId string, message *basicModel.Message) error {
	chat, err := basicDao.ChatDao.GetChatByUserAndPeer(db.RDB, userId, peerId)
	if err != nil {
		return err
	}
	if chat == nil {
		return basicDao.ChatDao.Create(db.RDB, &basicModel.Chat{
			ID:             utils.GenerateSfIDString(),
			UserID:         userId,
			PeerID:         peerId,
			SessionID:      sessionId,
			LastMsgContent: message,
			Type:           constant.ChatType.User,
			UnreadNum:      1,
		})
	}
	chat.LastMsgContent = message
	chat.UnreadNum = chat.UnreadNum + 1
	return basicDao.ChatDao.Update(db.RDB, chat)
}

func (s *chatService) SaveOrUpdate(userId, peerId, sessionId string, message *basicModel.Message) error {
	chat, err := basicDao.ChatDao.GetChatByUserAndPeer(db.RDB, userId, peerId)
	if err != nil {
		return err
	}
	if chat == nil {
		return basicDao.ChatDao.Create(db.RDB, &basicModel.Chat{
			ID:             utils.GenerateSfIDString(),
			UserID:         userId,
			PeerID:         peerId,
			SessionID:      sessionId,
			LastMsgContent: message,
			Type:           constant.ChatType.User,
		})
	}
	chat.LastMsgContent = message
	return basicDao.ChatDao.Update(db.RDB, chat)
}

func (s *chatService) ChatCreate(userId string, param *basicParam.ChatCreateParam) (*basicModel.Chat, error) {
	if !constant.ChatType.Validate(param.ChatType) {
		return nil, errors.New("param.type-not-exist")
	}
	chat, _ := basicDao.ChatDao.GetChatByUserAndPeer(db.RDB, userId, param.PeerId)
	if chat != nil {
		return chat, nil
	}
	chat = &basicModel.Chat{
		ID:     utils.GenerateSfIDString(),
		UserID: userId,
		PeerID: param.PeerId,
		Type:   param.ChatType,
	}
	err := basicDao.ChatDao.Create(db.RDB, chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *chatService) SetTop(userId string, param *basicParam.ChatSetTopParam) error {
	err := basicDao.ChatDao.SetIsTopByIdAndUserId(db.RDB, param.IsTop, userId, param.ChatId)
	return err
}
