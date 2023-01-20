package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	silr "bluebell/serializer"
	"crypto/md5"
	"encoding/hex"

	"go.uber.org/zap"
)

const secret = "cczjblog.top"

// UserRegister 用户注册
func UserRegister(user *models.UserRegister) (err error) {
	// 1. 校验用户名
	if err = mysql.CheckUsername(user.Username); err != nil {
		zap.L().Error("mysql checkUsername method err", zap.Error(err))
		return
	}
	// 2. 生成UserID
	user.UserId = snowflake.GenID()
	user.Password = encryptPassword(user.Password)
	// 3. 添加用户到数据库
	if err = mysql.InsertUser(user); err != nil {
		zap.L().Error("mysql User Insert method err", zap.Error(err))
	}
	return
}

// UserLogin 用户登录
func UserLogin(user *models.UserLogin) (atoken, rtoken string, err error) {
	pwdParam := encryptPassword(user.Password)
	if err = mysql.CheckLoginInfo(user); err != nil {
		zap.L().Error("mysql checkLoginInfo method err", zap.Error(err))
		return
	}
	if pwdParam != user.Password {
		zap.L().Error("user login password Compared", zap.Error(mysql.ErrorNotComparePwd))
		return
	}
	// 颁发token
	atoken, rtoken, err = jwt.GenToken(user.UserId, user.Username)
	if err != nil {
		return
	}
	// 将token存入redis中一份
	if err = redis.SetSingleUserToken(user.Username, atoken); err != nil {
		zap.L().Error("redis set userlogin token err", zap.Error(err))
	}
	return
}

// UserFollowBuild 用户关注表的构建
func UserFollowBuild(uid int64, follow *models.UserFollow) error {
	if follow.Agree {
		return redis.SetUserToFollow(uid, follow.ToUserId)
	}
	return redis.CancelUserToFollow(uid, follow.ToUserId)
}

// UserToFollowList 获取用户的关注列表
func UserToFollowList(uid int64) ([]*silr.ResponseUserFollow, error) {
	toFollows, err := redis.GetUserToFollows(uid)
	if err != nil {
		zap.L().Error("redis getUserToFollows method err", zap.Error(err))
		return []*silr.ResponseUserFollow{}, err
	}
	return mysql.GetUserFollows(toFollows)
}

// UserFollowList 获取用户的关注列表
func UserFollowList(uid int64) ([]*silr.ResponseUserFollow, error) {
	Follows, err := redis.GetUserFollows(uid)
	if err != nil {
		zap.L().Error("redis GetUserFollows method err", zap.Error(err))
		return []*silr.ResponseUserFollow{}, err
	}
	return mysql.GetUserFollows(Follows)
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

// encryptPassword 对password进行md5加密处理
func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
