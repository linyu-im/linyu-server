package response

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/i18n"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Ok(ctx *gin.Context, data ...interface{}) {
	resp := Response{
		Code: 0,
		Msg:  i18n.T(ctx, "param.success", nil),
	}
	if len(data) > 0 {
		resp.Data = data[0]
	}
	ctx.JSON(http.StatusOK, resp)
}

func OkMsg(ctx *gin.Context, msg string, data interface{}) {
	i18nMsg := i18n.T(ctx, msg, nil)
	ctx.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  i18nMsg,
		Data: data,
	})
}

func Fail(ctx *gin.Context, msg string) {
	i18nMsg := i18n.T(ctx, msg, nil)
	ctx.AbortWithStatusJSON(http.StatusOK, Response{
		Code: CodeServerError,
		Msg:  i18nMsg,
	})
}

func FailErrCode(ctx *gin.Context, code int, msg string) {
	i18nMsg := i18n.T(ctx, msg, nil)
	ctx.AbortWithStatusJSON(http.StatusOK, Response{
		Code: code,
		Msg:  i18nMsg,
	})
}

func FailWithData(ctx *gin.Context, code int, msg string, data interface{}) {
	i18nMsg := i18n.T(ctx, msg, nil)
	ctx.AbortWithStatusJSON(http.StatusOK, Response{
		Code: code,
		Msg:  i18nMsg,
		Data: data,
	})
}

func FailWithErrData(ctx *gin.Context, msg string, errData map[string]interface{}) {
	i18nMsg := i18n.T(ctx, msg, errData)
	ctx.AbortWithStatusJSON(http.StatusOK, Response{
		Code: CodeServerError,
		Msg:  i18nMsg,
	})
}
