package logic

import (
	"RedBlog/dao/mysql"
	"RedBlog/models"
	"RedBlog/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	//雪花算法生成post_id
	p.ID = snowflake.GenID()
	//插入post数据
	err = mysql.InsertPost(p)
	//返回响应
	return
}

// GetPostDetail 获取帖子详情
func GetPostDetail(id int64) (data *models.PostDetail, err error) {
	//根据id查询帖子详情
	post, err := mysql.GetPostDetailById(id)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById() failed", zap.Error(err))
		return
	}
	//根据author_id查询用户user信息
	authorID := post.AuthorID
	user, err := mysql.GetUserByID(authorID)
	if err != nil {
		zap.L().Error("mysql.GetUsernameByUserID() failed", zap.Error(err))
		return
	}
	//根据community_id查询社区(community)信息
	communityID := post.CommunityID
	communityDetail, err := mysql.GetCommunityDetailByID(communityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByCommunityID() failed", zap.Error(err))
		return
	}

	data = &models.PostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}

	//authorName := user.Username
	//data = new(models.PostDetail)
	//data.Post = post
	//data.AuthorName = authorName
	//data.CommunityDetail = communityDetail

	return
}
