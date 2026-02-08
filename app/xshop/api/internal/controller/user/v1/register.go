package user

import (
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/pkg/common/core"
	"github.com/gin-gonic/gin"
)

type UserRegisterRequest struct {
	Mobile   string `json:"mobile" binding:"required,mobile" `
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func (us *userServer) Register(ctx *gin.Context) {
	var cr UserRegisterRequest
	if err := ctx.ShouldBind(&cr); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}

	userDTO, err := us.sf.Users().Register(ctx, cr.Mobile, cr.Password, cr.Code)
	if err != nil {
		core.WriteErrResponse(ctx, err, nil)
		return
	}

	core.OkWithData(ctx, UserResponse{
		ID:        userDTO.ID,
		NickName:  userDTO.NickName,
		Token:     userDTO.Token,
		ExpiredAt: userDTO.ExpiresAt,
	})

}
