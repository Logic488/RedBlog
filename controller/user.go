package controller

//与用户相关

import (
	"RedBlog/dao/mysql"
	"RedBlog/logic"
	"RedBlog/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	//var p models.ParamSignUp
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}

		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		//c.JSON(http.StatusBadRequest, gin.H{
		//	"msg": errs.Translate(trans), //翻译错误
		//})
		return
	}
	//手动对请求参数进行详细的业务规则校验（ShouldBind只能进行简单的格式校验）
	if p.Password != p.RePassword {
		zap.L().Error("SignUp with invalid param", zap.Any("param", p))
		ResponseError(c, CodeInvalidPassword)
		//c.JSON(http.StatusBadRequest, gin.H{
		//	"msg": "两次密码不一致",
		//})
		return
	}
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("SignUp with invalid param", zap.Any("param", p))
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//}

	fmt.Println(p)
	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		//把真正的失败原因记进日志，否则响应里只有"注册失败"，排查时无从下手
		zap.L().Error("logic.SignUp failed", zap.Error(err), zap.String("username", p.Username))
		//c.JSON(http.StatusBadRequest, gin.H{
		//	"msg": "注册失败",
		//})
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success",
	//})
}

// LogInHandler 处理登录请求的函数
func LogInHandler(c *gin.Context) {
	//获取参数和校验参数
	p := new(models.ParamLogIn)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			//c.JSON(http.StatusBadRequest, gin.H{
			//	"msg": err.Error(),
			//})
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		//c.JSON(http.StatusBadRequest, gin.H{
		//	"msg": errs.Translate(trans), //翻译错误
		//})
		return
	}
	//业务处理
	token, err := logic.LogIn(p)
	if err != nil {
		zap.L().Error("logic.LogIn failed", zap.Error(err), zap.String("username", p.Username))
		//对外统一提示，不区分"用户不存在"和"密码错误"，避免被用来枚举用户名
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		//c.JSON(http.StatusUnauthorized, gin.H{
		//	"msg": "用户名或密码错误",
		//})
		return
	}
	//返回响应
	ResponseSuccess(c, token)
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success",
	//})
}
