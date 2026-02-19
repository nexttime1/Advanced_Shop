package v1

import v1 "Advanced_Shop/app/action/srv/internal/data/v1"

type ServiceFactory interface {
	Address() AddressSrv
	Collection() CollectionSrv
	Message() MessageSrv
}

type serviceFactory struct {
	data v1.DataFactory
}

func NewService(store v1.DataFactory) ServiceFactory {
	return &serviceFactory{data: store}
}

var _ ServiceFactory = &serviceFactory{}

func (s *serviceFactory) Address() AddressSrv {
	return newAddress(s)
}

func (s *serviceFactory) Collection() CollectionSrv {
	return newCollection(s)
}

func (s *serviceFactory) Message() MessageSrv {
	return newMessage(s)
}
