package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"

	"go.uber.org/zap"
)

const secret = "cczjblog.top"

var (
	ErrorUserExist     = errors.New("用户名已经存在")
	ErrorNotComparePwd = errors.New("用户密码不匹配")
)

// CheckUsername 检查用户名是否重复
func CheckUsername(username string) (err error) {
	var count int
	qStr := "select count(id) from bluebell.user where username = ?"
	err = db.Get(&count, qStr, username)
	if err != nil {
		zap.L().Error("CheckUsername db get failed", zap.Error(err))
		return err
	}
	if count > 0 {
		zap.L().Debug("username is exist", zap.Error(ErrorUserExist))
		return ErrorUserExist
	}
	return nil
}

// InsertUser 登记用户信息到数据库
func InsertUser(user *models.User) (err error) {
	iStr := "insert into bluebell.user(user_id,username,password,email,gender) values(?,?,?,?,?)"
	_, err = db.Exec(iStr,
		user.UserID,
		user.Username,
		encryptPassword(user.Password),
		user.Email,
		user.Gender,
	)
	if err != nil {
		zap.L().Error("create user method is failed", zap.Error(err))
		return err
	}
	return nil
}

// CheckLoginInfo 验证用户登录信息
func CheckLoginInfo(username, password string) error {
	// 1. 通过username找到password
	var oPassword string
	qStr := "select password from bluebell.user where username = ?"
	err := db.Get(&oPassword, qStr, username)
	if err == sql.ErrNoRows {
		zap.L().Error("LoginInfo is not compared", zap.Error(err))
		return err
	}
	if err != nil {
		zap.L().Error("CheckLoginInfo db get failed", zap.Error(err))
		return err
	}
	// 2. 验证password
	if oPassword != encryptPassword(password) {
		zap.L().Error("LoginInfo is not compared", zap.String("用户: ", username), zap.Error(ErrorNotComparePwd))
		return ErrorNotComparePwd
	}
	return nil
}

// encryptPassword 对password进行md5加密处理
func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
