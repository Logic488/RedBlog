package controller

//与社区相关

import (
	"RedBlog/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 查询得到所有社区
func CommunityHandler(c *gin.Context) {
	//业务操作：查询得到所有社区，(community_id,community_name),以列表形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() error", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错暴露给外面
		return
	}
	//返回得到的社区列表
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 查询社区详情
func CommunityDetailHandler(c *gin.Context) {
	//获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//业务操作：查询得到社区详情信息
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() error", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错暴露给外面
		return
	}
	//返回得到的社区详情信息
	ResponseSuccess(c, data)
}
