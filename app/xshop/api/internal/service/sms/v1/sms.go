package v1

import (
	"context"
	"fmt"
	"math/rand"
	"mxshop/app/pkg/options"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type SmsSrv interface {
	//
	// SendSms
	//  @Description: 发送短信验证码
	//  @param ctx
	//  @param mobile: 手机号码
	//  @param tpc: template code 消息模板编号
	//  @param tp: template param消息参数
	//  @return error
	//
	SendSms(ctx context.Context, mobile string, tpc, tp string) error
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

func (s *smsService) SendSms(ctx context.Context, mobile string, tpc, tp string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", s.smsOpts.APIKey, s.smsOpts.APISecret)
	if err != nil {
		panic(err)
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = mobile //手机号
	request.QueryParams["SignName"] = "慕学在线"     //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = tpc    //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = tp    //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
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
