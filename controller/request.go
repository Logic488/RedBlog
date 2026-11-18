package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const ContexUserIDKey = "UserID"

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUser 获取当前登录的用户ID
func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContexUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return

}
