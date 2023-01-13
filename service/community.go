package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

// CommunityCreate 创建社区
func CommunityCreate(cDetail *models.CommunityDetail) (err error) {
	cDetail.Status = mysql.CommunityPublish
	if err = mysql.CreateCommunity(cDetail); err != nil {
		zap.L().Error("mysql CreateCommunity method err",
			zap.Error(err))
		return
	}
	if err = redis.SetUserCommunity(cDetail.AuthorId, cDetail.ID); err != nil {
		zap.L().Error("redis SetUserCommunity method err",
			zap.Error(err))
	}
	return
}

// CommunityList 获取所有社区
func CommunityList(uid int64) ([]*models.Community, error) {
	if uid != 0 {
		pidNums, err := redis.GetUserCommunitys(uid)
		if err != nil {
			zap.L().Error("redis getUserCommunitys method err", zap.Error(err))
			return nil, err
		}
		return mysql.GetCommunityList(uid, pidNums)
	}
	cidNum, err := mysql.GetCommunitys()
	if err != nil {
		zap.L().Error("mysql GetCommunitys method err", zap.Error(err))
		return nil, err
	}
	return mysql.GetCommunityList(uid, cidNum)
}

// CommunityDetailById 指定获取某个社区详细信息
func CommunityDetailById(cid int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetail(cid)
}

// CommunityPostListInOrder 根据顺序获取社区的帖子列表
func CommunityPostListInOrder(page, size, cid int64, order string) (postList []*PostService, err error) {
	key := redis.KeyPostTimeZSet
	if order == OrderByScore {
		key = redis.KeyPostScoreZSet
	}
	ids, err := redis.GetCommunityPostIds(page, size, cid, key)
	if err != nil {
		zap.L().Error("redis GetPostList method is err",
			zap.Int64("page", page),
			zap.Int64("size", size),
			zap.String("order", order),
			zap.Error(err))
		return
	}
	return getPostListByIds(ids)
}
