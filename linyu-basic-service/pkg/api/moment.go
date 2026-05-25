package api

import (
	"github.com/gin-gonic/gin"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/basic/v1/moment/create", MomentCreateHandler)
	route.Register("POST", "/basic/v1/moment/page", MomentPageHandler)
	route.Register("POST", "/basic/v1/moment/delete", MomentDeleteHandler)
	route.Register("POST", "/basic/v1/moment/like/add", MomentLikeAddHandler)
	route.Register("POST", "/basic/v1/moment/like/cancel", MomentLikeCancelHandler)
	route.Register("POST", "/basic/v1/moment/comment/add", MomentCommentAddHandler)
	route.Register("POST", "/basic/v1/moment/comment/del", MomentCommentDelHandler)
}

func MomentCreateHandler(c *gin.Context) {
	param := &basicParam.MomentCreateParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.MomentService.CreateMoment(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

func MomentDeleteHandler(c *gin.Context) {
	param := &basicParam.MomentDeleteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.MomentService.MomentDelete(currentUserId, param.MomentId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

func MomentPageHandler(c *gin.Context) {
	param := &basicParam.MomentPageParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	pages, err := basicService.MomentService.PageMoment(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, pages)
}

func MomentLikeAddHandler(c *gin.Context) {
	param := &basicParam.MomentLikeAddParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	like, err := basicService.MomentService.MomentLikeAdd(currentUserId, param.MomentId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, like)
}

func MomentLikeCancelHandler(c *gin.Context) {
	param := &basicParam.MomentLikeCancelParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.MomentService.MomentLikeCancel(currentUserId, param.MomentId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

func MomentCommentAddHandler(c *gin.Context) {
	param := &basicParam.MomentCommentAddParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	comment, err := basicService.MomentService.MomentCommentAdd(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, comment)
}

func MomentCommentDelHandler(c *gin.Context) {
	param := &basicParam.MomentCommentDelParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.MomentService.MomentCommentDel(currentUserId, param.CommentId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
