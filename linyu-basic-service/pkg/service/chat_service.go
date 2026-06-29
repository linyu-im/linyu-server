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

func (s *chatService) SaveOrUpdateUserIncUnreadNum(userId string, message *basicModel.Message) ([]string, error) {
	//更新自己的会话
	err := ChatService.SaveOrUpdate(userId, message.ToID, message)
	if err != nil {
		return nil, err
	}
	if userId == message.ToID {
		return []string{userId}, nil
	}
	err = ChatService.SaveOrUpdateIncUnreadNum(message.ToID, userId, message)
	if err != nil {
		return nil, err
	}
	return []string{userId, message.ToID}, nil

}

func (s *chatService) SaveOrUpdateGroupIncUnreadNum(userId string, message *basicModel.Message) ([]string, error) {
	memberIds := GroupService.GetMemberUserIdsByGroupId(message.SessionID)
	for _, id := range memberIds {
		if id != userId {
			_ = ChatService.SaveOrUpdateIncUnreadNum(id, message.SessionID, message)
		}
	}
	err := ChatService.SaveOrUpdate(userId, message.SessionID, message)
	if err != nil {
		return nil, err
	}
	return memberIds, nil
}

func (s *chatService) SaveOrUpdateIncUnreadNum(userId, peerId string, message *basicModel.Message) error {
	chat, err := basicDao.ChatDao.GetChatByUserAndPeer(db.RDB, userId, peerId)
	if err != nil {
		return err
	}
	if chat == nil {
		return basicDao.ChatDao.Create(db.RDB, &basicModel.Chat{
			ID:             utils.GenerateSfIDString(),
			UserID:         userId,
			PeerID:         peerId,
			SessionID:      message.SessionID,
			LastMsgContent: message,
			SceneType:      message.SceneType,
			UnreadNum:      1,
		})
	}
	chat.LastMsgContent = message
	chat.UnreadNum = chat.UnreadNum + 1
	return basicDao.ChatDao.Update(db.RDB, chat)
}

func (s *chatService) SaveOrUpdate(userId, peerId string, message *basicModel.Message) error {
	chat, err := basicDao.ChatDao.GetChatByUserAndPeer(db.RDB, userId, peerId)
	if err != nil {
		return err
	}
	if chat == nil {
		return basicDao.ChatDao.Create(db.RDB, &basicModel.Chat{
			ID:             utils.GenerateSfIDString(),
			UserID:         userId,
			PeerID:         peerId,
			SessionID:      message.SessionID,
			LastMsgContent: message,
			SceneType:      message.SceneType,
		})
	}
	chat.LastMsgContent = message
	return basicDao.ChatDao.Update(db.RDB, chat)
}

func (s *chatService) ChatCreate(userId string, param *basicParam.ChatCreateParam) (*basicModel.Chat, error) {
	if !constant.SceneType.Validate(param.SceneType) {
		return nil, errors.New("param.type-not-exist")
	}
	chat, _ := basicDao.ChatDao.GetChatByUserAndPeer(db.RDB, userId, param.PeerId)
	if chat != nil {
		return chat, nil
	}
	chat = &basicModel.Chat{
		ID:        utils.GenerateSfIDString(),
		UserID:    userId,
		SessionID: MessageService.GetSessionIdByPeerIdAndSceneType(userId, param.PeerId, param.SceneType),
		PeerID:    param.PeerId,
		SceneType: param.SceneType,
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

func (s *chatService) SetMute(userId string, param *basicParam.ChatMuteParam) error {
	return basicDao.ChatDao.SetIsMuteByIdAndUserId(db.RDB, param.IsMute, userId, param.ChatId)
}

func (s *chatService) ChatDelete(userId string, param *basicParam.ChatDeleteParam) error {
	return basicDao.ChatDao.DeleteByIdAndUserId(db.RDB, userId, param.ChatId)
}

func (s *chatService) MarkRead(userId string, param *basicParam.ChatMarkReadParam) error {
	return basicDao.ChatDao.ClearUnreadByIdAndUserId(db.RDB, userId, param.ChatId)
}
