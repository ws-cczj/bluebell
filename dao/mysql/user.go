package mysql

import (
	"bluebell/models"
	silr "bluebell/serializer"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// CheckUsername 检查用户名是否重复
func CheckUsername(username string) (err error) {
	var count int
	qStr := `select count(id) from user where username = ?`
	if err = db.Get(&count, qStr, username); err != nil {
		zap.L().Error("CheckUsername db get failed", zap.Error(err))
		return
	}
	if count > 0 {
		zap.L().Debug("username is exist", zap.Error(ErrorUserExist))
		return ErrorUserExist
	}
	return
}

// InsertUser 登记用户信息到数据库
func InsertUser(user *models.UserRegister) (err error) {
	iStr := `insert into 
    			user(user_id,username,password,email,gender)
				values(?,?,?,?,?)`
	_, err = db.Exec(iStr,
		user.UserId,
		user.Username,
		user.Password,
		user.Email,
		user.Gender,
	)
	if err != nil {
		zap.L().Error("create user method is failed", zap.Error(err))
	}
	return
}

// CheckLoginInfo 验证用户登录信息
func CheckLoginInfo(user *models.UserLogin) (err error) {
	qStr := `select user_id,password
				from user 
				where username = ?`
	err = db.Get(user, qStr, user.Username)
	return
}

// GetUserById 根据用户ID查找到用户信息
func GetUserById(uid string) (user *models.UserRegister, err error) {
	user = new(models.UserRegister)
	qStr := `select user_id,username,password,email,gender
				from user 
				where user_id = ?`
	err = db.Get(user, qStr, uid)
	return
}

// GetUserFollows 获取用户的关注|粉丝列表
func GetUserFollows(uids []string) (data []*silr.ResponseUserFollow, err error) {
	if len(uids) == 0 {
		return
	}
	data = make([]*silr.ResponseUserFollow, 0, len(uids))
	qStr := `select user_id, username 
				from user
				where user_id in (?)`
	// TODO 无法按照关注|粉丝顺序返回，find_in_set 不走索引效率太低
	query, args, err := sqlx.In(qStr, uids)
	if err != nil {
		zap.L().Error("sqlx in method err", zap.Error(err))
		return nil, err
	}
	qey := db.Rebind(query)
	err = db.Select(&data, qey, args...)
	return
}

// GetUserCommunityList 获取用户管理的社区信息列表
func GetUserCommunityList(uid string, cidNums int64) (data []*models.Community, err error) {
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
func GetUserPostList(uid string, pidNums int64) (data []*models.Post, err error) {
	data = make([]*models.Post, 0, pidNums)
	qStr := `select post_id,author_name,title,content,status,create_time,update_time
				from post
				where author_id = ? AND status <> ?
				order by create_time DESC`
	if err = db.Select(&data, qStr, uid, PostDelete); err != nil {
		if err == ErrNoRows {
			zap.L().Warn("getCommunityList is null data")
			err = nil
		}
	}
	return
}
