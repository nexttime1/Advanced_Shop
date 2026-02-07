package v1

import (
	"Advanced_Shop/app/pkg/aliyun"
	"Advanced_Shop/app/pkg/options"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type SmsSrv interface {
	SendSms(ctx context.Context, mobile string) error
}

func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func (s *smsService) SendSms(ctx context.Context, mobile string) error {
	err := aliyun.SendCode(mobile, s.smsOpts)
	if err != nil {
		return err
	}
	return nil
}

type smsService struct {
	smsOpts *options.SmsOptions
}

func NewSmsService(smsOpts *options.SmsOptions) SmsSrv {
	return &smsService{smsOpts: smsOpts}
}
