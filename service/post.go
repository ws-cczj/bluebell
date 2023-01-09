package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/e"
	"bluebell/pkg/snowflake"
	silr "bluebell/serializer"
	"errors"

	"go.uber.org/zap"
)

const (
	OPostTime  = "time"
	OPostScore = "score"
)

var ErrPostExpired = errors.New("该帖子已经过期")

type PublishService struct {
	AuthorId    int64  `json:"author_id,string" form:"author_id"`
	CommunityId int64  `json:"community_id" form:"community_id" bidding:"required"`
	Title       string `json:"title" form:"title" bidding:"required"`
	Content     string `json:"content" form:"content" bidding:"required"`
}

type PostService struct {
	AuthorName string `json:"author_name"`
	*models.Post
	*models.CommunityDetail `json:"community_detail"`
}

// PublishPost 发布帖子
func (p PublishService) PublishPost() (silr.Response, error) {
	code := e.CodeSUCCESS
	post := &models.Post{
		PostId:      snowflake.GenID(),
		AuthorId:    p.AuthorId,
		CommunityId: p.CommunityId,
		Title:       p.Title,
		Content:     p.Content,
		Status:      mysql.PostPublish,
	}
	if err := mysql.CreatePost(post); err != nil {
		code = e.CodeServerBusy
		zap.L().Error("mysql CreatePost method is failed",
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	if err := redis.CreatePost(post.PostId, post.CommunityId); err != nil {
		code = e.CodeServerBusy
		zap.L().Error("redis CreatePost method is failed",
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	return silr.Response{Status: code, Msg: code.Msg()}, nil
}

// PostPut 修改帖子数据
func (p PublishService) PostPut(pid int64) (silr.Response, error) {
	code := e.CodeSUCCESS
	status, err := mysql.GetPostStatus(pid)
	if err != nil {
		code = e.CodeServerBusy
		zap.L().Error("mysql CheckPostStatus method is failed",
			zap.Error(err))
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	if status == mysql.PostExpired {
		code = e.CodePostVoteExpired
		return silr.Response{Status: code, Msg: code.Msg()}, ErrPostExpired
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
	p.Post, err = mysql.GetPostDetailById(pid)
	if err != nil {
		zap.L().Error("GetPostDetailById method is failed",
			zap.Int64("postid", pid),
			zap.Error(err))
		return err
	}
	p.CommunityDetail, err = mysql.GetCommunityDetail(p.Post.CommunityId)
	if err != nil {
		zap.L().Error("GetCommunityDetail method is failed",
			zap.Int64("community_id", p.Post.CommunityId),
			zap.Error(err))
		return err
	}
	user, err := mysql.GetUserById(p.Post.AuthorId)
	if err != nil {
		zap.L().Error("GetUserById method is failed",
			zap.Int64("author_id", p.Post.AuthorId),
			zap.Error(err))
		return err
	}
	if p.Post.Status != mysql.PostExpired {
		p.Post.VoteNum = redis.GetPostVote(pid)
	} else {
		if p.Post.VoteNum, err = mysql.GetPostVote(pid); err != nil {
			zap.L().Error("GetPostVote method is failed",
				zap.Int64("post_id", p.Post.PostId),
				zap.Error(err))
			return err
		}
	}
	p.AuthorName = user.Username
	return nil
}

// PostListInOrder 根据排序方法获取所有帖子列表
func (p *PostService) PostListInOrder(page, size int64, order string) (postList []*PostService, err error) {
	key := redis.KeyPostTimeZSet
	if order == OPostScore {
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

// CommunityPostListInOrder 根据顺序获取社区的帖子列表
func (p *PostService) CommunityPostListInOrder(page, size, cid int64, order string) (postList []*PostService, err error) {
	key := redis.KeyPostTimeZSet
	if order == OPostScore {
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

// DeletePost 删除帖子
func DeletePost(pid int64) (err error) {
	var cid int64
	if cid, err = mysql.FindCidByPid(pid); err == nil {
		if err = mysql.DeletePost(pid); err != nil {
			zap.L().Error("mysql DeletePost method is err",
				zap.Error(err))
			return
		}
		if status, _ := mysql.GetPostStatus(pid); status == mysql.PostExpired {
			return
		}
		if err = redis.DeletePost(pid, cid); err != nil {
			zap.L().Error("redis DeletePost method is err",
				zap.Error(err))
		}
	}
	return
}

// getPostListByIds 获取帖子列表根据ids
func getPostListByIds(ids []string) (postList []*PostService, err error) {
	if len(ids) == 0 {
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
	var user *models.User
	postList = make([]*PostService, 0, len(posts))
	for i, post := range posts {
		community, err = mysql.GetCommunityDetail(post.CommunityId)
		if err != nil {
			zap.L().Error("GetCommunityDetail method is err",
				zap.Int64("community_id", post.CommunityId),
				zap.Error(err))
			continue
		}
		user, err = mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("GetUserById method is err",
				zap.Int64("author_id", post.AuthorId),
				zap.Error(err))
			continue
		}
		post.VoteNum = tickets[i]
		plist := &PostService{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		postList = append(postList, plist)
	}
	return
}
