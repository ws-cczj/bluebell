package service

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/e"
	"bluebell/pkg/utils"
	silr "bluebell/serializer"
	"database/sql"
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
	//Token    string `json:"token"`
}

// Register 用户注册
func (service *RegisterService) Register() (silr.Response, error) {
	code := e.CodeSUCCESS
	// 1. 校验用户名
	if err := mysql.CheckUsername(service.Username); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			code = e.CodeExistUser
			return silr.Response{Status: code, Msg: err.Error()}, err
		}
		code = e.ErrorQueryDatabase
		return silr.Response{Status: code, Msg: code.Msg()}, err
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
		code = e.ErrorExecDatabase
		return silr.Response{Status: code, Msg: code}, err
	}
	return silr.Response{Status: code, Data: nil, Msg: code.Msg()}, nil
}

// Login 用户登录
func (service *LoginService) Login() (silr.Response, error) {
	code := e.CodeSUCCESS
	if err := mysql.CheckLoginInfo(service.Username, service.Password); err != nil {
		if errors.Is(err, mysql.ErrorNotComparePwd) {
			code = e.CodeNotComparePassword
		} else if err == sql.ErrNoRows {
			code = e.CodeNotExistUser
		} else {
			code = e.CodeServerBusy
		}
		return silr.Response{Status: code, Msg: code.Msg()}, err
	}
	return silr.Response{Status: code, Msg: code.Msg()}, nil
}
