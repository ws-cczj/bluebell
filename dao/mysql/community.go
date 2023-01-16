package mysql

import (
	"bluebell/models"

	"go.uber.org/zap"
)

const (
	CommunityCheck = iota
	CommunityPublish
	CommunityDelete
)

// CreateCommunity 创建社区
func CreateCommunity(cDetail *models.CommunityDetail) (int64, error) {
	iStr := `insert into 
    			community(author_id,author_name,community_name,introduction)
				values(?,?,?,?) `
	res, err := db.Exec(iStr,
		cDetail.AuthorId,
		cDetail.AuthorName,
		cDetail.Name,
		cDetail.Introduction)
	if err != nil {
		zap.L().Error("Create Community method err", zap.Error(err))
		return -1, err
	}
	return res.LastInsertId()
}

// GetCommunityList 获取社区列表
func GetCommunityList(pidNums int64) (data []*models.Community, err error) {
	data = make([]*models.Community, 0, pidNums)
	qStr := `select id,author_id,author_name,community_name
				from community
				order by create_time DESC`
	if err = db.Select(&data, qStr); err != nil {
		if err == ErrNoRows {
			zap.L().Warn("getCommunityList is null data")
			err = nil
		}
	}
	return
}

// GetCommunitys 获取社区数目
func GetCommunitys() (cidNum int64, err error) {
	qStr := `select COUNT(*)
				from community`
	err = db.Get(&cidNum, qStr)
	return
}

// GetCommunityDetail 获取社区的详细信息
func GetCommunityDetail(cid int64) (communityDeatil *models.CommunityDetail, err error) {
	communityDeatil = new(models.CommunityDetail)
	qStr := `select id,author_id,author_name,community_name,introduction,status,create_time,update_time
				from community
				where id = ?
			`
	err = db.Get(communityDeatil, qStr, cid)
	return
}
