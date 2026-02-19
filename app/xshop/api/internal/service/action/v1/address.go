package v1

import (
	pb "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/xshop/api/internal/data"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AddressSrv interface {
	GetAddressList(ctx context.Context, in *pb.AddressRequest) (*pb.AddressListResponse, error)
	CreateAddress(ctx context.Context, in *pb.AddressRequest) (*pb.AddressResponse, error)
	DeleteAddress(ctx context.Context, in *pb.AddressRequest) (*emptypb.Empty, error)
	UpdateAddress(ctx context.Context, in *pb.AddressRequest) (*emptypb.Empty, error)
}

type addressService struct {
	data data.DataFactory
}

func NewAddressService(data data.DataFactory) AddressSrv {
	return &addressService{
		data: data,
	}
}

func (as *addressService) GetAddressList(ctx context.Context, in *pb.AddressRequest) (*pb.AddressListResponse, error) {
	return as.data.Address().GetAddressList(ctx, in)
}

func (as *addressService) CreateAddress(ctx context.Context, in *pb.AddressRequest) (*pb.AddressResponse, error) {
	return as.data.Address().CreateAddress(ctx, in)
}

func (as *addressService) DeleteAddress(ctx context.Context, in *pb.AddressRequest) (*emptypb.Empty, error) {
	return as.data.Address().DeleteAddress(ctx, in)
}

func (as *addressService) UpdateAddress(ctx context.Context, in *pb.AddressRequest) (*emptypb.Empty, error) {
	return as.data.Address().UpdateAddress(ctx, in)
}

var _ AddressSrv = (*addressService)(nil)
