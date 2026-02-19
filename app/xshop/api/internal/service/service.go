package service

import (
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/app/xshop/api/internal/data"
	v3 "Advanced_Shop/app/xshop/api/internal/service/action/v1"
	v1 "Advanced_Shop/app/xshop/api/internal/service/goods/v1"
	v2 "Advanced_Shop/app/xshop/api/internal/service/inventory/v1"
	v14 "Advanced_Shop/app/xshop/api/internal/service/order/v1"
	v12 "Advanced_Shop/app/xshop/api/internal/service/sms/v1"
	v13 "Advanced_Shop/app/xshop/api/internal/service/user/v1"
)

type ServiceFactory interface {
	Goods() v1.GoodsSrv
	Users() v13.UserSrv
	Sms() v12.SmsSrv
	Order() v14.OrderSrv
	Inventory() v2.InventorySrv
	Address() v3.AddressSrv
	Collection() v3.CollectionSrv
	Message() v3.MessageSrv
}

type service struct {
	data data.DataFactory

	smsOpts *options.SmsOptions

	jwtOpts *options.JwtOptions
}

func (s *service) Address() v3.AddressSrv {
	return v3.NewAddressService(s.data)
}

func (s *service) Collection() v3.CollectionSrv {
	return v3.NewCollectionService(s.data)
}

func (s *service) Message() v3.MessageSrv {
	return v3.NewMessageService(s.data)
}

func (s *service) Inventory() v2.InventorySrv {
	return v2.NewInventoryService(s.data)
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

func (S *service) Order() v14.OrderSrv {
	return v14.NewOrderService(S.data)
}

func NewService(store data.DataFactory, smsOpts *options.SmsOptions, jwtOpts *options.JwtOptions) ServiceFactory {
	return &service{data: store,
		smsOpts: smsOpts,
		jwtOpts: jwtOpts,
	}
}

var _ ServiceFactory = &service{}
