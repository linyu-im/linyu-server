package api

import (
	"github.com/gin-gonic/gin"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
)

func init() {
	route.Register("POST", "/basic/v1/sticker/default/list", DefaultStickerListHandler)
	route.Register("POST", "/basic/v1/sticker/favorite/list", FavoriteStickerListHandler)
	route.Register("POST", "/basic/v1/sticker/user/pack/list", StickerUserPackListHandler)
}

// DefaultStickerListHandler 查询默认表情
func DefaultStickerListHandler(c *gin.Context) {
	list, err := basicService.StickerService.DefaultStickerList()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// FavoriteStickerListHandler 查询当前用户收藏的表情
func FavoriteStickerListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := basicService.StickerService.FavoriteStickerList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// StickerUserPackListHandler 查询用户表情分组内容
func StickerUserPackListHandler(c *gin.Context) {
	// TODO: 根据当前用户拥有的表情分组查询，暂用全部分组数据
	list, err := basicService.StickerService.StickerPackList()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
