package user

import (
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"github.com/gin-gonic/gin"
)

type UserRegisterRequest struct {
	Mobile   string `json:"mobile" binding:"required,mobile" `
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func (us *userServer) Register(ctx *gin.Context) error {
	var cr UserRegisterRequest
	if err := ctx.ShouldBind(&cr); err != nil {
		return gin2.HandleValidatorError(ctx, err, us.trans)

	}

	userDTO, err := us.sf.Users().Register(ctx, cr.Mobile, cr.Password, cr.Code)
	if err != nil {
		return err
	}

	common.OkWithData(ctx, UserResponse{
		ID:        userDTO.ID,
		NickName:  userDTO.NickName,
		Token:     userDTO.Token,
		ExpiredAt: userDTO.ExpiresAt,
	})
	return nil
}
