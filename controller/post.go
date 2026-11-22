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
	if err := logic.CreatePost(c.Request.Context(), p); err != nil {
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

// GetPostListHandler 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var (
		page int64
		size int64
		err  error
	)

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}

	//执行业务逻辑：获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() error", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版获取帖子列表（根据时间或分数查询）
// 根据前端传来的参数动态获取帖子列表
// 按创建时间排序 或者 按照帖子分数排序
// 1.获取参数
// 2.去redis查询id列表
// 3.根据id列表去数据库查询帖子详情
func GetPostListHandler2(c *gin.Context) {
	//GET请求参数：/api/v1/posts/2?page=1&size=10&order=time
	// 获取分页参数
	//初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, //默认按照时间排序（同时代码避免硬编码）
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2() error", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//pageStr := c.Query("page")
	//sizeStr := c.Query("size")
	//var (
	//	page int64
	//	size int64
	//	err  error
	//)
	//page, err = strconv.ParseInt(pageStr, 10, 64)
	//if err != nil {
	//	page = 1
	//}
	//size, err = strconv.ParseInt(sizeStr, 10, 64)
	//if err != nil {
	//	size = 10
	//}

	//执行业务逻辑：获取数据
	data, err := logic.GetPostList2(c.Request.Context(), p)
	if err != nil {
		zap.L().Error("logic.GetPostList() error", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}
