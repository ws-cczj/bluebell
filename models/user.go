package models

type User struct {
	UserID     int64  `db:"user_id"`
	Username   string `db:"username"`
	Password   string `db:"password"`
	Email      string `db:"email"`
	Gender     uint8  `db:"gender"`
	CreateTime string `db:"create_time"`
	UpdateTime string `db:"update_time"`
}
