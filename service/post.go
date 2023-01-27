package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/e"
	"bluebell/pkg/snowflake"
	silr "bluebell/serializer"
	"errors"
	"strconv"

	"go.uber.org/zap"
)

const (
	OrderByTime  = "time"
	OrderByScore = "score"
)

type Post struct {
}

func NewPostInstance() *Post {
	return &Post{}
}

var ErrPostDelete = errors.New("post already delete, please dont repeat to click")

// Put 修改帖子数据
func (Post) Put(pid string, p *models.PostPut) (silr.Response, error) {
	code := e.CodeSUCCESS
	status, err := mysql.GetPostStatus(pid)
	if err != nil {
		code = e.CodeServerBusy
		zap.L().Error("mysql GetPostStatus method is failed",
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

// ListInOrder 根据排序方法获取所有帖子列表
func (p Post) ListInOrder(page, size int64, order string) (postList []*PostAll, err error) {
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
	return p.getPostListByIds(ids)
}

// Delete 删除帖子
func (Post) Delete(uid string, post *models.PostDelete) (err error) {
	if post.Status == mysql.PostDelete {
		zap.L().Error("delete data is null", zap.Error(ErrPostDelete))
		return ErrPostDelete
	} else if post.Status == mysql.PostExpired {
		// 如果帖子已经过期了，只需要对他的状态从已过期更新为已删除即可
		if err = mysql.DeletePost(post.PostId); err != nil {
			zap.L().Error("mysql DeletePost method err",
				zap.Error(err))
			return
		}
		// 修改redis中post结构
		if err = redis.PostExpiredDelete(uid, post.PostId, strconv.Itoa(post.CommunityId)); err != nil {
			zap.L().Error("redis PostExpiredDelete method is err",
				zap.Error(err))
		}
	} else {
		// 如果帖子没有过期,在进行软删除之前需要将帖子的vote_num数据从redis中取出进行更新
		ticket := redis.GetPostVote(post.PostId)
		if err = mysql.UpdateAndDeletePost(post.PostId, int(ticket)); err != nil {
			zap.L().Error("mysql UpdateAndDeletePost method is err",
				zap.Error(err))
			return
		}
		if err = redis.PostDelete(uid, post.PostId, strconv.Itoa(post.CommunityId)); err != nil {
			zap.L().Error("redis DeletePost method is err",
				zap.Error(err))
		}
	}
	return
}

// getPostListByIds 获取帖子列表根据ids
func (Post) getPostListByIds(ids []string) (postList []*PostAll, err error) {
	if len(ids) <= 0 {
		zap.L().Warn("redis post data is null")
		return
	}
	// 获取到已经过期和未过期的帖子
	tickets, err := redis.GetPostVotes(ids)
	posts, err := mysql.GetPostListInOrder(ids)
	if err != nil {
		zap.L().Error("GetPostListInOrder method is err",
			zap.Error(err))
		return
	}
	postList = make([]*PostAll, 0, len(posts))
	for i, post := range posts {
		community, err := mysql.GetCommunityDetail(post.CommunityId)
		if err != nil {
			zap.L().Error("GetCommunityDetail method is err",
				zap.Int("community_id", post.CommunityId),
				zap.Error(err))
			continue
		}
		// 如果是还没有过期的，就直接将票数进行赋值
		if post.Status == mysql.PostPublish {
			post.VoteNum = tickets[i]
		}
		plist := &PostAll{
			Post:            post,
			CommunityDetail: community,
		}
		postList = append(postList, plist)
	}
	return
}

type Publish struct {
	CommunityId int    `form:"community_id" bidding:"required"`
	Title       string `form:"title" bidding:"required"`
	Content     string `form:"content" bidding:"required"`
}

// Publish 发布帖子
func (p Publish) Publish(uid, uname string) (err error) {
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
	if err = redis.PostCreate(post.AuthorId, post.PostId, strconv.Itoa(post.CommunityId)); err != nil {
		zap.L().Error("redis PostCreate method is failed",
			zap.Error(err))
	}
	return
}

type PostAll struct {
	*models.Post
	*models.CommunityDetail `json:"community_detail"`
}

// DetailById 根据帖子ID查询到帖子的详情
func (p *PostAll) DetailById(pid string) (err error) {
	if p.Post, err = mysql.GetPostDetailById(pid); err != nil {
		zap.L().Error("GetPostDetailById method is failed",
			zap.String("postid", pid),
			zap.Error(err))
		return
	}
	if p.CommunityDetail, err = mysql.GetCommunityDetail(p.Post.CommunityId); err != nil {
		zap.L().Error("GetCommunityDetail method is failed",
			zap.Int("community_id", p.Post.CommunityId),
			zap.Error(err))
		return
	}
	// 如果帖子还没有过期，就可以去redis中去查当前帖子的票数
	if p.Post.Status != mysql.PostExpired {
		p.Post.VoteNum = int(redis.GetPostVote(pid))
	}
	return
}
