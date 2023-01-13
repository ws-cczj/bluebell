package models

import "time"

type Community struct {
	ID         int64  `json:"id" db:"id"`
	AuthorId   int64  `json:"author_id,string" db:"author_id" form:"author_id"`
	AuthorName string `json:"author_name" db:"author_name" form:"author_name"`
	Name       string `json:"name" db:"community_name" form:"name" binding:"required"`
}

type CommunityDetail struct {
	*Community
	Status       uint8     `json:"status" db:"status" binding:"oneof=0 1 2"` // 规定 0为审核中, 1为已发布, 2为已删除
	Introduction string    `json:"introduction" db:"introduction" form:"introduction" binding:"required"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
	UpdateTime   time.Time `json:"update_time" db:"update_time"`
}
