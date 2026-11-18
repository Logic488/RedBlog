package mysql

import (
	"RedBlog/models"
	"database/sql"

	"go.uber.org/zap"
)

// InsertPost 插入新帖子
func InsertPost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id, author_id, community_id, title, content) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.AuthorID, p.CommunityID, p.Title, p.Content)
	return
}

// GetPostDetailById 根据id查询帖子详情
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
