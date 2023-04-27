package models

import "time"

type Community struct {
	ID         int    `json:"id" db:"id"`
	AuthorId   string `json:"author_id" db:"author_id" form:"author_id"`
	AuthorName string `json:"author_name" db:"author_name" form:"author_name"`
	Name       string `json:"name" db:"community_name" form:"name" binding:"required"`
}

type CommunityDetail struct {
	*Community
	Status       uint8     `json:"status" db:"status" binding:"oneof=0 1 2"` // 规定 0为审核中, 1为已发布, 2为已删除
	Introduction string    `json:"introduction" form:"introduction" db:"introduction" binding:"required"`
	CreateAt     time.Time `json:"createAt" db:"create_time"`
	UpdateAt     time.Time `json:"updateAt" db:"update_time"`
}

// NewCommunityDetail 实例化结构体
func NewCommunityDetail() *CommunityDetail {
	return &CommunityDetail{
		Community: new(Community),
	}
}
