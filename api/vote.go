package api

import (
	"bluebell/pkg/e"
	"bluebell/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PostVotesHandler 帖子投票情况
func PostVotesHandler(c *gin.Context) {
	v := new(service.PostVoteService)
	if err := c.ShouldBind(v); err != nil {
		errs := err.(validator.ValidationErrors)
		ResponseErrorWithMsg(c, e.CodeInvalidParams, errs.Translate(trans))
		return
	}
	v.VoteBuild()
	ResponseSuccess(c, nil)
}
