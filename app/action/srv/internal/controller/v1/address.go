package v1

import (
	pb "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"Advanced_Shop/app/action/srv/internal/domain/dto"
	v1 "Advanced_Shop/app/action/srv/internal/service/v1"
	gorm2 "Advanced_Shop/app/pkg/gorm"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type actionServer struct {
	pb.UnimplementedUserFavServer
	pb.UnimplementedAddressServer
	pb.UnimplementedMessageServer
	srv v1.ServiceFactory
}

func NewGoodsServer(srv v1.ServiceFactory) *actionServer {
	return &actionServer{srv: srv}
}

// GetAddressList 获取地址列表
func (o *actionServer) GetAddressList(ctx context.Context, request *pb.AddressRequest) (*pb.AddressListResponse, error) {
	// 调用业务层
	dtoList, err := o.srv.Address().GetAddressList(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	// DTO转换为Proto响应
	response := &pb.AddressListResponse{
		Total: int32(dtoList.TotalCount),
		Data:  make([]*pb.AddressResponse, 0, len(dtoList.Items)),
	}

	for _, dtoItem := range dtoList.Items {
		response.Data = append(response.Data, &pb.AddressResponse{
			Id:           dtoItem.ID,
			UserId:       dtoItem.UserId,
			Province:     dtoItem.Province,
			City:         dtoItem.City,
			District:     dtoItem.District,
			Address:      dtoItem.Address,
			SignerName:   dtoItem.SignerName,
			SignerMobile: dtoItem.SignerMobile,
		})
	}

	return response, nil
}

// CreateAddress 创建地址
func (o *actionServer) CreateAddress(ctx context.Context, request *pb.AddressRequest) (*pb.AddressResponse, error) {
	// Proto转换为DTO
	addressDTO := &dto.AddressDTO{
		AddressDO: do.AddressDO{
			UserId:       request.UserId,
			Province:     request.Province,
			City:         request.City,
			District:     request.District,
			Address:      request.Address,
			SignerName:   request.SignerName,
			SignerMobile: request.SignerMobile,
		},
	}

	// 调用业务层
	createdDTO, err := o.srv.Address().CreateAddress(ctx, addressDTO)
	if err != nil {
		return nil, err
	}

	// DTO转换为Proto响应
	return &pb.AddressResponse{
		Id: createdDTO.ID,
	}, nil
}

// DeleteAddress 删除地址
func (o *actionServer) DeleteAddress(ctx context.Context, request *pb.AddressRequest) (*emptypb.Empty, error) {
	err := o.srv.Address().DeleteAddress(ctx, uint(request.Id), request.UserId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateAddress 更新地址
func (o *actionServer) UpdateAddress(ctx context.Context, request *pb.AddressRequest) (*emptypb.Empty, error) {
	// Proto转换为DTO
	addressDTO := &dto.AddressDTO{
		AddressDO: do.AddressDO{
			Model:        gorm2.Model{ID: request.Id},
			UserId:       request.UserId,
			Province:     request.Province,
			City:         request.City,
			District:     request.District,
			Address:      request.Address,
			SignerName:   request.SignerName,
			SignerMobile: request.SignerMobile,
		},
	}

	// 调用业务层
	err := o.srv.Address().UpdateAddress(ctx, addressDTO)
	if err != nil {
		return nil, err

	}

	return &emptypb.Empty{}, nil
}

var _ pb.AddressServer = &actionServer{}
