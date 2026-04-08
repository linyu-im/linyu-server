package service

import (
	"errors"
	"github.com/linyu-im/linyu-server/linyu-auth/pkg/oauth2"
	"github.com/linyu-im/linyu-server/linyu-auth/pkg/result"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
)

var Oauth2Service = newOauth2Service()

type oauth2Service struct {
	oauth2s map[string]oauth2.Oauth2
}

func newOauth2Service() *oauth2Service {
	return &oauth2Service{}
}

func (s *oauth2Service) GetOauth2Service(oauth2Type string) oauth2.Oauth2 {
	if s.oauth2s == nil {
		s.oauth2s = map[string]oauth2.Oauth2{
			"gitee": oauth2.NewGiteeOauth2(),
		}
	}
	return s.oauth2s[oauth2Type]
}

func (s *oauth2Service) GetAuthURL(oauth2Type string) (*result.Oauth2UrlResult, error) {
	auth := s.GetOauth2Service(oauth2Type)
	if auth == nil {
		return nil, errors.New("param.error")
	}
	authURL, redirectUrl := auth.GetAuthURL()
	return &result.Oauth2UrlResult{
		AuthUrl:     authURL,
		RedirectUrl: redirectUrl,
	}, nil
}

func (s *oauth2Service) GetUserInfo(t string, code string) (*basicModel.User, error) {
	auth := s.GetOauth2Service(t)
	if auth == nil {
		return nil, errors.New("param.error")
	}
	return auth.GetUserInfo(code)
}
