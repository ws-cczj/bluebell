package redis

// redis Key 使用命名空间的方式，方便查询和分割
const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZest    = "post:time"   // zest;帖子及发表时间
	KeyPostScoreZest   = "post:score"  // zest;帖子及投票分数
	KeyPostVotedZestPF = "post:voted:" // zest;记录用户投票的参数类型;参数是 post id
)
