package models

// 用户注册结构体
type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

// 用户详情的结构体
type UserDetail struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Gender   int32  `db:"gender"`
}
