package models

type UserRegister struct {
	UserId     int64  `form:"user_id"`
	Username   string `form:"username" binding:"required"`
	Password   string `form:"password" binding:"required"`
	RePassword string `form:"re_password" binding:"required,eqfield=Password"`
	Email      string `form:"email" binding:"required"`
	Gender     uint8  `form:"gender" binding:"oneof=1 0 -1"` // 1 为男 0为女 -1为未知
}

type UserLogin struct {
	UserId   int64  `form:"user_id" db:"user_id"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" db:"password" binding:"required"`
}

type UserFollow struct {
	Agree    bool  `form:"agree"` // 规定 true 为关注, false 为取消关注
	ToUserId int64 `form:"to_user_id" binding:"required"`
}
