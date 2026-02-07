package aliyun

import (
	"Advanced_Shop/app/pkg/options"
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dypnsapi20170525 "github.com/alibabacloud-go/dypnsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
	"strings"
)

func createClient(opts *options.SmsOptions) (_result *dypnsapi20170525.Client, _err error) {
	// 设置 AccessKeyId 和 AccessKeySecret
	config := &openapi.Config{
		AccessKeyId:     tea.String(opts.APIKey),    // AccessKeyID
		AccessKeySecret: tea.String(opts.APISecret), // AccessKeySecret
		Endpoint:        tea.String("dypnsapi.aliyuncs.com"),
	}

	_result, _err = dypnsapi20170525.NewClient(config)
	return _result, _err
}

func SendCode(phone string, opts *options.SmsOptions) (_err error) {
	client, _err := createClient(opts)
	if _err != nil {
		return _err
	}

	sendSmsVerifyCodeRequest := &dypnsapi20170525.SendSmsVerifyCodeRequest{
		SignName:      tea.String("速通互联验证码"),
		TemplateCode:  tea.String("100001"),
		PhoneNumber:   tea.String(phone),
		TemplateParam: tea.String("{\"code\":\"##code##\",\"min\":\"5\"}"),
	}
	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		resp, _err := client.SendSmsVerifyCodeWithOptions(sendSmsVerifyCodeRequest, runtime)
		if _err != nil {
			return _err
		}
		// TODO
		zap.S().Infof("短信发送成功，响应结果：%#v\n", util.ToJSONString(resp))
		return nil

	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		zap.S().Errorf("短信发送失败，错误信息：", tea.StringValue(error.Message))
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			zap.S().Infof("阿里云建议： %#v", recommend)
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

// CheckSmsVerifyCode 验证验证码
func CheckSmsVerifyCode(phoneNumber, verifyCode string, opts *options.SmsOptions) (success bool, _err error) {
	client, _err := createClient(opts)
	if _err != nil {
		return false, _err
	}

	checkSmsVerifyCodeRequest := &dypnsapi20170525.CheckSmsVerifyCodeRequest{
		PhoneNumber: tea.String(phoneNumber),
		VerifyCode:  tea.String(verifyCode),
	}
	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		resp, _err := client.CheckSmsVerifyCodeWithOptions(checkSmsVerifyCodeRequest, runtime)
		if _err != nil {
			return _err
		}

		// 检查验证结果
		if resp.Body != nil {
			// 根据阿里云文档，验证成功会返回 200 状态码和相应的数据
			// 你可以根据实际返回结构来判断验证是否成功
			zap.S().Infof("验证码验证响应：%v\n", util.ToJSONString(resp))

			// 通常验证成功会返回 Model 数据
			if resp.Body.Model != nil {
				success = true
			}
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		zap.S().Errorf("验证码验证失败，错误信息：", tea.StringValue(error.Message))
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			zap.S().Infof("阿里云建议：", recommend)
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return false, _err
		}
	}
	return success, _err
}
