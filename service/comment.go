package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

type Comment struct {
}

func NewCommentInstance() *Comment {
	return &Comment{}
}

// Publish 发布评论
func (Comment) Publish(comment *models.CommentDetail) (err error) {
	// 1. 向数据库中增加一条评论
	commentId, err := mysql.CreateComment(comment)
	if err != nil {
		zap.L().Error("mysql CreateComment method err", zap.Error(err))
		return
	}
	// 2. 根据评论类型向redis中创建评论信息
	if comment.Type == mysql.CommentPost {
		if err = redis.CreatePostComment(comment.PostId, int(commentId)); err != nil {
			zap.L().Error("redis CreateFComment method err", zap.Error(err))
		}
	} else { // 如果不是对帖子，那么要么对评论，要么对父评论
		if err = redis.CreateChildComment(comment.FatherId, int(commentId)); err != nil {
			zap.L().Error("redis CreateCComment method err", zap.Error(err))
		}
	}
	return
}

// FavoriteBuild 点赞创建
func (Comment) FavoriteBuild(f *models.Favorite, uid string) error {
	// 对帖子评论需要计算分值，而对人评论不需要计算分值
	if f.Type == mysql.CommentPost {
		if !f.Agree {
			return redis.DeleteFFavorite(f.PostId, uid, f.Id)
		}
		return redis.CreateFFavorite(f.PostId, uid, f.ToAuthorId, f.Id)
	}
	if !f.Agree {
		return redis.DeleteCFavorite(uid, f.Id)
	}
	return redis.CreateCFavorite(uid, f.ToAuthorId, f.Id)
}

// Delete 删除评论 如果是子评论仅仅删除数据库即可
func (Comment) Delete(commentD *models.CommentDelete) (err error) {
	if err = mysql.DeleteComment(commentD.Id); err != nil {
		zap.L().Error("mysql deleteComment method err", zap.Error(err))
		return
	}
	// Type_id 如果删除的是对帖子的评论，则type_id为 post_id,否则为 fCommentId
	// 如果是对帖子的评论还需要删除mysql中子评论 子评论用定时任务进行删除，因为没有帖子就不会因为子评论影响到业务
	if commentD.Type == mysql.CommentPost {
		if err = redis.DeleteFatherComment(commentD.TypeId, commentD.Id); err != nil {
			zap.L().Error("redis deleteComment method err", zap.Error(err))
		}
		return
	}
	// 剩下的就是去删除子评论
	return redis.DeleteChildComment(commentD.TypeId, commentD.Id)
}

type CommentAll struct {
	*models.Comment
	Children []*models.CommentDetail `json:"children,omitempty"`
}

// GetListAll 根据排序获取评论
func (Comment) GetListAll(pid, order string) (commentDatas []*CommentAll, err error) {
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
	if len(fList) <= 0 {
		return
	}
	// 3. 无需对子评论进行排序处理，只需要根据排序过后的父评论去查找子评论即可
	cList, err := redis.GetAllCommentId(fList)
	if err != nil {
		zap.L().Error("redis GetAllCommentId method err", zap.Error(err))
		return
	}
	// 批量取出父评论
	fcomments, err := mysql.GetCommentByIds(fList)
	if err != nil {
		zap.L().Error("mysql GetCommentById method err", zap.Error(err))
		return
	}
	// 批量取出父评论点赞数
	favorites, err := redis.GetFavoriteList(fList)
	if err != nil {
		zap.L().Error("redis GetFavoriteList method err", zap.Error(err))
		return
	}
	commentDatas = make([]*CommentAll, 0, len(fList))
	var cData []*models.CommentDetail
	// 4. 遍历父id，查找mysql中的数据
	for i, _ := range fList {
		// 取出父评论点赞数
		fData := fcomments[i]
		fData.FavoriteNum = favorites[i]
		if len(cList[i]) > 0 {
			// 取出该父评论的子评论列表
			cData, err = mysql.GetCommentList(cList[i])
			if err != nil {
				zap.L().Error("mysql GetCommentList method err", zap.Error(err))
				err = nil
			}
			// 取出该付评论的子评论的点赞数列表
			favoriteList, err := redis.GetFavoriteList(cList[i])
			if err != nil {
				zap.L().Error("redis GetFavoriteList method err", zap.Error(err))
				err = nil
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
