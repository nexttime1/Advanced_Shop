package user

import (
	"Advanced_Shop/app/pkg/common"
	"Advanced_Shop/gnova/server/restserver/middlewares"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

type userDetailResponse struct {
	NickName string `json:"nick_name"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
	Mobile   string `json:"mobile"`
}

func (us *userServer) GetUserDetail(c *gin.Context) error {
	log.Info("GetUserDetail")
	userID, _ := c.Get(middlewares.KeyUserID)

	ctx := c.Request.Context()
	userDTO, err := us.sf.Users().Get(ctx, uint64(userID.(float64)))
	if err != nil {
		return err
	}
	common.OkWithData(c, userDetailResponse{
		NickName: userDTO.NickName,
		Birthday: userDTO.Birthday.Format("2006-01-02"),
		Gender:   userDTO.Gender,
		Mobile:   userDTO.Mobile,
	})
	return nil
}
