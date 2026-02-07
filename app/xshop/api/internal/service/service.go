package service

import (
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/app/xshop/api/internal/data"
	v1 "Advanced_Shop/app/xshop/api/internal/service/goods/v1"
	v12 "Advanced_Shop/app/xshop/api/internal/service/sms/v1"
	v13 "Advanced_Shop/app/xshop/api/internal/service/user/v1"
)

type ServiceFactory interface {
	Goods() v1.GoodsSrv
	Users() v13.UserSrv
	Sms() v12.SmsSrv
}

type service struct {
	data data.DataFactory

	smsOpts *options.SmsOptions

	jwtOpts *options.JwtOptions
}

func (s *service) Sms() v12.SmsSrv {
	return v12.NewSmsService(s.smsOpts)
}

func (s *service) Goods() v1.GoodsSrv {
	return v1.NewGoods(s.data)
}

func (s *service) Users() v13.UserSrv {
	return v13.NewUserService(s.data, s.jwtOpts)
}

func NewService(store data.DataFactory, smsOpts *options.SmsOptions, jwtOpts *options.JwtOptions) *service {
	return &service{data: store,
		smsOpts: smsOpts,
		jwtOpts: jwtOpts,
	}
}

var _ ServiceFactory = &service{}
