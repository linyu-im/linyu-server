package service

import (
	"github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	"github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
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

func (s messageService) SendMessage(userId string, param *basicParam.SendMessageParam) error {
	message := &model.Message{
		ID:      utils.GenerateSfIDString(),
		FromID:  userId,
		ToID:    param.ToUserId,
		Source:  constant.MessageSource.User,
		Content: param.Content,
		Status:  constant.MessageStatus.Unread,
		Type:    constant.MessageType.Text,
	}
	err := dao.MessageDao.Create(db.MysqlDB, message)
	if err != nil {
		return err
	}
	_ = eventbus.DefaultEventBus.Publish(event.MessageEvent{
		FromUserId: userId,
		ToUserId:   param.ToUserId,
		Content:    param.Content,
	})
	return nil
}
