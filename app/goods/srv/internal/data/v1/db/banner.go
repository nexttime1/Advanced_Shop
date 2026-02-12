package db

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"Advanced_Shop/app/goods/srv/internal/domain/service"
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/app/pkg/struct_to_map"
	code2 "Advanced_Shop/gnova/code"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type banners struct {
	db *gorm.DB
}

func newBanner(factory *mysqlFactory) *banners {
	banners := &banners{
		db: factory.db,
	}
	return banners
}

func (b *banners) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BannerList, error) {
	var bannerModels []*do.BannerDO
	limit := opts.GetLimit()
	offset := opts.GetOffset()
	query := b.db.Model(&do.BannerDO{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Errorf("mysql query error: %v", err)
		return nil, errors.WithCode(code2.ErrDatabase, err.Error()) // 统计总数出错
	}
	// 处理排序
	if len(orderby) > 0 {
		query = query.Order(strings.Join(orderby, ","))
	}

	if err := query.Limit(limit).Offset(offset).Find(&bannerModels).Error; err != nil {
		log.Errorf("mysql query error: %v", err)
		return nil, errors.WithCode(code2.ErrDatabase, err.Error()) // 查找错误
	}
	return &do.BannerList{
		Items:      bannerModels,
		TotalCount: total,
	}, nil
}

func (b *banners) Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	err := b.db.Create(&banner).Error
	if err != nil {
		log.Errorf("mysql create error: %v", err)
		return errors.WithCode(code2.ErrDatabase, err.Error()) // 查找错误
	}
	return nil
}

func (b *banners) Update(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	var model do.BannerDO
	err := b.db.Take(&model, banner.ID).Error
	if err != nil {
		log.Errorf("banner not found : %v", err)
		return errors.WithCode(code.ErrBannerNotFound, err.Error()) // 找不到
	}
	updateMap := service.BannerUpdateServiceMap{
		Image: banner.Image,
		Url:   banner.Url,
		Index: banner.Index,
	}
	toMap := struct_to_map.StructToMap(updateMap)
	err = b.db.Model(&model).Updates(toMap).Error
	if err != nil {
		log.Errorf("mysql create error: %v", err)
		return errors.WithCode(code2.ErrDatabase, err.Error()) // 查找错误
	}
	return nil
}

func (b *banners) Delete(ctx context.Context, ID uint64) error {
	var model do.BannerDO
	err := b.db.Take(&model, ID).Error
	if err != nil {
		zap.S().Error(err.Error())
		return errors.WithCode(code.ErrBannerNotFound, err.Error()) // 找不到
	}
	err = b.db.Delete(&model).Error
	if err != nil {
		zap.S().Error(err.Error())
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

var _ v1.BannerStore = &banners{}
