package service

import (
	"errors"

	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event/eventbus"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	voipParam "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/param"
	voipResult "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/result"
)

var AvCallService = newAvCallService()

func newAvCallService() *avCallService {
	return &avCallService{}
}

type avCallService struct{}

// FriendInvite 发起好友通话邀请
func (s *avCallService) FriendInvite(userId string, param *voipParam.AvCallFriendInviteParam) (*voipResult.AvCallInviteResult, error) {
	if param.UserId == "" || param.UserId == userId {
		return nil, errors.New("param.error")
	}
	if !constant.AvCallType.Validate(param.CallType) {
		return nil, errors.New("param.type-error")
	}
	// 验证双方是否都是好友
	if !basicService.ContactsService.IsFriendBothAnd(userId, param.UserId) {
		return nil, errors.New("basic.contacts.you-no-other-friend")
	}

	sessionId := utils.GenerateSessionID(userId, param.UserId, constant.SceneType.User)
	inviteId := utils.GenerateSfIDString()
	content := &voipResult.AvCallWsContent{
		Action:    constant.AvCallAction.Invite,
		SessionId: sessionId,
		FromId:    userId,
		ToUserIds: []string{param.UserId},
		CallType:  param.CallType,
		SceneType: constant.SceneType.User,
	}
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  []string{param.UserId},
		Data: &event.WsData{
			SeqId:   inviteId,
			Type:    constant.WsDataType.Call,
			Content: content,
		},
	})

	return &voipResult.AvCallInviteResult{
		SessionId: sessionId,
		UserId:    param.UserId,
		CallType:  param.CallType,
		SceneType: constant.SceneType.User,
	}, nil
}

// GroupInvite 发起群聊通话邀请
func (s *avCallService) GroupInvite(userId string, param *voipParam.AvCallGroupInviteParam) (*voipResult.AvCallInviteResult, error) {
	if param.GroupId == "" || len(param.UserIds) == 0 {
		return nil, errors.New("param.error")
	}
	if !constant.AvCallType.Validate(param.CallType) {
		return nil, errors.New("param.type-error")
	}
	if !basicService.GroupService.IsGroupMember(param.GroupId, userId) {
		return nil, errors.New("param.error")
	}

	// 仅邀请传入用户中的群成员（排除自己、去重）
	toUserIds := s.filterGroupMemberUserIds(param.GroupId, userId, param.UserIds)
	if len(toUserIds) == 0 {
		return nil, errors.New("param.error")
	}

	sessionId := utils.GenerateSessionID(userId, param.GroupId, constant.SceneType.Group)
	content := &voipResult.AvCallWsContent{
		Action:    constant.AvCallAction.Invite,
		SessionId: sessionId,
		FromId:    param.GroupId,
		ToUserIds: toUserIds,
		CallType:  param.CallType,
		SceneType: constant.SceneType.Group,
	}
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  toUserIds,
		Data: &event.WsData{
			SeqId:   utils.GenerateSfIDString(),
			Type:    constant.WsDataType.Call,
			Content: content,
		},
	})
	// 成员变更通知：全群成员（含自己）
	s.publishCallChange(userId, sessionId, param.GroupId, param.CallType)

	return &voipResult.AvCallInviteResult{
		SessionId: sessionId,
		GroupId:   param.GroupId,
		CallType:  param.CallType,
		SceneType: constant.SceneType.Group,
	}, nil
}

// UserHangup 单聊通话挂断，通知对端
func (s *avCallService) UserHangup(userId string, param *voipParam.AvCallUserHangupParam) error {
	if param.UserId == "" || param.UserId == userId {
		return errors.New("param.error")
	}

	sessionId := utils.GenerateSessionID(userId, param.UserId, constant.SceneType.User)
	content := &voipResult.AvCallWsContent{
		Action:    constant.AvCallAction.Hangup,
		SessionId: sessionId,
		FromId:    userId,
		SceneType: constant.SceneType.User,
	}
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  []string{param.UserId},
		Data: &event.WsData{
			SeqId:   utils.GenerateSfIDString(),
			Type:    constant.WsDataType.Call,
			Content: content,
		},
	})
	return nil
}

// GroupHangup 群聊通话挂断，通知指定群成员
func (s *avCallService) GroupHangup(userId string, param *voipParam.AvCallGroupHangupParam) error {
	if param.GroupId == "" || len(param.UserIds) == 0 {
		return errors.New("param.error")
	}
	if !basicService.GroupService.IsGroupMember(param.GroupId, userId) {
		return errors.New("param.error")
	}

	toUserIds := s.filterGroupMemberUserIds(param.GroupId, userId, param.UserIds)
	if len(toUserIds) == 0 {
		return errors.New("param.error")
	}

	sessionId := utils.GenerateSessionID(userId, param.GroupId, constant.SceneType.Group)
	content := &voipResult.AvCallWsContent{
		Action:    constant.AvCallAction.Hangup,
		SessionId: sessionId,
		FromId:    userId,
		SceneType: constant.SceneType.Group,
	}
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  toUserIds,
		Data: &event.WsData{
			SeqId:   utils.GenerateSfIDString(),
			Type:    constant.WsDataType.Call,
			Content: content,
		},
	})
	// 成员变更通知：全群成员（含自己）
	s.publishCallChange(userId, sessionId, param.GroupId, "")
	return nil
}

// publishCallChange 向全群成员推送 change（含自己）
func (s *avCallService) publishCallChange(userId, sessionId, groupId, callType string) {
	allUserIds := basicService.GroupService.GetMemberUserIdsByGroupId(groupId)
	if len(allUserIds) == 0 {
		return
	}
	content := &voipResult.AvCallWsContent{
		Action:    constant.AvCallAction.Change,
		SessionId: sessionId,
		FromId:    groupId,
		ToUserIds: allUserIds,
		CallType:  callType,
		SceneType: constant.SceneType.Group,
	}
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  allUserIds,
		Data: &event.WsData{
			SeqId:   utils.GenerateSfIDString(),
			Type:    constant.WsDataType.Call,
			Content: content,
		},
	})
}

// filterGroupMemberUserIds 过滤出群成员（排除自己、去重、非群成员）
func (s *avCallService) filterGroupMemberUserIds(groupId, currentUserId string, userIds []string) []string {
	seen := make(map[string]struct{}, len(userIds))
	result := make([]string, 0, len(userIds))
	for _, id := range userIds {
		if id == "" || id == currentUserId {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if !basicService.GroupService.IsGroupMember(groupId, id) {
			continue
		}
		result = append(result, id)
	}
	return result
}
