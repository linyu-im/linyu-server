package api

import (
	"github.com/gin-gonic/gin"
	aiParam "github.com/linyu-im/linyu-server/linyu-ai/pkg/param"
	aiService "github.com/linyu-im/linyu-server/linyu-ai/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/ai/v1/skill/list", SkillListHandler)
}

// SkillListHandler 查询技能列表
func SkillListHandler(c *gin.Context) {
	param := &aiParam.SkillListParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	list, err := aiService.AiSkillService.List(param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
