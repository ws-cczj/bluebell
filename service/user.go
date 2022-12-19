package service

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/e"
	"bluebell/pkg/utils"
	"bluebell/serializer"
)

type RegisterService struct {
	Username   string `json:"username" form:"username" binding:"required"`
	Password   string `json:"password"  form:"password" binding:"required"`
	RePassword string `json:"re_password" form:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email" form:"email" binding:"required"`
	Gender     uint8  `json:"gender" form:"gender" binding:"required"`
}

type LoginService struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token"`
}

// Register 用户注册
func (service *RegisterService) Register() (serializer.Response, error) {
	code := e.SUCCESS
	// 1. 校验用户名
	if err := mysql.CheckUsername(service.Username); err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}, err
	}
	// 2. 生成UserID
	id := utils.GenID()
	// 3. 生成token
	// 4. 添加用户到数据库
	u := &models.User{
		UserID:   id,
		Username: service.Username,
		Password: service.Password,
		Email:    service.Email,
		Gender:   service.Gender,
	}
	if err := mysql.InsertUser(u); err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}, err
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}, nil
}

// Login 用户登录
func (service *LoginService) Login() (serializer.Response, error) {
	return serializer.Response{}, nil
}
