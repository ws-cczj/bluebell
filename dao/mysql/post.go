package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	iStr := `insert into post(
				post_id,title,content,author_id,community_id,status)
				values (?,?,?,?,?,?)`
	_, err = db.Exec(iStr,
		post.PostId,
		post.Title,
		post.Content,
		post.AuthorId,
		post.CommunityId,
		post.Status)
	return
}

// GetPostDetailById 根据ID获取帖子
func GetPostDetailById(id int64) (data *models.Post, err error) {
	data = new(models.Post)
	qStr := `select 
				post_id,community_id,author_id,title,content,status,create_time
				from post
				where post_id = ?`
	err = db.Get(data, qStr, id)
	if err == ErrNoRows {
		zap.L().Error("getPostDetail data is null", zap.Error(err))
		err = ErrorInvalidParam
	}
	return
}

// GetPostList 获取所有帖子
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	qStr := `select 
				post_id,community_id,author_id,title,content,status,create_time
				from post
				order by create_time DESC
				limit ?,?`
	posts = make([]*models.Post, 0, size)
	err = db.Select(&posts, qStr, (page-1)*size, size)
	return
}

// GetPostListInOrder 根据指定顺序查询帖子
func GetPostListInOrder(ids []string) (posts []*models.Post, err error) {
	qStr := `select 
				post_id,community_id,author_id,title,content,status,create_time
				from post
				where post_id in (?)
				order by FIND_IN_SET(post_id, ?)`
	posts = make([]*models.Post, 0, len(ids))
	query, args, err := sqlx.In(qStr, ids, strings.Join(ids, ","))
	query = db.Rebind(query)
	err = db.Select(&posts, query, args...)
	return
}
