package user

import (
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginRequest struct {
	Mobile    string `json:"mobile" binding:"required,mobile" ` //以使用 binding:"mobile" 这样的标签  自动调用验证函数
	Password  string `json:"password" binding:"required"`
	CaptchaId string `json:"captcha_id" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}

type UserResponse struct {
	ID        uint64 `json:"id"`
	NickName  string `json:"nick_name"`
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}

func (us *userServer) Login(ctx *gin.Context) {
	log.Info("login is called")

	var cr UserLoginRequest
	if err := ctx.ShouldBind(&cr); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}

	//验证码验证
	verifyResult := store.Verify(cr.CaptchaId, cr.Answer, true)
	if !verifyResult {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	userDTO, err := us.sf.Users().MobileLogin(ctx, cr.Mobile, cr.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "登录失败",
		})
		return
	}
	core.OkWithData(ctx, UserResponse{
		ID:        userDTO.ID,
		NickName:  userDTO.NickName,
		Token:     userDTO.Token,
		ExpiredAt: userDTO.ExpiresAt,
	})

}
