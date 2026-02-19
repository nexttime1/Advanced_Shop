package v1

import (
	v1 "Advanced_Shop/app/action/srv/internal/data/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"Advanced_Shop/app/action/srv/internal/domain/dto"
	gorm2 "Advanced_Shop/app/pkg/gorm"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AddressSrv 地址业务逻辑层接口
type AddressSrv interface {
	// GetAddressList 根据用户ID获取地址列表
	GetAddressList(ctx context.Context, userID int32) (*dto.AddressDTOList, error)

	// CreateAddress 创建新地址
	CreateAddress(ctx context.Context, addressDTO *dto.AddressDTO) (*dto.AddressDTO, error)

	// UpdateAddress 更新地址信息
	UpdateAddress(ctx context.Context, addressDTO *dto.AddressDTO) error

	// DeleteAddress 删除地址
	DeleteAddress(ctx context.Context, ID uint, userID int32) error

	// GetAddressByID 根据ID和用户ID获取地址详情
	GetAddressByID(ctx context.Context, ID uint, userID int32) (*dto.AddressDTO, error)
}

type addressService struct {
	//工厂
	data v1.DataFactory
}

func newAddress(srv *serviceFactory) AddressSrv {
	return &addressService{
		data: srv.data,
	}
}

// GetAddressList 根据用户ID获取地址列表
func (s *addressService) GetAddressList(ctx context.Context, userID int32) (*dto.AddressDTOList, error) {

	// 调用数据层获取DO列表
	addressDOs, err := s.data.Address().ListByUserID(ctx, userID)
	if err != nil {
		log.Errorf("获取地址列表失败: %v", err)
		return nil, err
	}

	// DO转换为DTO
	dtoList := &dto.AddressDTOList{
		TotalCount: len(addressDOs),
		Items:      make([]*dto.AddressDTO, 0, len(addressDOs)),
	}

	for _, doItem := range addressDOs {
		dtoItem := &dto.AddressDTO{
			AddressDO: *doItem,
		}
		dtoList.Items = append(dtoList.Items, dtoItem)
	}

	return dtoList, nil
}

// CreateAddress 创建新地址
func (s *addressService) CreateAddress(ctx context.Context, addressDTO *dto.AddressDTO) (*dto.AddressDTO, error) {

	// DTO转换为DO
	addressDO := &do.AddressDO{
		UserId:       addressDTO.UserId,
		Province:     addressDTO.Province,
		City:         addressDTO.City,
		District:     addressDTO.District,
		Address:      addressDTO.Address,
		SignerName:   addressDTO.SignerName,
		SignerMobile: addressDTO.SignerMobile,
	}

	// 调用数据层创建
	err := s.data.Address().Create(ctx, addressDO)
	if err != nil {
		log.Errorf("创建地址失败: %v", err)
		return nil, err
	}

	// 设置创建后的ID并返回
	addressDTO.ID = addressDO.ID
	return addressDTO, nil
}

// UpdateAddress 更新地址信息
func (s *addressService) UpdateAddress(ctx context.Context, addressDTO *dto.AddressDTO) error {

	// 先查询确认地址存在
	_, err := s.data.Address().GetByID(ctx, uint(addressDTO.ID), addressDTO.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.S().Errorf("地址不存在: id=%d, user_id=%d", addressDTO.ID, addressDTO.UserId)
			return err
		}
		zap.S().Errorf("查询地址失败: %v", err)
		return err
	}

	// DTO转换为DO
	addressDO := &do.AddressDO{
		Model:        gorm2.Model{ID: addressDTO.ID},
		UserId:       addressDTO.UserId,
		Province:     addressDTO.Province,
		City:         addressDTO.City,
		District:     addressDTO.District,
		Address:      addressDTO.Address,
		SignerName:   addressDTO.SignerName,
		SignerMobile: addressDTO.SignerMobile,
	}

	// 调用数据层更新
	err = s.data.Address().Update(ctx, addressDO)
	if err != nil {
		log.Errorf("更新地址失败: %v", err)
		return err
	}

	return nil
}

// DeleteAddress 删除地址
func (s *addressService) DeleteAddress(ctx context.Context, ID uint, userID int32) error {

	// 调用数据层删除
	err := s.data.Address().Delete(ctx, ID, userID)
	if err != nil {
		log.Errorf("删除地址失败: %v", err)
		return err
	}

	return nil
}

// GetAddressByID 根据ID和用户ID获取地址详情
func (s *addressService) GetAddressByID(ctx context.Context, ID uint, userID int32) (*dto.AddressDTO, error) {

	// 调用数据层获取DO
	addressDO, err := s.data.Address().GetByID(ctx, ID, userID)
	if err != nil {
		log.Errorf("获取地址详情失败: %v", err)
		return nil, err
	}

	// DO转换为DTO
	addressDTO := &dto.AddressDTO{
		AddressDO: *addressDO,
	}

	return addressDTO, nil
}

var _ AddressSrv = &addressService{}
