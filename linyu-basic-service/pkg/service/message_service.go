package service

import (
	"fmt"
	"strings"

	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event/eventbus"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct{}

func (s messageService) SendMessageToSession(currentUserId string, message *basicModel.Message) (*basicModel.Message, error) {
	var toUserIds []string
	switch message.SceneType {
	case constant.SceneType.User:
		toUserIds, _ = ChatService.SaveOrUpdateUserIncUnreadNum(currentUserId, message)
	case constant.SceneType.Group:
		toUserIds, _ = ChatService.SaveOrUpdateGroupIncUnreadNum(currentUserId, message)
	}
	err := basicDao.MessageDao.Create(db.RDB, message)
	if err != nil {
		return nil, err
	}
	//发送消息
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: currentUserId,
		ToUserIds:  toUserIds,
		Data: &event.WsData{
			SeqId:   message.ID,
			Type:    constant.WsDataType.Message,
			Content: message,
		},
	})
	return message, nil
}

func (s messageService) SendMessage(userId string, param *basicParam.SendMessageToUserParam) (*basicModel.Message, error) {
	content, err := basicModel.ParseMsgContent(param.MsgType, param.Content)
	if err != nil {
		return nil, err
	}
	sceneType, toId, err := s.VerifySessionSceneType(userId, param.SessionId)
	message := &basicModel.Message{
		ID:             utils.GenerateSfIDString(),
		SessionID:      param.SessionId,
		FromID:         userId,
		ToID:           toId,
		SceneType:      sceneType,
		Content:        param.Content,
		KeywordContent: content.GetKeywordContent(),
		MsgType:        param.MsgType,
		FromType:       constant.MessageFromType.User,
		IsShowTime:     param.IsShowTime,
		QuoteMsgId:     param.QuoteMsgId,
	}
	if err != nil {
		message.CreatedAt = localtime.Now()
		message.Status = "failed"
		message.FailReason = err.Error()
		_ = ChatService.SaveOrUpdate(userId, message.ToID, message)
		return message, nil
	}
	//分发消息
	return s.SendMessageToSession(userId, message)
}

func (s messageService) GetMessageBySessionId(sessionId string, num int) []*basicModel.Message {
	return basicDao.MessageDao.GetLatestMessagesBySessionID(db.RDB, sessionId, num)
}

func (s messageService) VerifySessionSceneType(userId, session string) (string, string, error) {
	if strings.Contains(session, userId) {
		ids := utils.Split1v1SessionID(session)
		toId := ids[0]
		if toId == userId {
			toId = ids[1]
		}
		if toId == userId {
			return constant.SceneType.User, userId, nil
		}
		// 判断自己是否是对方好友
		if ContactsService.IsFriend(toId, userId) {
			return constant.SceneType.User, toId, nil
		} else {
			return constant.SceneType.User, toId, fmt.Errorf("basic.contacts.you-no-other-friend")
		}
	} else if GroupService.IsGroupMember(session, userId) {
		return constant.SceneType.Group, session, nil
	}
	return constant.SceneType.Group, session, fmt.Errorf("param.error")
}

func (s messageService) GetSessionIdByPeerIdAndSceneType(userId, peerId, sceneType string) string {
	if sceneType == constant.SceneType.User {
		return utils.Generate1v1SessionID(userId, peerId)
	} else {
		return peerId
	}
}

func (s messageService) MessageList(userId string, param *basicParam.MessageListParam) ([]*basicModel.Message, error) {
	_, _, err := s.VerifySessionSceneType(userId, param.SessionId)
	if err != nil {
		return nil, err
	}
	messages, err := basicDao.MessageDao.ListMessagesBySessionIDSinceMsgID(db.RDB, param.SessionId, param.SinceMsgId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s messageService) MessagePage(userId string, param *basicParam.MessagePageParam) (*response.PageResult[*basicModel.Message], error) {
	var sessionId string
	if ContactsService.IsFriendBothOr(userId, param.ToId) {
		sessionId = utils.Generate1v1SessionID(userId, param.ToId)
	} else if GroupService.IsGroupMember(param.ToId, userId) {
		sessionId = param.ToId
	} else {
		return nil, fmt.Errorf("param.error")
	}
	messages, total, err := basicDao.MessageDao.PageMessagesBySessionID(db.RDB, sessionId, param.Page, param.PageSize)
	if err != nil {
		return nil, err
	}
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.PageSize <= 0 {
		param.PageSize = 10
	}
	totalPage := int(total) / param.PageSize
	if int(total)%param.PageSize > 0 {
		totalPage++
	}
	result := &response.PageResult[*basicModel.Message]{
		Records:   messages,
		Total:     total,
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalPage: totalPage,
	}
	return result, nil
}

func (s messageService) ForwardMessage(currentUserId string, param *basicParam.ForwardMessageParam) error {
	content, err := basicModel.ParseMsgContent(param.Message.MsgType, param.Message.Content)
	if err != nil {
		return err

	}
	for _, peerInfo := range param.ToPeerInfo {
		message := &basicModel.Message{
			ID:             utils.GenerateSfIDString(),
			FromID:         currentUserId,
			ToID:           peerInfo.PeerId,
			Content:        param.Message.Content,
			MsgType:        param.Message.MsgType,
			KeywordContent: content.GetKeywordContent(),
			SceneType:      peerInfo.PeerSceneType,
			FromType:       constant.MessageFromType.User,
		}
		switch peerInfo.PeerSceneType {
		case constant.SceneType.User:
			message.SessionID = utils.Generate1v1SessionID(currentUserId, peerInfo.PeerId)
		case constant.SceneType.Group:
			message.SessionID = peerInfo.PeerId
		default:
			return fmt.Errorf("param.error")
		}
		_, _ = s.SendMessageToSession(currentUserId, message)
	}
	return nil
}
