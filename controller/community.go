package controller

import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// ----跟社区相关的----

func CommunityHandler(c *gin.Context) {
	// 期望查询到所有的社区(community_id,conmmunity_name)以切片形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList()", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 根据id查询到社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 查询该社区详细内容
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail()", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端暴露给外面
		return
	}
	ResponseSuccess(c, data)
}
