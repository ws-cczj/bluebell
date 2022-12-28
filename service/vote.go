package service

type PostVoteService struct {
	PostId    int64 `json:"post_id,string" form:"post_id" bidding:"required"`
	Direction int8  `json:"direction" form:"direction" bidding:"required,oneof=1 0 -1"` // 规定 1为赞成，0为取消投票，-1为反对
}

func (v PostVoteService) VoteBuild() {

}
