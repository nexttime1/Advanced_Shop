package v1

import (
	pb "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"Advanced_Shop/app/action/srv/internal/domain/dto"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetFavList 获取用户收藏列表
func (o *actionServer) GetFavList(ctx context.Context, request *pb.UserFavRequest) (*pb.UserFavListResponse, error) {
	// 调用业务层
	dtoList, err := o.srv.Collection().GetFavList(ctx, request.UserId, request.GoodsId)
	if err != nil {
		return nil, err
	}

	// DTO转换为Proto响应
	response := &pb.UserFavListResponse{
		Total: int32(dtoList.TotalCount),
		Data:  make([]*pb.UserFavResponse, 0, len(dtoList.Items)),
	}

	for _, dtoItem := range dtoList.Items {
		response.Data = append(response.Data, &pb.UserFavResponse{
			UserId:  dtoItem.UserId,
			GoodsId: dtoItem.GoodId,
		})
	}

	return response, nil
}

// AddUserFav 添加用户收藏
func (o *actionServer) AddUserFav(ctx context.Context, request *pb.UserFavRequest) (*emptypb.Empty, error) {
	// Proto转换为DTO
	collectionDTO := &dto.CollectionDTO{
		UserCollectionDO: do.UserCollectionDO{
			UserId: request.UserId,
			GoodId: request.GoodsId,
		},
	}

	// 调用业务层
	err := o.srv.Collection().AddUserFav(ctx, collectionDTO)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteUserFav 删除用户收藏
func (o *actionServer) DeleteUserFav(ctx context.Context, request *pb.UserFavRequest) (*emptypb.Empty, error) {
	err := o.srv.Collection().DeleteUserFav(ctx, request.UserId, request.GoodsId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetUserFavDetail 检查用户收藏详情
func (o *actionServer) GetUserFavDetail(ctx context.Context, request *pb.UserFavRequest) (*emptypb.Empty, error) {
	err := o.srv.Collection().GetUserFavDetail(ctx, request.UserId, request.GoodsId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
