package mysql

import (
	"bluebell/models"

	"go.uber.org/zap"
)

// GetCommunityList 获取社区信息列表
func GetCommunityList() (data []*models.Community, err error) {
	qStr := `select community_id,community_name from community order by create_time ASC`
	err = db.Select(&data, qStr)
	if err != nil {
		if err == ErrNoRows {
			zap.L().Warn("getCommunityList is null data")
			err = nil
		}
	}
	return
}

func GetCommunityDetail(id int64) (communityDeatil *models.CommunityDetail, err error) {
	communityDeatil = new(models.CommunityDetail)
	qStr := `select community_id,community_name,introduction,create_time
				from community
				where community_id = ?
			`
	err = db.Get(communityDeatil, qStr, id)
	if err == ErrNoRows {
		zap.L().Warn("GetCommunityDetail is null data")
		err = ErrorInvalidParam
	}
	return
}
