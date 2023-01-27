package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/pkg/e"
	silr "bluebell/serializer"

	"go.uber.org/zap"
)

type PostVote struct {
	PostId    string `form:"post_id" bidding:"required"`
	Direction int8   `form:"direction" bidding:"oneof=1 0 -1"` // 规定 1为赞成，0为取消投票，-1为反对
}

// Build 投票构建
func (v PostVote) Build(uid string) (silr.Response, error) {
	code := e.CodeSUCCESS
	// 1. 判断帖子状态
	status, err := mysql.GetPostStatus(v.PostId)
	if err != nil {
		code = e.CodeServerBusy
		zap.L().Error(code.Msg(),
			zap.String("postId", v.PostId),
			zap.String("userId", uid),
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	// 如果帖子已经被删除直接返回参数错误
	if status == mysql.PostDelete {
		code = e.CodeInvalidParams
		zap.L().Error(code.Msg(),
			zap.String("postId", v.PostId),
			zap.String("userId", uid),
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, mysql.ErrNoRows
	} else if status == mysql.PostExpired {
		code = e.CodePostVoteExpired
		zap.L().Error(code.Msg(),
			zap.String("postId", v.PostId),
			zap.String("userId", uid),
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, redis.ErrVoteTimeExpired
	}
	// 2. 判断帖子投票情况
	diff := redis.PostVoteDirect(v.PostId, uid, float64(v.Direction))
	// 3. 更改帖子分数
	if err = redis.ChangeVoteInfo(v.PostId, uid, diff, float64(v.Direction)); err != nil {
		code = e.CodeServerBusy
		zap.L().Error("ChangePostInfo method pipe exec is failed",
			zap.String("postId", v.PostId),
			zap.String("userId", uid),
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	return silr.Response{}, nil
}
