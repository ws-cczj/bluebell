package service

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func CommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

func CommunityDetailById(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetail(id)
}
