package service

import (
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event/eventbus"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct{}

func (s messageService) SendMessageToUser(userId string, param *basicParam.SendMessageToUserParam) error {
	//创建消息
	message := &basicModel.Message{
		ID:        utils.GenerateSfIDString(),
		SessionID: utils.Generate1v1SessionID(userId, param.ToUserId),
		FromID:    userId,
		ToID:      param.ToUserId,
		MsgScene:  constant.MessageScene.User,
		Content:   param.Content,
		MsgType:   constant.MessageType.Text,
	}
	err := basicDao.MessageDao.Create(db.RDB, message)
	if err != nil {
		return err
	}
	//更新对方的聊天会话（增加未读数）
	if param.ToUserId != userId {
		_ = ChatService.SaveOrUpdateIncUnreadNum(param.ToUserId, userId, message)
	}
	//更新自己的会话
	_ = ChatService.SaveOrUpdate(userId, param.ToUserId, message)
	//发送消息
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  []string{param.ToUserId},
		Data: &event.WsData{
			SeqId:   message.ID,
			Type:    constant.WsDataType.Message,
			Content: message,
		},
	})
	return nil
}

func (s messageService) SendMessageToGroup(userId string, param *basicParam.SendMessageToGroupParam) error {
	//创建消息
	message := &basicModel.Message{
		ID:        utils.GenerateSfIDString(),
		SessionID: param.ToGroupId,
		FromID:    userId,
		ToID:      param.ToGroupId,
		MsgScene:  constant.MessageScene.User,
		Content:   param.Content,
		MsgType:   constant.MessageType.Text,
	}
	err := basicDao.MessageDao.Create(db.RDB, message)
	if err != nil {
		return err
	}
	memberIds := GroupService.GetMemberUserIdsByGroupId(param.ToGroupId)
	//更新群成员的会话
	for _, id := range memberIds {
		if id != userId {
			_ = ChatService.SaveOrUpdateIncUnreadNum(id, param.ToGroupId, message)
		}
	}
	//更新自己的会话
	_ = ChatService.SaveOrUpdate(userId, param.ToGroupId, message)
	//发送消息
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  memberIds,
		Data: &event.WsData{
			SeqId:   message.ID,
			Type:    constant.WsDataType.Message,
			Content: message,
		},
	})
	return nil
}

func (s messageService) GetMessageBySessionId(sessionId string, num int) []*basicModel.Message {
	return basicDao.MessageDao.GetLatestMessagesBySessionID(db.RDB, sessionId, num)
}
