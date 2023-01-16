package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/e"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	silr "bluebell/serializer"
	"errors"
)

type RegisterService struct {
	Username   string `json:"username" form:"username" binding:"required"`
	Password   string `json:"password"  form:"password" binding:"required"`
	RePassword string `json:"re_password" form:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email" form:"email" binding:"required"`
	Gender     uint8  `json:"gender" form:"gender" binding:"required"`
}

type LoginService struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// Register 用户注册
func (r RegisterService) Register() (silr.Response, error) {
	code := e.CodeSUCCESS
	// 1. 校验用户名
	if err := mysql.CheckUsername(r.Username); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			code = e.CodeExistUser
			return silr.Response{Status: code, Msg: err.Error()}, err
		}
		code = e.CodeServerBusy
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	// 2. 生成UserID
	id := snowflake.GenID()
	// 3. 添加用户到数据库
	u := &models.User{
		UserID:   id,
		Username: r.Username,
		Password: r.Password,
		Email:    r.Email,
		Gender:   r.Gender,
	}
	if err := mysql.InsertUser(u); err != nil {
		code = e.CodeServerBusy
		return silr.Response{Status: code, Msg: code}, err
	}
	return silr.Response{Status: code, Data: nil, Msg: code.Msg()}, nil
}

// Login 用户登录
func (l LoginService) Login() (silr.Response, error) {
	code := e.CodeSUCCESS
	user := &models.User{
		Username: l.Username,
		Password: l.Password,
	}
	if err := mysql.CheckLoginInfo(user); err != nil {
		if errors.Is(err, mysql.ErrorNotComparePwd) {
			code = e.CodeNotComparePassword
		} else if err == mysql.ErrNoRows {
			code = e.CodeNotExistUser
		} else {
			code = e.CodeServerBusy
		}
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	// 颁发token
	token, rtoken, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		code = e.TokenFailGenerate
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	// 将token存入redis中一份
	if err = redis.SetSingleUserToken(user.Username, token); err != nil {
		code = e.CodeServerBusy
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	return silr.Response{Data: []string{token, rtoken}}, nil
}

// UserCommunityList 获取用户管理的所有社区
func UserCommunityList(uid int64) ([]*models.Community, error) {
	cidNums := redis.GetUserCommunitys(uid)
	return mysql.GetUserCommunityList(uid, cidNums)
}

// UserPostList 获取用户发布的所有帖子列表
func UserPostList(uid int64) ([]*models.Post, error) {
	pidNums := redis.GetUserPostNums(uid)
	return mysql.GetUserPostList(uid, pidNums)
}
