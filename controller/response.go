package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//{
//	"code":10000, //程序中的错误码
//	"msg":xx, //显示信息
//	"data":{}, //数据
//}

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"` // omitempty:当字段为null，就不会在json序列化的时候展示这个字段
}

func ResponseError(c *gin.Context, code ResCode) {
	rd := ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
