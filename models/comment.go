package models

import "time"

type Comment struct {
	Type        int8      `json:"type" db:"type" form:"type" binding:"oneof=0 1"` // 规定 0 为对帖子评论, 1 为对人评论
	Id          int64     `json:"id" db:"id" form:"id"`
	PostId      int64     `json:"post_id,string"  db:"post_id" form:"post_id" binding:"required"`
	AuthorId    int64     `json:"author_id,string" db:"author_id"` // 该条评论的作者
	FavoriteNum int64     `json:"favorite_num"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	Content     string    `json:"content" db:"content" form:"content" binding:"required"`
	CreateAt    time.Time `json:"CreateAt" db:"create_time"`
	UpdateAt    time.Time `json:"UpdateAt" db:"update_time"`
}

type CommentDetail struct {
	FatherId     int64  `json:"father_id,omitempty" db:"father_id" form:"father_id"`
	ToAuthorId   int64  `json:"to_author_id,string,omitempty" db:"to_author_id" form:"to_author_id"` // 被回复的作者
	ToAuthorName string `json:"to_author_name,omitempty" db:"to_author_name" form:"to_author_name"`
	*Comment
}

type Favorite struct {
	Agree      bool  `form:"agree"`                    // true 表示点赞，false表示取消点赞
	Type       int8  `form:"type" binding:"oneof=0 1"` // 表示对人或者对帖子点赞
	Id         int64 `form:"id" binding:"required"`
	PostId     int64 `form:"post_id"` // 对于子评论点赞无需帖子Id
	ToAuthorId int64 `form:"to_author_id" binding:"required"`
}

type CommentDelete struct {
	Type   int8  `form:"type" binding:"oneof=0 1"`
	Id     int64 `form:"id" binding:"required"`
	TypeId int64 `form:"type_id" binding:"required"` // 如果type是对人则id为fId，否则id为post_id
}
