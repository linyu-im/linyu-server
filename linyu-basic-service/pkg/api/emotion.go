package api

import (
	"github.com/gin-gonic/gin"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
)

func init() {
	route.Register("POST", "/basic/v1/emotion/list", EmotionListHandler)
}

func EmotionListHandler(c *gin.Context) {
	list, err := basicService.EmotionService.EmotionList()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
