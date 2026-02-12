package db

import (
	"Advanced_Shop/app/goods/srv/internal/domain/service"
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/app/pkg/struct_to_map"
	code2 "Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"context"
	"go.uber.org/zap"
	"strings"

	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/log"
	"gorm.io/gorm"
)

type brands struct {
	db *gorm.DB
}

func newBrands(factory *mysqlFactory) *brands {
	brands := &brands{
		db: factory.db,
	}
	return brands
}

func (b *brands) Get(ctx context.Context, ID uint64) (*do.BrandsDO, error) {
	var brandModel do.BrandsDO
	// 根据ID查询品牌
	err := b.db.Take(&brandModel, ID).Error
	if err != nil {
		log.Errorf("brand not found, ID: %d, error: %v", ID, err)
		// 注意：需确保code.ErrBrandNotFound错误码已定义（替换为你实际的错误码）
		return nil, errors.WithCode(code.ErrBrandNotFound, err.Error())
	}
	return &brandModel, nil
}

func (b *brands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error) {
	var brandModels []*do.BrandsDO

	limit := opts.GetLimit()
	offset := opts.GetOffset()

	query := b.db.Model(&do.BrandsDO{})
	var total int64

	if err := query.Count(&total).Error; err != nil {
		log.Errorf("mysql query brand count error: %v", err)
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	// 处理排序
	if len(orderby) > 0 {
		query = query.Order(strings.Join(orderby, ","))
	}

	// 分页查询列表数据
	if err := query.Limit(limit).Offset(offset).Find(&brandModels).Error; err != nil {
		log.Errorf("mysql query brand list error: %v", err)
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	// 构造返回结果
	return &do.BrandsDOList{
		Items:      brandModels,
		TotalCount: total,
	}, nil
}

func (b *brands) Create(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	// 优先使用传入的事务db，无事务则用默认db
	db := b.db
	if txn != nil {
		db = txn
	}
	// 执行创建操作
	err := db.Create(&brands).Error
	if err != nil {
		log.Errorf("mysql create brand error, brand name: %s, error: %v", brands.Name, err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

func (b *brands) Update(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	// 优先使用传入的事务 db
	db := b.db
	if txn != nil {
		db = txn
	}
	// 先校验品牌是否存在
	var model do.BrandsDO
	err := db.Take(&model, brands.ID).Error
	if err != nil {
		log.Errorf("brand not found, ID: %d, error: %v", brands.ID, err)
		return errors.WithCode(code.ErrBrandNotFound, err.Error())
	}

	updateMap := service.BrandUpdateServiceMap{
		Name: brands.Name,
		Logo: brands.Logo,
	}
	// 结构体转map
	toMap := struct_to_map.StructToMap(updateMap)

	// 执行更新操作
	err = db.Model(&model).Updates(toMap).Error
	if err != nil {
		log.Errorf("mysql update brand error, ID: %d, error: %v", brands.ID, err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

func (b *brands) Delete(ctx context.Context, ID uint64) error {
	var model do.BrandsDO
	// 校验品牌是否存在
	err := b.db.Take(&model, ID).Error
	if err != nil {
		zap.S().Errorf("brand not found, ID: %d, error: %v", ID, err)
		return errors.WithCode(code.ErrBrandNotFound, err.Error())
	}

	// 执行删除操作
	err = b.db.Delete(&model).Error
	if err != nil {
		zap.S().Errorf("mysql delete brand error, ID: %d, error: %v", ID, err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

var _ v1.BrandsStore = &brands{}
