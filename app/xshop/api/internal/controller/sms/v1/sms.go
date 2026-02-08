package v1

import (
	"Advanced_Shop/app/pkg/code"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/service"
	v1 "Advanced_Shop/app/xshop/api/internal/service/sms/v1"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/storage"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"time"
)

type SendSmsRequest struct {
	Mobile string `json:"mobile" binding:"required,mobile"`
	Role   string `json:"role" binding:"required,oneof=1 2 3"` // 1 代表 管理员 2代表普通用户 3代表游客
}

type SmsController struct {
	sf    service.ServiceFactory
	trans ut.Translator
}

func NewSmsController(sf service.ServiceFactory, trans ut.Translator) *SmsController {
	return &SmsController{sf, trans}
}

func (sc *SmsController) SendSms(c *gin.Context) {
	var cr SendSmsRequest
	if err := c.ShouldBind(&cr); err != nil {
		gin2.HandleValidatorError(c, err, sc.trans)
	}

	smsCode := v1.GenerateSmsCode(6)
	err := sc.sf.Sms().SendSms(c, cr.Mobile)
	if err != nil {
		core.WriteErrResponse(c, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}

	//将验证码保存起来 - redis
	rstore := storage.RedisCluster{}
	err = rstore.SetKey(c, cr.Mobile, smsCode, 5*time.Minute)
	if err != nil {
		core.WriteErrResponse(c, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}

	core.OkWithMessage(c, "发送成功")
}
