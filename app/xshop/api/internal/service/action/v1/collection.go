package v1

import (
	pb "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/xshop/api/internal/data"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CollectionSrv interface {
	GetFavList(ctx context.Context, in *pb.UserFavRequest) (*pb.UserFavListResponse, error)
	AddUserFav(ctx context.Context, in *pb.UserFavRequest) (*emptypb.Empty, error)
	DeleteUserFav(ctx context.Context, in *pb.UserFavRequest) (*emptypb.Empty, error)
	GetUserFavDetail(ctx context.Context, in *pb.UserFavRequest) (*emptypb.Empty, error)
}

type collectionService struct {
	data data.DataFactory
}

func NewCollectionService(data data.DataFactory) CollectionSrv {
	return &collectionService{
		data,
	}
}

func (cs *collectionService) GetFavList(ctx context.Context, in *pb.UserFavRequest) (*pb.UserFavListResponse, error) {
	return cs.data.Collection().GetFavList(ctx, in)
}

func (cs *collectionService) AddUserFav(ctx context.Context, in *pb.UserFavRequest) (*emptypb.Empty, error) {
	return cs.data.Collection().AddUserFav(ctx, in)
}

func (cs *collectionService) DeleteUserFav(ctx context.Context, in *pb.UserFavRequest) (*emptypb.Empty, error) {
	return cs.data.Collection().DeleteUserFav(ctx, in)
}

func (cs *collectionService) GetUserFavDetail(ctx context.Context, in *pb.UserFavRequest) (*emptypb.Empty, error) {
	return cs.data.Collection().GetUserFavDetail(ctx, in)
}

var _ CollectionSrv = (*collectionService)(nil)
