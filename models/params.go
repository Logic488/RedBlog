package models

//用来定义请求参数的结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 用户注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// ParamLogIn 用户登录请求参数
type ParamLogIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`              // 贴子id
	Direction int    `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票(1)还是反对票(-1)还是取消投票(0)
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	Page  int64
	Size  int64
	Order string
}
