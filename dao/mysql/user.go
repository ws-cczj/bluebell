package mysql

import (
	"bluebell/models"
	"bluebell/pkg/e"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"go.uber.org/zap"
)

const secret = "cczjblog.top"

func CheckUsername(username string) (err error) {
	var count int
	qStr := "select count(id) from bluebell.user where username = ?"
	err = db.Get(&count, qStr, username)
	if err != nil {
		zap.L().Error("checkUsername get method is failed", zap.Error(err))
		return err
	}
	if count > 0 {
		return errors.New(e.GetMsg(e.ErrorExistUser))
	}
	return nil
}

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

// encryptPassword 对password进行md5加密处理
func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
