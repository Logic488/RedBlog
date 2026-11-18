package logic

import (
	"RedBlog/dao/mysql"
	"RedBlog/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	//查询所有社区
	data, err = mysql.GetCommunityList()
	return
}

func GetCommunityDetail(id int64) (data *models.CommunityDetail, err error) {
	data, err = mysql.GetCommunityDetailByID(id)
	return
}
