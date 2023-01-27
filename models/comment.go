package models

import "time"

type Comment struct {
	Type        int8      `json:"type" db:"type" form:"type" binding:"oneof=0 1 2"` // 规定 0 为对帖子评论, 1 为对父评论评论, 2为对人评论
	Id          int       `json:"id" db:"id" form:"id"`
	FavoriteNum int       `json:"favorite_num"`
	PostId      string    `json:"post_id"  db:"post_id" form:"post_id" binding:"required"`
	AuthorId    string    `json:"author_id" db:"author_id"` // 该条评论的作者
	AuthorName  string    `json:"author_name" db:"author_name"`
	Content     string    `json:"content" db:"content" form:"content" binding:"required"`
	CreateAt    time.Time `json:"CreateAt" db:"create_time"`
	UpdateAt    time.Time `json:"UpdateAt" db:"update_time"`
}

type CommentDetail struct {
	FatherId     int    `json:"father_id,omitempty" db:"father_id" form:"father_id"`
	ToAuthorId   string `json:"to_author_id,omitempty" db:"to_author_id" form:"to_author_id"` // 被回复的作者
	ToAuthorName string `json:"to_author_name,omitempty" db:"to_author_name" form:"to_author_name"`
	*Comment
}

type Favorite struct {
	Agree      bool   `form:"agree"`                    // true 表示点赞，false表示取消点赞
	Type       int8   `form:"type" binding:"oneof=0 1"` // 表示对人或者对帖子点赞
	Id         int    `form:"id" binding:"required"`
	PostId     string `form:"post_id"` // 对于子评论点赞无需帖子Id
	ToAuthorId string `form:"to_author_id" binding:"required"`
}

type CommentDelete struct {
	Type   int8   `form:"type" binding:"oneof=0 1 2"`
	Id     int    `form:"id" binding:"required"`
	TypeId string `form:"type_id" binding:"required"` // 如果type是对人则id为fId，否则id为post_id
}
