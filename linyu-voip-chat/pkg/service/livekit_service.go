package service

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"github.com/livekit/protocol/auth"
	"time"
)

var LivekitService = newLiveKitService()

func newLiveKitService() *liveService {
	return &liveService{}
}

type liveService struct{}

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
