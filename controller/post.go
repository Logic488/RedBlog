package controller

import (
	"RedBlog/logic"
	"RedBlog/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	//接收参数,参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Create Post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c取到当前发送请求的用户ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() fail", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//执行业务逻辑：创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("Create Post failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}

// PostDetailHandler 获取帖子详情
func PostDetailHandler(c *gin.Context) {
	//从请求路径得到post_id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//执行业务逻辑：获取帖子详情
	data, err := logic.GetPostDetail(id)
	if err != nil {
		zap.L().Error("logic.GetPostDetail() error", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)

}
