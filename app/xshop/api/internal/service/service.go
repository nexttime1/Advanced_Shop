package service

import (
	"mxshop/app/mxshop/api/internal/data"
	v1 "mxshop/app/mxshop/api/internal/service/goods/v1"
	v12 "mxshop/app/mxshop/api/internal/service/sms/v1"
	v13 "mxshop/app/mxshop/api/internal/service/user/v1"
	"mxshop/app/pkg/options"
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
