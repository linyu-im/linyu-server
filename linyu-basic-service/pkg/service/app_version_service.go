package service

import (
	"errors"

	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicResult "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/result"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var AppVersionService = newAppVersionService()

func newAppVersionService() *appVersionService {
	return &appVersionService{}
}

type appVersionService struct{}

func (s *appVersionService) Check(param *basicParam.AppVersionCheckParam) (*basicResult.AppVersionCheckResult, error) {
	if !constant.AppPlatform.Validate(param.Platform) {
		return nil, errors.New("param.type-error")
	}

	result := &basicResult.AppVersionCheckResult{
		Platform: param.Platform,
	}

	cfg, err := basicDao.AppVersionDao.GetByPlatform(db.RDB, param.Platform)
	if err != nil {
		return nil, err
	}
	if cfg == nil || !cfg.Enabled {
		return result, nil
	}

	result.LatestVersion = cfg.LatestVersion
	result.LatestVersionCode = cfg.LatestVersionCode
	result.MinSupportVersion = cfg.MinSupportVersion
	result.MinSupportVersionCode = cfg.MinSupportVersionCode
	result.DownloadUrl = cfg.DownloadUrl
	result.UpdateDesc = cfg.UpdateDesc
	return result, nil
}

// EnsureMinSupport 登录前校验：低于最小支持版本则拒绝；无配置/未启用则放行
func (s *appVersionService) EnsureMinSupport(platform string, versionCode int) error {
	if !constant.AppPlatform.Validate(platform) {
		return errors.New("param.type-error")
	}
	cfg, err := basicDao.AppVersionDao.GetByPlatform(db.RDB, platform)
	if err != nil {
		return err
	}
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if versionCode < cfg.MinSupportVersionCode {
		return errors.New("auth.version-too-low")
	}
	return nil
}
