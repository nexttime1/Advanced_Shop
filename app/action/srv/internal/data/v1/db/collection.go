package db

import (
	v1 "Advanced_Shop/app/action/srv/internal/data/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	code2 "Advanced_Shop/app/pkg/code"
	"Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"gorm.io/gorm"
)

type collectionData struct {
	db *gorm.DB
}

func newCollection(factory *mysqlFactory) v1.CollectionStore {
	collection := &collectionData{
		db: factory.db,
	}
	return collection
}

// ListByUserID 根据用户ID获取收藏列表（支持商品ID筛选）
func (s *collectionData) ListByUserID(ctx context.Context, userID int32, goodID int32) ([]*do.UserCollectionDO, int64, error) {
	var collections []*do.UserCollectionDO
	tx := s.db.WithContext(ctx).Model(&do.UserCollectionDO{}).Where("user_id = ?", userID)

	// 商品ID筛选
	if goodID > 0 {
		tx = tx.Where("good_id = ?", goodID)
	}

	// 统计总数
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		log.Errorf("ListByUserID count err:%v", err)
		return nil, 0, errors.WithCode(code.ErrDatabase, err.Error())
	}

	// 查询列表
	if err := tx.Find(&collections).Error; err != nil {
		log.Errorf("ListByUserID find err:%v", err)
		return nil, 0, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return collections, count, nil
}

// Create 创建用户收藏
func (s *collectionData) Create(ctx context.Context, collection *do.UserCollectionDO) error {
	err := s.db.WithContext(ctx).Create(collection).Error
	if err != nil {
		log.Errorf("Create collection err:%v", err)
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

// Delete 根据用户ID和商品ID删除收藏
func (s *collectionData) Delete(ctx context.Context, userID int32, goodID int32) (int64, error) {
	result := s.db.WithContext(ctx).Unscoped().
		Where("user_id = ? AND good_id = ?", userID, goodID).
		Delete(&do.UserCollectionDO{})

	if result.Error != nil {
		log.Errorf("Delete collection err:%v", result.Error)
		return 0, errors.WithCode(code.ErrDatabase, result.Error.Error())
	}

	return result.RowsAffected, nil
}

// GetByUserAndGoodID 检查用户是否收藏了某个商品
func (s *collectionData) GetByUserAndGoodID(ctx context.Context, userID int32, goodID int32) (*do.UserCollectionDO, error) {
	var collection do.UserCollectionDO
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND good_id = ?", userID, goodID).
		Take(&collection).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Errorf("GetByUserAndGoodID not found: user_id=%d, good_id=%d", userID, goodID)
			return nil, errors.WithCode(code2.ErrRecordNotFound, err.Error())
		}
		log.Errorf("GetByUserAndGoodID err:%v", err)
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return &collection, nil
}

var _ v1.CollectionStore = &collectionData{}
