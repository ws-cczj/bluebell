package models

import "time"

type Post struct {
	PostId      int64     `json:"post_id,string" db:"post_id"`
	AuthorId    int64     `json:"author_id,string" db:"author_id"`
	CommunityId int64     `json:"community_id" db:"community_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Status      uint8     `json:"status" db:"status"` // 规定 0 为审核中,1 为已发布,2 为已保存
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}