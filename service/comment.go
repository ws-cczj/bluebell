package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

type CommentAll struct {
	*models.Comment
	Children []*models.CommentDetail `json:"children,omitempty"`
}

// PublishComment 发布评论
func PublishComment(comment *models.CommentDetail) (err error) {
	// 1. 向数据库中增加一条评论
	commentId, err := mysql.CreateComment(comment)
	if err != nil {
		zap.L().Error("mysql CreateComment method err", zap.Error(err))
		return
	}
	// 2. 根据评论类型向redis中创建评论信息
	if comment.Type == mysql.CommentPost {
		if err = redis.CreatePostComment(comment.PostId, commentId); err != nil {
			zap.L().Error("redis CreateFComment method err", zap.Error(err))
		}
	} else if comment.Type == mysql.CommentPeople {
		if err = redis.CreateChildComment(comment.FatherId, commentId); err != nil {
			zap.L().Error("redis CreateCComment method err", zap.Error(err))
		}
	}
	return
}

// FavoriteBuild 点赞创建
func FavoriteBuild(f *models.Favorite, uid int64) error {
	if !f.Agree {
		return redis.DeleteFavorite(f.PostId, uid, f.Id)
	}
	return redis.CreateFavorite(f.PostId, uid, f.ToAuthorId, f.Id)
}

// DeleteComment 删除评论 如果是子评论仅仅删除数据库即可
func DeleteComment(commentD *models.CommentDelete) (err error) {
	if err = mysql.DeleteComment(commentD.Id); err != nil {
		zap.L().Error("mysql deleteComment method err", zap.Error(err))
		return
	}
	// 如果是对帖子的评论还需要删除redis 子评论不用删除, 否则开销太大
	if commentD.Type != mysql.CommentPeople {
		if err = redis.DeleteComment(commentD.PostId, commentD.Id); err != nil {
			zap.L().Error("redis deleteComment method err", zap.Error(err))
		}
	}
	return
}

// GetCommentList 根据排序获取评论
func GetCommentList(pid int64, order string) (commentDatas []*CommentAll, err error) {
	// 1. 判断排序类别
	key := redis.KeyCommentScoreZSet
	if order == OrderByTime {
		key = redis.KeyCommentTimeZSet
	}
	// 2. 根据排序去redis中查找该帖子的所有父评论
	fList, err := redis.GetFatherCommentId(pid, key)
	if err != nil {
		zap.L().Error("redis GetFatherCommentId method err", zap.Error(err))
		return
	}
	// 3. 无需对子评论进行排序处理，只需要根据排序过后的父评论去查找子评论即可
	cList, err := redis.GetAllCommentId(fList)
	if err != nil {
		zap.L().Error("redis GetAllCommentId method err", zap.Error(err))
		return
	}
	commentDatas = make([]*CommentAll, 0, len(fList))
	var cData []*models.CommentDetail
	var favoriteList []int64
	var fData *models.Comment
	// 4. 遍历父id，查找mysql中的数据
	for i, fid := range fList {
		// 取出父评论
		fData, err = mysql.GetCommentById(fid)
		if err != nil {
			zap.L().Error("mysql GetCommentById method err", zap.Error(err))
			continue
		}
		// 取出父评论点赞数
		fData.FavoriteNum = redis.GetFavorites(fData.Id)
		if len(cList[i]) > 0 {
			// 取出该父评论的子评论列表
			cData, err = mysql.GetCommentList(cList[i])
			if err != nil {
				zap.L().Error("mysql GetCommentList method err", zap.Error(err))
				err = nil
				continue
			}
			// 取出该付评论的子评论的点赞数列表
			favoriteList, err = redis.GetFavoriteList(cList[i])
			if err != nil {
				zap.L().Error("redis GetFavoriteList method err", zap.Error(err))
				err = nil
				continue
			}
			for j, data := range cData {
				data.Comment.FavoriteNum = favoriteList[j]
			}
		}
		commentData := &CommentAll{
			Comment:  fData,
			Children: cData,
		}
		commentDatas = append(commentDatas, commentData)
	}
	return
}
