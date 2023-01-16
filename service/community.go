package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

// CommunityCreate 创建社区
func CommunityCreate(cDetail *models.CommunityDetail) error {
	cDetail.Status = mysql.CommunityPublish
	cNums, err := mysql.CreateCommunity(cDetail)
	if err != nil {
		zap.L().Error("mysql CreateCommunity method err",
			zap.Error(err))
		return err
	}
	// 设置社区数目缓存和用户社区信息
	if err = redis.SetUserCommunity(cNums, cDetail.AuthorId, cDetail.ID); err != nil {
		zap.L().Error("redis SetUserCommunity method err",
			zap.Error(err))
	}
	return err
}

// CommunityList 获取社区列表
func CommunityList() ([]*models.Community, error) {
	// 如果缓存有效，直接通过缓存数量进行查询
	if cidNum, err := redis.GetCommunitys(); err == nil {
		Num, _ := strconv.Atoi(cidNum)
		return mysql.GetCommunityList(int64(Num))
	}
	// 如果缓存无效或者方法错误，就通过mysql查询
	cidNum, err := mysql.GetCommunitys()
	if err != nil {
		zap.L().Error("mysql GetCommunitys method err", zap.Error(err))
		return nil, err
	}
	// 更新缓存，错误也不影响
	_ = redis.SetCommunityNums(cidNum)
	return mysql.GetCommunityList(cidNum)
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
