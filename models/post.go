package models

import "time"

type Post struct {
	PostId      int64     `json:"post_id,string" db:"post_id"`
	AuthorId    int64     `json:"author_id,string" db:"author_id"`
	CommunityId int64     `json:"community_id" db:"community_id"`
	VoteNum     int64     `json:"vote_num" db:"vote_num"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Status      uint8     `json:"status" db:"status"` // 规定 0 为审核中, 1为已发布, 2为已保存, 3为已经过期
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}
