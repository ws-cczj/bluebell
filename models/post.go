package models

import "time"

type Post struct {
	CommunityId int       `json:"community_id,omitempty" db:"community_id"`
	VoteNum     int       `json:"vote_num" db:"vote_num"`
	Status      uint8     `json:"status" db:"status"` // 规定 0为审核中, 1为已发布, 2为已保存, 3为已经过期, 4为已删除
	PostId      string    `json:"post_id" db:"post_id"`
	AuthorId    string    `json:"author_id,omitempty" db:"author_id"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	CreateAt    time.Time `json:"createAt" db:"create_time"`
	UpdateAt    time.Time `json:"updateAt" db:"update_time"`
}

type PostPut struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
}

type PostDelete struct {
	CommunityId int    `form:"community_id" binding:"required"`
	Status      uint8  `form:"status" binding:"required"`
	PostId      string `form:"post_id" binding:"required"`
}
