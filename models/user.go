package models

import "time"

type User struct {
	UserID     int64     `db:"user_id"`
	Username   string    `db:"username"`
	Password   string    `db:"password"`
	Email      string    `db:"email"`
	Gender     uint8     `db:"gender"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}
