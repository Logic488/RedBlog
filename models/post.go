package models

import "time"

// Post 创建帖子接口的结构体
type Post struct {
	ID          int64     `json:"id" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// PostDetail 帖子详情接口的结构体
type PostDetail struct {
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"`
	*Post                               //嵌入帖子结构体
	*CommunityDetail `json:"community"` //嵌入社区信息结构体
}
