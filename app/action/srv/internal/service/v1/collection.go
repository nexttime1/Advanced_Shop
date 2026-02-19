package v1

import (
	proto "Advanced_Shop/api/goods/v1"
	v1 "Advanced_Shop/app/action/srv/internal/data/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"Advanced_Shop/app/action/srv/internal/domain/dto"
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
)

// CollectionSrv 用户收藏业务逻辑层接口
type CollectionSrv interface {
	// GetFavList 获取用户收藏列表
	GetFavList(ctx context.Context, userID int32, goodID int32) (*dto.CollectionDTOList, error)

	// AddUserFav 添加用户收藏（包含商品存在性校验）
	AddUserFav(ctx context.Context, collectionDTO *dto.CollectionDTO) error

	// DeleteUserFav 删除用户收藏
	DeleteUserFav(ctx context.Context, userID int32, goodID int32) error

	// GetUserFavDetail 检查用户是否收藏了某个商品
	GetUserFavDetail(ctx context.Context, userID int32, goodID int32) error
}

type collectionService struct {
	data v1.DataFactory
}

func newCollection(srv *serviceFactory) CollectionSrv {
	return &collectionService{
		data: srv.data,
	}
}

// GetFavList 获取用户收藏列表
func (s *collectionService) GetFavList(ctx context.Context, userID int32, goodID int32) (*dto.CollectionDTOList, error) {
	// 调用数据层获取DO列表和总数
	collectionDOs, count, err := s.data.Collection().ListByUserID(ctx, userID, goodID)
	if err != nil {
		log.Errorf("GetFavList failed: %v", err)
		return nil, err
	}

	// DO转换为DTO
	dtoList := &dto.CollectionDTOList{
		TotalCount: int(count),
		Items:      make([]*dto.CollectionDTO, 0, len(collectionDOs)),
	}

	for _, doItem := range collectionDOs {
		dtoItem := &dto.CollectionDTO{
			UserCollectionDO: *doItem,
		}
		dtoList.Items = append(dtoList.Items, dtoItem)
	}

	return dtoList, nil
}

// AddUserFav 添加用户收藏（包含商品存在性校验）
func (s *collectionService) AddUserFav(ctx context.Context, collectionDTO *dto.CollectionDTO) error {

	_, err := s.data.Goods().GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: collectionDTO.GoodId,
	})
	if err != nil {
		log.Errorf("goods not found: good_id=%d, err=%v", collectionDTO.GoodId, err)
		return errors.WithCode(code.ErrGoodsNotFound, err.Error())
	}

	// 2. 转换DTO到DO
	collectionDO := &do.UserCollectionDO{
		UserId: collectionDTO.UserId,
		GoodId: collectionDTO.GoodId,
	}

	// 3. 调用数据层创建收藏
	err = s.data.Collection().Create(ctx, collectionDO)
	if err != nil {
		log.Errorf("Create collection failed: %v", err)
		return err
	}

	return nil
}

// DeleteUserFav 删除用户收藏
func (s *collectionService) DeleteUserFav(ctx context.Context, userID int32, goodID int32) error {
	// 调用数据层删除
	rowsAffected, err := s.data.Collection().Delete(ctx, userID, goodID)
	if err != nil {
		log.Errorf("Delete collection failed: %v", err)
		return err
	}

	// 检查是否删除成功
	if rowsAffected == 0 {
		log.Errorf("collection not found: user_id=%d, good_id=%d", userID, goodID)
		return errors.WithCode(code.ErrRecordNotFound, "收藏记录不存在")
	}

	return nil
}

// GetUserFavDetail 检查用户是否收藏了某个商品
func (s *collectionService) GetUserFavDetail(ctx context.Context, userID int32, goodID int32) error {
	// 调用数据层查询
	_, err := s.data.Collection().GetByUserAndGoodID(ctx, userID, goodID)
	if err != nil {
		log.Errorf("GetUserFavDetail failed: user_id=%d, good_id=%d, err=%v", userID, goodID, err)
		return err
	}

	return nil
}

// 确保实现了接口
var _ CollectionSrv = &collectionService{}
