package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	CommentPost = iota
	CommentFComment
	CommentPeople
)

// CreateComment 创建一条评论
func CreateComment(comment *models.CommentDetail) (int64, error) {
	if comment.Type == CommentPost {
		iStr := `insert into 
    			comment(post_id,type,author_id,author_name,content)
    			values (?,?,?,?,?)`
		iRes, err := db.Exec(iStr,
			comment.Comment.PostId,
			comment.Comment.Type,
			comment.Comment.AuthorId,
			comment.Comment.AuthorName,
			comment.Comment.Content)
		if err != nil {
			return 0, err
		}
		return iRes.LastInsertId()
	}
	iStr := `insert into 
    			comment(father_id,post_id,type,author_id,author_name,to_author_id,to_author_name,content)
    			values (?,?,?,?,?,?,?,?)`
	iRes, err := db.Exec(iStr,
		comment.FatherId,
		comment.Comment.PostId,
		comment.Comment.Type,
		comment.Comment.AuthorId,
		comment.Comment.AuthorName,
		comment.ToAuthorId,
		comment.ToAuthorName,
		comment.Comment.Content)
	if err != nil {
		return 0, err
	}
	return iRes.LastInsertId()
}

// DeleteComment 删除一条评论
func DeleteComment(commentId int) (err error) {
	dStr := `delete from comment 
				where id = ?`
	_, err = db.Exec(dStr, commentId)
	return
}

// CrontabCheckComment 检查不合格的子评论 指的是子评论的父评论是否还在
func CrontabDeleteComment() (err error) {
	dStr := `delete 
				from comment
				where id in
      			(select tb.id
       			from (select id
             			from comment
             			where father_id not in (select id
             				                        from comment
             				                        )
             				) as tb)`
	_, err = db.Exec(dStr)
	return
}

// GetCommentById 通过fid查找一条父评论
func GetCommentById(id string) (Fcomment *models.Comment, err error) {
	qStr := `select id,post_id,type,author_id,author_name,content,create_time,update_time
				from comment
				where id = ?`
	Fcomment = new(models.Comment)
	err = db.Get(Fcomment, qStr, id)
	return
}

// GetCommentByIds 通过fids查找多条父评论
func GetCommentByIds(fids []string) (Fcomments []*models.Comment, err error) {
	qStr := `select id,post_id,type,author_id,author_name,content,create_time,update_time
				from comment
				where id in (?)
				order by FIND_IN_SET(id,?)`
	Fcomments = make([]*models.Comment, 0, len(fids))
	query, args, err := sqlx.In(qStr, fids, strings.Join(fids, ","))
	if err != nil {
		return
	}
	qey := db.Rebind(query)
	err = db.Select(&Fcomments, qey, args...)
	return
}

// GetCommentList 获取评论集合
func GetCommentList(ids []string) (clist []*models.CommentDetail, err error) {
	qStr := `select *
				from comment
				where id in (?)
				order by FIND_IN_SET(id,?)`
	var build strings.Builder
	for i, id := range ids {
		build.WriteString(id)
		if i != len(ids)-1 {
			build.WriteString(",")
		}
	}
	query, args, _ := sqlx.In(qStr, ids, build.String())
	clist = make([]*models.CommentDetail, len(ids), len(ids))
	err = db.Select(&clist, query, args...)
	return
}
