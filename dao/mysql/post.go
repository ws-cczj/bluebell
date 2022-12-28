package mysql

import (
	"bluebell/models"
	"database/sql"

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
	if err == sql.ErrNoRows {
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
				limit ?,?`
	posts = make([]*models.Post, 0, size)
	err = db.Select(&posts, qStr, (page-1)*size, size)
	return
}
