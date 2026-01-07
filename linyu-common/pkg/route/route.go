package route

import "github.com/gin-gonic/gin"

var Routers []*ApiRoute

type ApiRoute struct {
	Path    string
	Method  string
	Handler func(c *gin.Context)
	IsWhite bool
}

func Register(method string, path string, handler func(c *gin.Context), isWhite bool) {
	Routers = append(Routers, &ApiRoute{
		Path:    path,
		Method:  method,
		Handler: handler,
		IsWhite: isWhite,
	})
}
