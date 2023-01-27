package models

import "time"

type Post struct {
	PostId      int64     `json:"post_id,string" db:"post_id"`
	AuthorId    int64     `json:"author_id,string,omitempty" db:"author_id"`
	CommunityId int64     `json:"community_id,omitempty" db:"community_id"`
	Status      uint8     `json:"status" db:"status"` // 规定 0为审核中, 1为已发布, 2为已保存, 3为已经过期, 4为已删除
	VoteNum     uint32    `json:"vote_num" db:"vote_num"`
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
	Status      uint8 `form:"status" binding:"required"`
	PostId      int64 `form:"post_id" binding:"required"`
	CommunityId int64 `form:"community_id" binding:"required"`
}
