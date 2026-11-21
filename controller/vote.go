package controller

import (
	"RedBlog/logic"
	"RedBlog/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostVoteController 帖子投票功能
func PostVoteController(c *gin.Context) {
	//参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译并去除掉错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	//拿到当前请求的用户的id
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//处理业务逻辑：投票
	if err := logic.VoteForPost(c.Request.Context(), userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err), zap.Int64("userID", userID))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}
