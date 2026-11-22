package mysql

import (
	"RedBlog/models"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// InsertPost 插入新帖子
func InsertPost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id, author_id, community_id, title, content) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.AuthorID, p.CommunityID, p.Title, p.Content)
	return
}

// GetPostDetailById 根据id查询单篇帖子详情
func GetPostDetailById(id int64) (p *models.Post, err error) {
	p = new(models.Post)
	sqlStr := `select post_id, author_id, title, content, community_id, status, create_time from post where post_id = ?`
	if err = db.Get(p, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("can not find the post by id")
			err = nil
		}
	}
	return
}

// GetPostList 分页查询帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, author_id, title, content, community_id, status, create_time 
		from post 
	    ORDER BY create_time DESC
	    limit ?,?
	    `
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据id列表查询帖子
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
from post
where post_id in (?)
order by FIND_IN_SET(post_id, ?)`

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)

	err = db.Select(&postList, query, args...)
	return
}
