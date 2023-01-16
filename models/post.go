package models

import "time"

type Post struct {
	PostId      int64     `json:"post_id,string,omitempty" db:"post_id"`
	AuthorId    int64     `json:"author_id,string,omitempty" db:"author_id"`
	CommunityId int64     `json:"community_id,omitempty" db:"community_id"`
	Status      uint8     `json:"status" db:"status"` // 规定 0为审核中, 1为已发布, 2为已保存, 3为已经过期, 4为已删除
	VoteNum     uint32    `json:"vote_num,omitempty" db:"vote_num"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
}
