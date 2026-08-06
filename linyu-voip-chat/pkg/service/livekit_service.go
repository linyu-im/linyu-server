package service

import (
	"context"
	"strings"
	"sync"
	"time"

	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	voipResult "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/result"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

var LivekitService = newLiveKitService()

func newLiveKitService() *liveService {
	return &liveService{}
}

type liveService struct{}

var (
	roomServiceClientMu sync.Mutex
	roomServiceClient   *lksdk.RoomServiceClient
)

func (s *liveService) LivekitToken(room, userId string) (string, error) {
	at := auth.NewAccessToken(config.C.Livekit.Key, config.C.Livekit.Secret)
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         room,
		CanPublish:   utils.BoolPtr(true),
		CanSubscribe: utils.BoolPtr(true),
	}

	at.SetVideoGrant(grant).
		SetIdentity(userId).
		SetValidFor(time.Hour * 5)

	token, err := at.ToJWT()
	if err != nil {
		return "", err
	}
	return token, nil
}

func roomServiceURL(host string) string {
	host = strings.TrimSpace(host)
	switch {
	case strings.HasPrefix(host, "https://"), strings.HasPrefix(host, "http://"):
		return host
	case strings.HasPrefix(host, "wss://"):
		return "https://" + strings.TrimPrefix(host, "wss://")
	case strings.HasPrefix(host, "ws://"):
		return "http://" + strings.TrimPrefix(host, "ws://")
	default:
		return "http://" + host
	}
}

func (s *liveService) getRoomServiceClient() *lksdk.RoomServiceClient {
	roomServiceClientMu.Lock()
	defer roomServiceClientMu.Unlock()
	if roomServiceClient != nil {
		return roomServiceClient
	}
	roomServiceClient = lksdk.NewRoomServiceClient(
		roomServiceURL(config.C.Livekit.Host),
		config.C.Livekit.Key,
		config.C.Livekit.Secret,
	)
	return roomServiceClient
}

func (s *liveService) resetRoomServiceClient() *lksdk.RoomServiceClient {
	roomServiceClientMu.Lock()
	defer roomServiceClientMu.Unlock()
	roomServiceClient = lksdk.NewRoomServiceClient(
		roomServiceURL(config.C.Livekit.Host),
		config.C.Livekit.Key,
		config.C.Livekit.Secret,
	)
	return roomServiceClient
}

// ListRoomUsers 查询房间内在线用户列表（room 名为 sessionId，单聊/群聊均可）
func (s *liveService) ListRoomUsers(sessionId string) ([]*voipResult.LivekitRoomUser, error) {
	participants, err := s.listParticipants(sessionId)
	if err != nil {
		// 失败时清空并重建 client 后重试一次
		s.resetRoomServiceClient()
		participants, err = s.listParticipants(sessionId)
		if err != nil {
			return nil, err
		}
	}

	result := make([]*voipResult.LivekitRoomUser, 0, len(participants))
	for _, p := range participants {
		if p == nil || p.Identity == "" {
			continue
		}
		item := &voipResult.LivekitRoomUser{
			UserId: p.Identity,
			State:  int32(p.State),
		}
		if user := basicService.UserService.GetUserById(p.Identity); user != nil {
			item.Username = user.Username
			item.Avatar = user.Avatar
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *liveService) listParticipants(room string) ([]*livekit.ParticipantInfo, error) {
	client := s.getRoomServiceClient()
	resp, err := client.ListParticipants(context.Background(), &livekit.ListParticipantsRequest{
		Room: room,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Participants) == 0 {
		return []*livekit.ParticipantInfo{}, nil
	}
	return resp.Participants, nil
}
