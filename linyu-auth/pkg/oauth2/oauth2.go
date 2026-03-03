package oauth2

import basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"

type Oauth2 interface {
	GetAuthURL() string
	GetUserInfo(code string) (*basicModel.User, error)
}
