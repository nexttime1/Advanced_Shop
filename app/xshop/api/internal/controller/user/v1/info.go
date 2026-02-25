package user

import (
	"Advanced_Shop/gnova/server/restserver/middlewares"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

type userDetailResponse struct {
	NickName string `json:"nick_name"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
	Mobile   string `json:"mobile"`
}

func (us *userServer) GetUserDetail(ctx *gin.Context) {
	log.Info("GetUserDetail")
	userID, _ := ctx.Get(middlewares.KeyUserID)
	userDTO, err := us.sf.Users().Get(ctx, uint64(userID.(float64)))
	if err != nil {
		core.WriteErrResponse(ctx, err, nil)
		return
	}
	core.OkWithData(ctx, userDetailResponse{
		NickName: userDTO.NickName,
		Birthday: userDTO.Birthday.Format("2006-01-02"),
		Gender:   userDTO.Gender,
		Mobile:   userDTO.Mobile,
	})
}
