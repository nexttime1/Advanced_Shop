package user

import (
	"net/http"

	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

type CaptchaResponse struct {
	CaptchaId     string `json:"captcha_id"`
	CaptchaBase64 string `json:"captcha_base64"`
}

func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	captchaId, base64s, answer, err := cp.Generate()
	if err != nil {
		log.Errorf("生成验证码错误,: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
	}

	log.Infof("Answer : %s", answer)
	ctx.JSON(http.StatusOK, CaptchaResponse{
		CaptchaId:     captchaId,
		CaptchaBase64: base64s,
	})
}
