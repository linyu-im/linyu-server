package route

import "github.com/gin-gonic/gin"

var Routers []*ApiRoute

type ApiRoute struct {
	Path    string
	Method  string
	Handler func(c *gin.Context)
	IsWhite bool
}

// registerRoute
func registerRoute(method string, path string, handler func(c *gin.Context), isWhite bool) {
	Routers = append(Routers, &ApiRoute{
		Path:    path,
		Method:  method,
		Handler: handler,
		IsWhite: isWhite, // 差异化参数
	})
}

// RegisterWhite 注册白名单路由
func RegisterWhite(method string, path string, handler func(c *gin.Context)) {
	registerRoute(method, path, handler, true)
}

// Register 注册普通路由
func Register(method string, path string, handler func(c *gin.Context)) {
	registerRoute(method, path, handler, false)
}
