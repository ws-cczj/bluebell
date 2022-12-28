package service

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/e"
	"bluebell/pkg/snowflake"
	silr "bluebell/serializer"

	"go.uber.org/zap"
)

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
		Status:      1,
	}
	if err := mysql.CreatePost(post); err != nil {
		code = e.CodeServerBusy
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	return silr.Response{Status: code, Msg: code.Msg()}, nil
}

// PostDetailById 根据帖子ID查询到帖子的详情
func (p *PostService) PostDetailById(pid int64) error {
	post, err := mysql.GetPostDetailById(pid)
	if err != nil {
		zap.L().Error("GetPostDetailById method is failed",
			zap.Int64("postid", pid),
			zap.Error(err))
		return err
	}
	community, err := mysql.GetCommunityDetail(post.CommunityId)
	if err != nil {
		zap.L().Error("GetCommunityDetail method is failed",
			zap.Int64("community_id", post.CommunityId),
			zap.Error(err))
		return err
	}
	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		zap.L().Error("GetUserById method is failed",
			zap.Int64("author_id", post.AuthorId),
			zap.Error(err))
		return err
	}
	p.AuthorName = user.Username
	p.Post = post
	p.CommunityDetail = community
	return nil
}

// PostList 获取所有帖子列表
func (p *PostService) PostList(page, size int64) ([]*PostService, error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("GetPostDetailById method is failed",
			zap.Int64("page", page),
			zap.Int64("size", size),
			zap.Error(err))
		return nil, err
	}
	var community *models.CommunityDetail
	var user *models.User
	postList := make([]*PostService, 0, len(posts))
	for _, post := range posts {
		community, err = mysql.GetCommunityDetail(post.CommunityId)
		if err != nil {
			zap.L().Error("GetCommunityDetail method is failed",
				zap.Int64("community_id", post.CommunityId),
				zap.Error(err))
			continue
		}
		user, err = mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("GetUserById method is failed",
				zap.Int64("author_id", post.AuthorId),
				zap.Error(err))
			continue
		}
		plist := &PostService{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		postList = append(postList, plist)
	}
	return postList, nil
}
