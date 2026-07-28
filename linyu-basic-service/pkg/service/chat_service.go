package service

import (
	"errors"
	"fmt"

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
	ids := utils.Split1v1SessionID(message.SessionID)
	toId := ids[0]
	if toId == userId {
		toId = ids[1]
	}
	//更新自己的会话
	err := ChatService.SaveOrUpdate(userId, toId, message)
	if err != nil {
		return nil, err
	}
	if userId == toId {
		return []string{userId}, nil
	}
	err = ChatService.SaveOrUpdateIncUnreadNum(toId, userId, message)
	if err != nil {
		return nil, err
	}
	return []string{userId, toId}, nil
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
	isActiveSession := s.isUserActiveSession(userId, message.SessionID)
	if isActiveSession {
		chat.UnreadNum = 0
	} else {
		chat.UnreadNum++
	}
	return basicDao.ChatDao.Update(db.RDB, chat)
}

func (s *chatService) isUserActiveSession(userId string, sessionId string) bool {
	if sessionId == "" {
		return false
	}

	devices, err := db.CacheDB.SMembers(fmt.Sprintf(constant.RedisKey.UserOnline, userId))
	if err != nil || len(devices) == 0 {
		return false
	}
	for _, device := range devices {
		var activeSessionIds []string
		err := db.CacheDB.GetObject(
			fmt.Sprintf(constant.RedisKey.UserActiveSession, userId, device),
			&activeSessionIds,
		)
		if err != nil {
			continue
		}
		for _, activeSessionId := range activeSessionIds {
			if activeSessionId == sessionId {
				return true
			}
		}
	}
	return false
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

func (s *chatService) SetActiveSession(userId string, device string, activeSessionIds []string) error {
	key := fmt.Sprintf(constant.RedisKey.UserActiveSession, userId, device)
	if len(activeSessionIds) == 0 {
		return db.CacheDB.Del(key)
	}
	return db.CacheDB.SetObject(key, activeSessionIds, 0)
}
