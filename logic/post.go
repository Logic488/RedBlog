package logic

import (
	"RedBlog/dao/mysql"
	"RedBlog/dao/redis"
	"RedBlog/models"
	"RedBlog/pkg/snowflake"
	"context"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(ctx context.Context, p *models.Post) (err error) {
	//雪花算法生成post_id
	p.ID = snowflake.GenID()
	//插入post数据
	err = mysql.InsertPost(p)

	err = redis.CreatePost(ctx, p.ID)

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
		zap.L().Error("mysql.GetUserByID() failed", zap.Error(err))
		return
	}
	//根据community_id查询社区(community)信息
	communityID := post.CommunityID
	communityDetail, err := mysql.GetCommunityDetailByID(communityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Error(err))
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

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.PostDetail, err error) {
	//查询数据库，获取Post信息
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return
	}

	data = make([]*models.PostDetail, 0, len(posts))

	for _, post := range posts {
		//根据author_id查询用户user信息
		authorID := post.AuthorID
		user, err := mysql.GetUserByID(authorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Error(err))
			continue
		}
		//根据community_id查询社区(community)信息
		communityID := post.CommunityID
		communityDetail, err := mysql.GetCommunityDetailByID(communityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Error(err))
			continue
		}

		//组装数据
		postDetail := &models.PostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}

		data = append(data, postDetail)
	}
	//返回数据
	return
}

func GetPostList2(ctx context.Context, p *models.ParamPostList) (data []*models.PostDetail, err error) {
	// 2.去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(ctx, p)
	if err != nil {
		return
	}
	//如果没有帖子，返回空
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() return 0 data", zap.Any("param", p))
		return
	}
	// 3.根据id列表去数据库查询帖子详情
	// 返回的数据顺序要按照我给定id的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ctx, ids)
	if err != nil {
		return
	}

	data = make([]*models.PostDetail, 0, len(posts))

	//将帖子的作者以及分区信息查出来填充到帖子中
	for idx, post := range posts {
		//根据author_id查询用户user信息
		authorID := post.AuthorID
		user, err := mysql.GetUserByID(authorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Error(err))
			continue
		}
		//根据community_id查询社区(community)信息
		communityID := post.CommunityID
		communityDetail, err := mysql.GetCommunityDetailByID(communityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Error(err))
			continue
		}

		//组装数据
		postDetail := &models.PostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: communityDetail,
		}

		data = append(data, postDetail)
	}
	//返回数据
	return

}
