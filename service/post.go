package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/e"
	"bluebell/pkg/snowflake"
	silr "bluebell/serializer"

	"go.uber.org/zap"
)

const (
	OrderByTime  = "time"
	OrderByScore = "score"
	NULLData     = 0
)

type Publish struct {
	CommunityId int64  `json:"community_id" form:"community_id" bidding:"required"`
	Title       string `json:"title" form:"title" bidding:"required"`
	Content     string `json:"content" form:"content" bidding:"required"`
}

type PostService struct {
	*models.Post
	*models.CommunityDetail `json:"community_detail"`
}

// PublishPost 发布帖子
func (p Publish) PublishPost(uid int64, uname string) (err error) {
	post := &models.Post{
		PostId:      snowflake.GenID(),
		AuthorId:    uid,
		AuthorName:  uname,
		CommunityId: p.CommunityId,
		Title:       p.Title,
		Content:     p.Content,
		Status:      mysql.PostPublish,
	}
	if err = mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql CreatePost method is failed",
			zap.Error(err))
		return
	}
	if err = redis.CreatePost(post.AuthorId, post.PostId, post.CommunityId); err != nil {
		zap.L().Error("redis CreatePost method is failed",
			zap.Error(err))
	}
	return
}

// PostPut 修改帖子数据
func (p Publish) PostPut(pid int64) (silr.Response, error) {
	code := e.CodeSUCCESS
	status, err := mysql.GetPostStatus(pid)
	if err != nil {
		code = e.CodeServerBusy
		zap.L().Error("mysql CheckPostStatus method is failed",
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	// 要修改帖子，只能在帖子处于：0：待审核，1：已发布，2：保存状态
	if status == mysql.PostExpired {
		code = e.CodePostVoteExpired
		return silr.Response{Status: code, Msg: code.Msg()}, redis.ErrVoteTimeExpired
	} else if status == mysql.PostDelete {
		code = e.CodeInvalidParams
		return silr.Response{Status: code, Msg: code.Msg()}, mysql.ErrNoRows
	}
	if err = mysql.UpdatePost(pid, p.Title, p.Content); err != nil {
		code = e.CodeServerBusy
		zap.L().Error("mysql UpdatePost method is failed",
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	return silr.Response{Status: code, Msg: code.Msg()}, nil
}

// PostDetailById 根据帖子ID查询到帖子的详情
func (p *PostService) PostDetailById(pid int64) (err error) {
	if p.Post, err = mysql.GetPostDetailById(pid); err != nil {
		zap.L().Error("GetPostDetailById method is failed",
			zap.Int64("postid", pid),
			zap.Error(err))
		return
	}
	if p.CommunityDetail, err = mysql.GetCommunityDetail(p.Post.CommunityId); err != nil {
		zap.L().Error("GetCommunityDetail method is failed",
			zap.Int64("community_id", p.Post.CommunityId),
			zap.Error(err))
		return
	}
	user, err := mysql.GetUserById(p.Post.AuthorId)
	if err != nil {
		zap.L().Error("GetUserById method is failed",
			zap.Int64("author_id", p.Post.AuthorId),
			zap.Error(err))
		return
	}
	// 如果帖子还没有过期，就可以去redis中去查当前帖子的票数
	if p.Post.Status != mysql.PostExpired {
		p.Post.VoteNum = uint32(redis.GetPostVote(pid))
	}
	p.AuthorName = user.Username
	return nil
}

// PostListInOrder 根据排序方法获取所有帖子列表
func PostListInOrder(page, size int64, order string) (postList []*PostService, err error) {
	key := redis.KeyPostTimeZSet
	if order == OrderByScore {
		key = redis.KeyPostScoreZSet
	}
	ids, err := redis.GetPostIds(page, size, key)
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

// DeletePost 删除帖子
func DeletePost(pid int64) (err error) {
	var cid int64
	if cid, err = mysql.FindCidByPid(pid); err == nil {
		if status, _ := mysql.GetPostStatus(pid); status != mysql.PostDelete {
			if status == mysql.PostExpired {
				// 如果帖子已经过期了，只需要对他的状态从已过期更新为已删除即可
				if err = mysql.DeletePost(pid); err != nil {
					zap.L().Error("mysql DeletePost method is err",
						zap.Error(err))
				}
				// 这里可以不修改redis，因为一旦帖子id找不到也就说明该帖子下的所有评论都会找不到
				return
			}
			// 如果帖子没有过期,在进行软删除之前需要将帖子的vote_num数据从redis中取出进行更新
			ticket := redis.GetPostVote(pid)
			if err = mysql.UpdateCtbPost(pid, uint32(ticket)); err != nil {
				zap.L().Error("mysql UpdateCtbPost method is err",
					zap.Error(err))
				return
			}
			if err = mysql.DeletePost(pid); err != nil {
				zap.L().Error("mysql DeletePost method is err",
					zap.Error(err))
				return
			}
			if err = redis.DeletePost(pid, cid); err != nil {
				zap.L().Error("redis DeletePost method is err",
					zap.Error(err))
			}
		}
	}
	return
}

// getPostListByIds 获取帖子列表根据ids
func getPostListByIds(ids []string) (postList []*PostService, err error) {
	if len(ids) == NULLData {
		zap.L().Warn("redis post data is null")
		return
	}
	tickets, err := redis.GetPostVotes(ids)
	posts, err := mysql.GetPostListInOrder(ids)
	if err != nil {
		zap.L().Error("GetPostListInOrder method is err",
			zap.Error(err))
		return
	}
	var community *models.CommunityDetail
	postList = make([]*PostService, NULLData, len(posts))
	for i, post := range posts {
		community, err = mysql.GetCommunityDetail(post.CommunityId)
		if err != nil {
			zap.L().Error("GetCommunityDetail method is err",
				zap.Int64("community_id", post.CommunityId),
				zap.Error(err))
			continue
		}
		post.VoteNum = tickets[i]
		plist := &PostService{
			Post:            post,
			CommunityDetail: community,
		}
		postList = append(postList, plist)
	}
	return
}
