package logic

import (
	"RedBlog/dao/mysql"
	"RedBlog/models"
	"RedBlog/pkg/jwt"
	"RedBlog/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//1.判断用户存不存在
	err = mysql.QueryUserByUsername(p.Username)
	if err != nil {
		//数据库查询出错
		return err
	}
	//2.生成UID
	userID := snowflake.GenID()
	//构造一个user实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存进数据库
	err = mysql.InsertUser(user)
	return
}

func LogIn(p *models.ParamLogIn) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，就能拿到user.UserID
	if err := mysql.CheckUserLogin(user); err != nil {
		return "", err
	}
	//生成JWT
	return jwt.GenToken(user.UserID, user.Username)

}
