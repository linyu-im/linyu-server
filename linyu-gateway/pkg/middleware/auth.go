package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/jwt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"golang.org/x/text/language"
	"strings"
)

func AddWhiteListPaths(paths []string) {
	for _, path := range paths {
		AddWhiteListPath(path)
	}
}

func AddWhiteListPath(Path string) {
	whiteListPaths[Path] = true
}

// 白名单接口
var whiteListPaths = map[string]bool{}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//设置语言
		SetLanguage(ctx)
		ctx.Set("device", DetectDeviceType(ctx.Request.UserAgent()))
		ctx.Set("ip", ctx.ClientIP())

		path := ctx.FullPath()
		if whiteListPaths[path] {
			ctx.Next()
			return
		}

		// 获取 token
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			tokenString, _ = ctx.GetQuery("Authorization")
		}
		if tokenString == "" {
			response.FailErrCode(ctx, response.CodeUnauthorized, "Unauthorized")
			return
		}

		// 解析用户信息
		userInfo, err := jwt.ParseJwtToken(tokenString)
		if err != nil {
			response.FailErrCode(ctx, response.CodeUnauthorized, "Identity information error")
			return
		}

		ctx.Set("userID", userInfo.UserID)
		ctx.Set("roleID", userInfo.RoleID)

		ctx.Next()
	}
}

func SetLanguage(ctx *gin.Context) {
	acceptLang := ctx.GetHeader("Accept-Language")
	if acceptLang == "" {
		acceptLang = "zh"
	}
	tags, _, err := language.ParseAcceptLanguage(acceptLang)
	lang := "zh"
	if err == nil && len(tags) > 0 {
		base, _ := tags[0].Base()
		lang = strings.ToLower(base.String())
	}
	ctx.Set("lang", lang)
}

func DetectDeviceType(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "mobile"), strings.Contains(ua, "android"), strings.Contains(ua, "iphone"):
		return constant.Device.Mobile
	case strings.Contains(ua, "electron"), strings.Contains(ua, "tauri"):
		return constant.Device.Desktop
	case strings.Contains(ua, "mozilla"), strings.Contains(ua, "chrome"), strings.Contains(ua, "safari"):
		return constant.Device.Web
	default:
		return constant.Device.Unknown
	}
}
