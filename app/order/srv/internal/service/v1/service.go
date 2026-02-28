package service

import (
	v1 "Advanced_Shop/app/order/srv/internal/data/v1"
	"Advanced_Shop/app/pkg/options"
)

type ServiceFactory interface {
	Orders() OrderSrv
	Cart() CartSrv
}

type service struct {
	data    v1.DataFactory
	dtmopts *options.DtmOptions
	MqOpts  *options.RocketMQOptions
}

func (s *service) Cart() CartSrv {
	return NewCartService(s)
}

func (s *service) Orders() OrderSrv {
	return newOrderService(s)
}

var _ ServiceFactory = &service{}

func NewService(data v1.DataFactory, dtmopts *options.DtmOptions, mqOpts *options.RocketMQOptions) ServiceFactory {
	return &service{data: data, dtmopts: dtmopts, MqOpts: mqOpts}
}
