package service

import (
	"bluebell/dao/redis"
	"bluebell/pkg/e"
	silr "bluebell/serializer"

	"go.uber.org/zap"
)

type PostVoteService struct {
	PostId    int64 `json:"post_id,string" form:"post_id" bidding:"required"`
	Direction int8  `json:"direction" form:"direction" bidding:"oneof=1 0 -1"` // 规定 1为赞成，0为取消投票，-1为反对
}

// VoteBuild 投票构建
func (v PostVoteService) VoteBuild(uid int64) (silr.Response, error) {
	code := e.CodeSUCCESS
	// 1. 判断帖子投票时间是否过期
	if err := redis.CheckVoteTime(v.PostId); err != nil {
		code = e.CodePostVoteExpired
		zap.L().Error(code.Msg(),
			zap.Int64("postId", v.PostId),
			zap.Int64("userId", uid),
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	// 2. 判断帖子投票情况
	diff := redis.PostVoteDirect(v.PostId, uid, float64(v.Direction))
	// 3. 更改帖子分数
	if err := redis.ChangePostInfo(v.PostId, uid, diff, float64(v.Direction)); err != nil {
		code = e.CodeServerBusy
		zap.L().Error("ChangePostInfo method pipe exec is failed",
			zap.Int64("postId", v.PostId),
			zap.Int64("userId", uid),
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	return silr.Response{}, nil
}
