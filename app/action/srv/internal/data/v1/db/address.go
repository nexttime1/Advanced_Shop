package db

import (
	v1 "Advanced_Shop/app/action/srv/internal/data/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"gorm.io/gorm"
)

type addressData struct {
	db *gorm.DB
}

func newAddress(factory *mysqlFactory) *addressData {
	Address := &addressData{
		db: factory.db,
	}
	return Address
}

// GetByID 根据ID和用户ID获取地址信息
func (s *addressData) GetByID(ctx context.Context, ID uint, userID int32) (*do.AddressDO, error) {
	var address do.AddressDO
	err := s.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", ID, userID).
		Take(&address).Error

	if err != nil {
		log.Errorf("GetByID err:%v", err)
		return nil, errors.WithCode(code.ErrDatabase, err.Error())

	}
	return &address, nil
}

// ListByUserID 根据用户ID获取地址列表
func (s *addressData) ListByUserID(ctx context.Context, userID int32) ([]*do.AddressDO, error) {
	var addresses []*do.AddressDO
	err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&addresses).Error

	if err != nil {
		log.Errorf("ListByUserID err:%v", err)
		return nil, errors.WithCode(code.ErrDatabase, err.Error())

	}
	return addresses, nil
}

// Create 创建新地址
func (s *addressData) Create(ctx context.Context, address *do.AddressDO) error {
	err := s.db.WithContext(ctx).Create(address).Error
	if err != nil {
		log.Errorf("Create err:%v", err)
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

// Update 更新地址信息
func (s *addressData) Update(ctx context.Context, address *do.AddressDO) error {
	err := s.db.WithContext(ctx).Save(address).Error
	if err != nil {
		log.Errorf("Update err:%v", err)
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

// Delete 根据ID和用户ID删除地址
func (s *addressData) Delete(ctx context.Context, ID uint, userID int32) error {
	// 先查询确认记录存在
	var address do.AddressDO
	err := s.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", ID, userID).
		Take(&address).Error

	if err != nil {
		log.Errorf("Delete err:%v", err)
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	// 执行删除
	err = s.db.WithContext(ctx).Delete(&address).Error
	if err != nil {
		log.Errorf("Delete err:%v", err)
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

var _ v1.AddressStore = &addressData{}
