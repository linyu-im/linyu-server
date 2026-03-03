package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strconv"
)

type GiteeOauth2 struct {
	oauth *oauth2.Config
}

type GiteeUserResult struct {
	Id        *int   `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

func NewGiteeOauth2() Oauth2 {
	return &GiteeOauth2{
		oauth: &oauth2.Config{
			ClientID:     config.C.Auth.Gitee.ClientID,
			ClientSecret: config.C.Auth.Gitee.ClientSecret,
			RedirectURL:  config.C.Auth.Gitee.RedirectURL,
			Scopes:       []string{"user_info"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://gitee.com/oauth/authorize",
				TokenURL: "https://gitee.com/oauth/token",
			},
		},
	}
}

func (o *GiteeOauth2) GetAuthURL() string {
	return o.oauth.AuthCodeURL(config.C.Auth.Gitee.RedirectURL)
}

func (o *GiteeOauth2) GetUserInfo(code string) (*basicModel.User, error) {
	token, err := o.oauth.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	client := o.oauth.Client(context.Background(), token)
	resp, err := client.Get("https://gitee.com/api/v5/user")
	if err != nil {
		return nil, fmt.Errorf("get user info failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitee api returned status: %d", resp.StatusCode)
	}

	giteeUserInfo := &GiteeUserResult{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, giteeUserInfo); err != nil {
		return nil, err
	}
	if giteeUserInfo.Id == nil {
		return nil, fmt.Errorf("gitee user info is empty")
	}
	//验证用户是否存在
	userInfo := basicService.UserService.GetUserByGitee(strconv.Itoa(*giteeUserInfo.Id))
	//创建用户
	if userInfo == nil {
		userInfo, err = basicService.UserService.CreateUserByKV("gitee", *giteeUserInfo.Id,
			"username", giteeUserInfo.Name,
			"avatar", giteeUserInfo.AvatarUrl)
		if err != nil {
			return nil, err
		}
	}
	return userInfo, nil
}
