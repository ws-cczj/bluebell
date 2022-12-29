package redis

import (
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	OneWeekPostTime = 7 * 24 * 3600
	OneTicketScore  = 432 // 86400/200 当有200个赞成的时候，就将帖子置为热门贴
)

// CheckVoteTime 查看帖子投票时间是否过期
func CheckVoteTime(pid int64) error {
	postTime := rdb.ZScore(addKeyPrefix(KeyPostTimeZest), stvI64toa(pid)).Val()
	if float64(time.Now().Unix())-postTime > OneWeekPostTime {
		return ErrVoteTimeExpired
	}
	return nil
}

// PostVoteDirect 判断该用户对帖子的投票情况
func PostVoteDirect(pid, uid int64, direct float64) (diff float64) {
	oldDirect := rdb.ZScore(addKeyPrefix(KeyPostVotedZestPF+stvI64toa(pid)), stvI64toa(uid)).Val()
	if oldDirect == direct {
		return
	}
	diff = math.Abs(oldDirect - direct)
	if direct > oldDirect {
		return
	}
	return -diff
}

// ChangePostInfo 更改帖子分数和用户投票情况
func ChangePostInfo(pid, uid int64, diff, direct float64) (err error) {
	pipe := rdb.TxPipeline()
	if diff == 0 && direct != 0 {
		pipe.ZRem(addKeyPrefix(KeyPostVotedZestPF+stvI64toa(pid)), stvI64toa(uid))
		pipe.ZIncrBy(addKeyPrefix(KeyPostScoreZest), -direct*OneTicketScore, stvI64toa(pid))
	} else {
		pipe.ZIncrBy(addKeyPrefix(KeyPostScoreZest), diff*OneTicketScore, stvI64toa(pid))
		pipe.ZAdd(addKeyPrefix(KeyPostVotedZestPF+stvI64toa(pid)), redis.Z{
			Score:  direct,
			Member: uid,
		})
	}
	_, err = pipe.Exec()
	return
}
