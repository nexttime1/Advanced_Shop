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

type categoryBrands struct {
	db *gorm.DB
}

func (cb *categoryBrands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandList, error) {
	var gcbModels []*do.GoodsCategoryBrandDO

	// 获取分页参数
	limit := opts.GetLimit()
	offset := opts.GetOffset()

	// 构建基础查询器
	query := cb.db.Model(&do.GoodsCategoryBrandDO{})
	var total int64

	// 统计符合条件的总记录数
	if err := query.Count(&total).Error; err != nil {
		log.Errorf("mysql query category-brand count error: %v", err)
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	// 处理排序
	if len(orderby) > 0 {
		query = query.Order(strings.Join(orderby, ","))
	}

	// 第三步：执行分页查询
	if err := query.Limit(limit).Offset(offset).Find(&gcbModels).Error; err != nil {
		log.Errorf("mysql query category-brand list error: %v", err)
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	// 构造并返回列表结果
	return &do.GoodsCategoryBrandList{
		TotalCount: total,
		Items:      gcbModels,
	}, nil
}

func (cb *categoryBrands) Create(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error {
	db := cb.db
	if txn != nil {
		db = txn
	}
	// 执行创建操作（注意唯一索引idx_category_brand会约束分类+品牌组合唯一）
	err := db.Create(&gcb).Error
	if err != nil {
		log.Errorf("mysql create category-brand error, CategoryID: %d, BrandsID: %d, error: %v",
			gcb.CategoryID, gcb.BrandsID, err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

func (cb *categoryBrands) Update(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error {
	// 优先使用传入的事务DB
	db := cb.db
	if txn != nil {
		db = txn
	}
	// 校验分类-品牌关联记录是否存在
	var model do.GoodsCategoryBrandDO
	err := db.Take(&model, gcb.ID).Error
	if err != nil {
		log.Errorf("category-brand not found, ID: %d, error: %v", gcb.ID, err)
		// 需确保code.ErrCategoryBrandNotFound错误码已定义（替换为你实际的错误码）
		return errors.WithCode(code.ErrCategoryBrandNotFound, err.Error())
	}

	// 构造需要更新的字段（
	updateMap := service.CategoryBrandUpdateServiceMap{
		CategoryID: gcb.CategoryID,
		BrandsID:   gcb.BrandsID,
	}

	// 结构体转map
	toMap := struct_to_map.StructToMap(updateMap)

	// 执行更新操作
	err = db.Model(&model).Updates(toMap).Error
	if err != nil {
		log.Errorf("mysql update category-brand error, ID: %d, CategoryID: %d, BrandsID: %d, error: %v",
			gcb.ID, gcb.CategoryID, gcb.BrandsID, err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

func (cb *categoryBrands) Delete(ctx context.Context, ID uint64) error {
	// 校验记录是否存在
	var model do.GoodsCategoryBrandDO
	err := cb.db.Take(&model, ID).Error
	if err != nil {
		zap.S().Errorf("category-brand not found, ID: %d, error: %v", ID, err)
		return errors.WithCode(code.ErrCategoryBrandNotFound, err.Error())
	}

	// 执行删除操作
	err = cb.db.Delete(&model).Error
	if err != nil {
		zap.S().Errorf("mysql delete category-brand error, ID: %d, error: %v", ID, err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

var _ v1.GoodsCategoryBrandStore = &categoryBrands{}
