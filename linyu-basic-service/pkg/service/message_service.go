package service

import (
	"fmt"

	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event/eventbus"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct{}

func (s messageService) SendMessageToSession(currentUserId string, sessionId string, msgScene string, message *basicModel.Message) (*basicModel.Message, error) {
	var toUserIds []string
	switch msgScene {
	case constant.MessageScene.User:
		toUserIds, _ = ChatService.SaveOrUpdateUserIncUnreadNum(currentUserId, sessionId, message)
	case constant.MessageScene.Group:
		toUserIds, _ = ChatService.SaveOrUpdateGroupIncUnreadNum(currentUserId, sessionId, message)
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

func (s messageService) SendMessageToUser(userId string, param *basicParam.SendMessageToUserParam) (*basicModel.Message, error) {
	if _, err := basicModel.ParseMsgContent(param.MsgType, param.Content); err != nil {
		return nil, err
	}
	sessionId := utils.Generate1v1SessionID(userId, param.ToUserId)
	message := &basicModel.Message{
		ID:        utils.GenerateSfIDString(),
		SessionID: sessionId,
		FromID:    userId,
		ToID:      param.ToUserId,
		MsgScene:  constant.MessageScene.User,
		Content:   param.Content,
		MsgType:   param.MsgType,
		FromType:  constant.MessageFromType.User,
	}
	//分发消息
	return s.SendMessageToSession(userId, sessionId, constant.MessageScene.User, message)
}

func (s messageService) SendMessageToGroup(userId string, param *basicParam.SendMessageToGroupParam) (*basicModel.Message, error) {
	if _, err := basicModel.ParseMsgContent(param.MsgType, param.Content); err != nil {
		return nil, err
	}
	message := &basicModel.Message{
		ID:        utils.GenerateSfIDString(),
		SessionID: param.ToGroupId,
		FromID:    userId,
		ToID:      param.ToGroupId,
		MsgScene:  constant.MessageScene.Group,
		Content:   param.Content,
		MsgType:   constant.MessageType.Text,
		FromType:  constant.MessageFromType.User,
	}
	//分发消息
	return s.SendMessageToSession(userId, param.ToGroupId, constant.MessageScene.Group, message)
}

func (s messageService) GetMessageBySessionId(sessionId string, num int) []*basicModel.Message {
	return basicDao.MessageDao.GetLatestMessagesBySessionID(db.RDB, sessionId, num)
}

func (s messageService) MessagePage(userId string, param *basicParam.MessagePageParam) (*response.PageResult[*basicModel.Message], error) {
	var sessionId string
	if ContactsService.IsContacts(userId, param.ToId) {
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
	for _, peerInfo := range param.ToPeerInfo {
		message := &basicModel.Message{
			ID:       utils.GenerateSfIDString(),
			FromID:   currentUserId,
			ToID:     peerInfo.PeerId,
			Content:  param.Message.Content,
			MsgType:  param.Message.MsgType,
			MsgScene: peerInfo.PeerMessageScene,
			FromType: constant.MessageFromType.User,
		}
		if peerInfo.PeerMessageScene == constant.MessageScene.User {
			message.SessionID = utils.Generate1v1SessionID(currentUserId, peerInfo.PeerId)
		} else if peerInfo.PeerMessageScene == constant.MessageScene.Group {
			message.SessionID = peerInfo.PeerId
		} else {
			return fmt.Errorf("param.error")
		}
		_, _ = s.SendMessageToSession(currentUserId, message.SessionID, peerInfo.PeerMessageScene, message)
	}
	return nil
}
