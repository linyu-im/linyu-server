package gateway

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/linyu-im/linyu-server/linyu-ai/pkg/api"
	_ "github.com/linyu-im/linyu-server/linyu-application/pkg/api"
	_ "github.com/linyu-im/linyu-server/linyu-auth/pkg/api"
	_ "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/api"
	_ "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/api"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-gateway/pkg/middleware"
	_ "github.com/linyu-im/linyu-server/linyu-im/pkg/api"
	_ "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/api"
)

func Run() {
	engine := gin.New()
	engine.Use(middleware.Cors())
	engine.Use(middleware.Auth())
	apiGroup := engine.Group(config.C.Server.RoutePrefix)
	apiGroup.Use(middleware.ReqLogger())
	// 加载路由 使用 import _ 将相关模块引入
	for _, apiRoute := range route.Routers {
		if apiRoute == nil || apiRoute.Path == "" || apiRoute.Handler == nil {
			continue
		}
		// 添加白名单
		if apiRoute.IsWhite {
			middleware.AddWhiteListPath(config.C.Server.RoutePrefix + apiRoute.Path)
		}
		switch apiRoute.Method {
		case "GET":
			apiGroup.GET(apiRoute.Path, apiRoute.Handler)
		case "POST":
			apiGroup.POST(apiRoute.Path, apiRoute.Handler)
		case "PUT":
			apiGroup.PUT(apiRoute.Path, apiRoute.Handler)
		case "DELETE":
			apiGroup.DELETE(apiRoute.Path, apiRoute.Handler)

		}
	}
	go func() {
		err := engine.Run(fmt.Sprintf(":%d", config.C.Server.Port))
		if err != nil {
			panic("app start error: " + err.Error())
		}
	}()
}
