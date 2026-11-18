package mysql

import (
	"RedBlog/models"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorUserExist    = errors.New("用户已存在")
	ErrorUserNotExist = errors.New("用户不存在")
	ErrorPassword     = errors.New("密码错误")
)

// QueryUserByUsername 注册时查重：用户名已存在则返回错误
func QueryUserByUsername(username string) error {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) error {
	// 加密失败（比如密码超过 72 字节）必须中断，不能带着空密码入库
	pwd, err := encryptPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = pwd
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return err
}

//// CheckUserLogin 检查用户登录信息
//func CheckUserLogin(user *models.User) (err error) {
//	oPassword := user.Password //用户登录的密码
//	sqlStr := `select user_id, username, password from user where username=?`
//	err = db.Get(user, sqlStr, user.Username)
//	if err != nil {
//		return
//	}
//	//判断密码是否正确
//	password := encryptPassword(oPassword)
//	if password != user.Password {
//		return errors.New("密码错误")
//	}
//	return
//}

// CheckUserLogin 检查用户登录信息
func CheckUserLogin(user *models.User) (err error) {
	oPassword := user.Password // 必须先存下来：下面的 db.Get 会把 user.Password 覆盖成库里的哈希
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	// 把库里的哈希和用户输入的明文交给 bcrypt 比对，不要自己重新加密
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oPassword))
	if err != nil {
		return ErrorPassword
	}
	return nil
}

func encryptPassword(oPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(oPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
