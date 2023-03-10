package mysql

import (
	"bluebell/models"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

const (
	PostCheck = iota
	PostPublish
	PostSave
	PostExpired
	PostDelete
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	iStr := `insert into post(
				post_id,title,content,author_id,author_name,community_id,status)
				values (?,?,?,?,?,?,?)`
	_, err = db.Exec(iStr,
		post.PostId,
		post.Title,
		post.Content,
		post.AuthorId,
		post.AuthorName,
		post.CommunityId,
		post.Status)
	return
}

// DeletePost 软删除帖子
func DeletePost(pid string) (err error) {
	uStr := `update post 
				set status = ?
				where post_id = ?`
	_, err = db.Exec(uStr, PostDelete, pid)
	return
}

// UpdatePost 更新帖子数据
func UpdatePost(pid string, title, content string) (err error) {
	uStr := `update post 
				set title = ?, content = ? 
				where post_id = ?`
	_, err = db.Exec(uStr, title, content, pid)
	return
}

// UpdateCtbPost 更新帖子的票数 -> 过期状态
func UpdateCtbPost(pid string, vote_num int) (err error) {
	uStr := `update post 
				set vote_num = ?, status = ? 
				where post_id = ?`
	_, err = db.Exec(uStr, vote_num, PostExpired, pid)
	return
}

// UpdateAndDeletePost 更新并且删除帖子 -> 软删除状态
func UpdateAndDeletePost(pid string, vote_num int) (err error) {
	uStr := `update post 
				set vote_num = ?, status = ? 
				where post_id = ?`
	_, err = db.Exec(uStr, vote_num, PostDelete, pid)
	return
}

// GetPostDetailById 根据ID获取帖子
func GetPostDetailById(id string) (data *models.Post, err error) {
	data = new(models.Post)
	qStr := `select 
				post_id,community_id,author_id,author_name,title,content,vote_num,status,create_time
				from post
				where post_id = ? and status <> ?`
	err = db.Get(data, qStr, id, PostDelete)
	return
}

// GetPostStatus 根据Id获取帖子状态
func GetPostStatus(pid string) (status uint8, err error) {
	qStr := `select status 
				from post 
				where post_id = ?`
	err = db.Get(&status, qStr, pid)
	return
}

// GetPostListInOrder 根据指定顺序查询帖子
func GetPostListInOrder(ids []string) (posts []*models.Post, err error) {
	qStr := `select 
				post_id,community_id,author_id,author_name,vote_num,title,content,status,create_time,update_time
				from post
				where post_id in (?)
				order by FIND_IN_SET(post_id, ?)`
	posts = make([]*models.Post, 0, len(ids))
	query, args, _ := sqlx.In(qStr, ids, strings.Join(ids, ","))
	query = db.Rebind(query)
	err = db.Select(&posts, query, args...)
	if err == ErrNoRows {
		zap.L().Warn("GetPostList method data is null")
		err = nil
	}
	return
}

// CrontabPostDelete 定实获取这段时间被删除的帖子
func CrontabPostDelete(preT, nowT time.Time) (pids []string, err error) {
	qStr := `select post_id
				from post
				where status = ?
				AND update_time >= ?
				AND update_time <= ?`
	err = db.Select(&pids, qStr, PostDelete, preT, nowT)
	if err != nil {
		zap.L().Error("db query method err", zap.Error(err))
	}
	return
}
