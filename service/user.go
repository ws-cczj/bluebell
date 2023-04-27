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

type User struct {
}

func NewUserInstance() *User {
	return &User{}
}

// Register 用户注册
func (u User) Register(user *models.UserRegister) (err error) {
	// 1. 校验用户名
	if err = mysql.CheckUsername(user.Username); err != nil {
		zap.L().Error("mysql checkUsername method err", zap.Error(err))
		return
	}
	// 2. 生成UserID
	user.UserId = snowflake.GenID()
	user.Password = Md5Password(user.Password)
	// 3. 添加用户到数据库
	if err = mysql.InsertUser(user); err != nil {
		zap.L().Error("mysql User Insert method err", zap.Error(err))
	}
	return
}

// Login 用户登录
func (User) Login(user *models.UserLogin) (atoken, rtoken string, err error) {
	pwdParam := Md5Password(user.Password)
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

// FollowBuild 用户关注表的构建
func (User) FollowBuild(uid string, follow *models.UserFollow) error {
	if follow.Agree {
		return redis.SetUserToFollow(uid, follow.ToUserId)
	}
	return redis.CancelUserToFollow(uid, follow.ToUserId)
}

// ToFollowList 获取用户的关注列表
func (User) ToFollowList(uid string) ([]*silr.ResponseUserFollow, error) {
	toFollows, err := redis.GetUserToFollows(uid)
	if err != nil {
		zap.L().Error("redis getUserToFollows method err", zap.Error(err))
		return []*silr.ResponseUserFollow{}, err
	}
	return mysql.GetUserFollows(toFollows)
}

// FollowList 获取用户的关注列表
func (User) FollowList(uid string) ([]*silr.ResponseUserFollow, error) {
	Follows, err := redis.GetUserFollows(uid)
	if err != nil {
		zap.L().Error("redis GetUserFollows method err", zap.Error(err))
		return []*silr.ResponseUserFollow{}, err
	}
	return mysql.GetUserFollows(Follows)
}

// CommunityList 获取用户管理的所有社区
func (User) CommunityList(uid string) ([]*models.Community, error) {
	cidNums := redis.GetUserCommunitys(uid)
	return mysql.GetUserCommunityList(uid, cidNums)
}

// PostList 获取用户发布的所有帖子列表
func (User) PostList(uid string) ([]*models.Post, error) {
	pidNums := redis.GetUserPostNums(uid)
	return mysql.GetUserPostList(uid, pidNums)
}

const secret = "cczjblog.top"

// Md5Password 对password进行md5加密处理
func Md5Password(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
