package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"encoding/hex"

	"go.uber.org/zap"
)

const secret = "cczjblog.top"

// CheckUsername 检查用户名是否重复
func CheckUsername(username string) (err error) {
	var count int
	qStr := `select count(id) from user where username = ?`
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
	iStr := `insert into user(user_id,username,password,email,gender) values(?,?,?,?,?)`
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
func CheckLoginInfo(user *models.User) error {
	// 1. 通过username找到password
	var oPassword = user.Password
	qStr := `select user_id,username,password 
				from user 
				where username = ?`
	err := db.Get(user, qStr, user.Username)
	if err == ErrNoRows {
		zap.L().Error("LoginInfo is not compared", zap.Error(err))
		return err
	}
	if err != nil {
		zap.L().Error("CheckLoginInfo db get failed", zap.Error(err))
		return err
	}
	// 2. 验证password
	if encryptPassword(oPassword) != user.Password {
		zap.L().Error("LoginInfo is not compared", zap.String("用户名: ", user.Username), zap.Error(ErrorNotComparePwd))
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

// GetUserById 根据用户ID查找到用户信息
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	qStr := `select user_id,username,password,email,gender
				from user 
				where user_id = ?`
	err = db.Get(user, qStr, uid)
	return
}

// GetUserCommunityList 获取用户管理的社区信息列表
func GetUserCommunityList(uid, cidNums int64) (data []*models.Community, err error) {
	data = make([]*models.Community, 0, cidNums)
	qStr := `select id,author_id,author_name,community_name
				from community
				where author_id = ?
				order by create_time DESC`
	if err = db.Select(&data, qStr, uid); err != nil {
		if err == ErrNoRows {
			zap.L().Warn("getCommunityList is null data")
			err = nil
		}
	}
	return
}

// GetUserPostList 获取用户管理的帖子列表
func GetUserPostList(uid, pidNums int64) (data []*models.Post, err error) {
	data = make([]*models.Post, 0, pidNums)
	qStr := `select author_name,title,content,status,create_time,update_time
				from post
				where author_id = ?
				order by create_time DESC`
	if err = db.Select(&data, qStr, uid); err != nil {
		if err == ErrNoRows {
			zap.L().Warn("getCommunityList is null data")
			err = nil
		}
	}
	return
}
