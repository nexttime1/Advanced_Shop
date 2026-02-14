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
	"gorm.io/gorm"
)

type goods struct {
	db *gorm.DB
}

func (g *goods) Begin() *gorm.DB {
	return g.db.Begin()
}

func newGoods(factory *mysqlFactory) *goods {
	goods := &goods{
		db: factory.db,
	}
	return goods
}

func (g *goods) CreateInTxn(ctx context.Context, txn *gorm.DB, goods *v1.GoodsInfo) error {
	// 商品表
	tx := txn.Create(goods.GoodsDO)
	if tx.Error != nil {
		log.Errorf("mysql create goods error: %v", tx.Error)
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}

	// web 已经上传了 七牛云 这里就是 url
	// 添加第三章表  图片
	// 主图
	var ImagesModels []*do.GoodsImageModel

	err := txn.Create(&do.GoodsImageModel{
		GoodsID:   goods.GoodsDO.ID,
		ImageURL:  goods.GoodsFrontImage,
		Sort:      0,
		IsMain:    true,
		ImageType: do.MainImageType, //（1=主图，2=详情图，3=其他）
	}).Error
	if err != nil {
		log.Errorf("mysql create good Images error: %v", err)
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	ImagesModels = append(ImagesModels, &do.GoodsImageModel{
		ImageURL: goods.GoodsFrontImage,
	})

	for i, image := range goods.DescImages {
		err = txn.Debug().Create(&do.GoodsImageModel{
			GoodsID:   goods.GoodsDO.ID,
			ImageURL:  image,
			Sort:      i + 1,
			IsMain:    true,
			ImageType: do.DetailImageType, //（1=主图，2=详情图，3=其他）
		}).Error
		if err != nil {
			log.Errorf("mysql create good Images error: %v", err)
			return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
		}
		ImagesModels = append(ImagesModels, &do.GoodsImageModel{
			ImageURL: image,
		})
	}

	for i, image := range goods.Images {
		err = txn.Create(&do.GoodsImageModel{
			GoodsID:   goods.GoodsDO.ID,
			ImageURL:  image,
			Sort:      i + 1,
			IsMain:    true,
			ImageType: do.OtherImageType, //（1=主图，2=详情图，3=其他）
		}).Error
		if err != nil {
			log.Errorf("mysql create good Images error: %v", err)
			return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
		}

	}

	return nil
}

func (g *goods) UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *v1.GoodsInfo) error {
	var model do.GoodsDO
	err := txn.Where("id = ?", goods.GoodsDO.ID).Take(&model).Error
	if err != nil {
		log.Errorf("mysql good not found error: %v", err)
		return errors.WithCode(code.ErrGoodsNotFound, err.Error())
	}

	if goods.GoodsDO.BrandsID != 0 {
		var brand do.BrandsDO
		err := txn.Where("id = ?", goods.GoodsDO.BrandsID).Take(&brand).Error
		if err != nil {
			log.Errorf("mysql brand not found error: %v", err)
			return errors.WithCode(code.ErrBrandNotFound, err.Error())
		}

	}
	if goods.GoodsDO.CategoryID != 0 {
		var category do.CategoryDO
		err := txn.Where("id = ?", goods.GoodsDO.CategoryID).Take(&category).Error
		if err != nil {
			log.Errorf("mysql category not found error: %v", err)
			return errors.WithCode(code.ErrCategoryNotFound, err.Error())
		}

	}

	// 修改 第三章表
	if goods.GoodsFrontImage != "" {
		err = txn.Where("goods_id = ? and is_main = 1", goods.GoodsDO.ID).Delete(&do.GoodsImageModel{}).Error
		if err != nil {
			log.Errorf("mysql GoodsImage not found error: %v", err)
			return errors.WithCode(code.ErrGoodsImageNotFound, err.Error())
		}
		err = txn.Create(&do.GoodsImageModel{
			GoodsID:   model.ID,
			ImageURL:  goods.GoodsFrontImage,
			Sort:      0,
			IsMain:    true,
			ImageType: 1, //（1=主图，2=详情图，3=其他）
		}).Error
		if err != nil {
			log.Errorf("mysql create error: %v", err)
			return errors.WithCode(code2.ErrDatabase, err.Error())
		}
	}
	if goods.DescImages != nil {
		err = txn.Where("goods_id = ? and image_type = 2", goods.GoodsDO.ID).Delete(&do.GoodsImageModel{}).Error
		if err != nil {
			log.Errorf("mysql GoodsImage not found error: %v", err)
			return errors.WithCode(code.ErrGoodsImageNotFound, err.Error())
		}
		for i, image := range goods.DescImages {
			err = txn.Create(&do.GoodsImageModel{
				GoodsID:   model.ID,
				ImageURL:  image,
				Sort:      i + 1,
				IsMain:    true,
				ImageType: 2, //（1=主图，2=详情图，3=其他）
			}).Error
			if err != nil {
				log.Errorf("mysql create error: %v", err)
				return errors.WithCode(code2.ErrDatabase, err.Error())
			}
		}
	}
	if goods.Images != nil {
		err = txn.Where("goods_id = ? and image_type = 3", goods.GoodsDO.ID).Delete(&do.GoodsImageModel{}).Error
		if err != nil {
			log.Errorf("mysql GoodsImage not found error: %v", err)
			return errors.WithCode(code.ErrGoodsImageNotFound, err.Error())
		}

		for i, image := range goods.Images {
			err = txn.Create(&do.GoodsImageModel{
				GoodsID:   model.ID,
				ImageURL:  image,
				Sort:      i + 1,
				IsMain:    true,
				ImageType: 3, //（1=主图，2=详情图，3=其他）
			}).Error
			if err != nil {
				log.Errorf("mysql create error: %v", err)
				return errors.WithCode(code2.ErrDatabase, err.Error())
			}
		}

	}

	// 修改 商品表
	StructMap := service.GoodUpdateServiceMap{
		Name:        goods.GoodsDO.Name,
		GoodsSn:     goods.GoodsDO.GoodsSn,
		MarketPrice: goods.GoodsDO.MarketPrice,
		ShopPrice:   goods.GoodsDO.ShopPrice,
		GoodsBrief:  goods.GoodsDO.GoodsBrief,
		ShipFree:    goods.GoodsDO.ShipFree,
		IsNew:       goods.GoodsDO.IsNew,
		IsHot:       goods.GoodsDO.IsHot,
		OnSale:      goods.GoodsDO.OnSale,
		CategoryId:  goods.GoodsDO.CategoryID,
		Brand:       goods.GoodsDO.BrandsID,
	}

	toMap := struct_to_map.StructToMap(StructMap)
	err = txn.Model(&model).Updates(toMap).Error
	if err != nil {
		log.Errorf("mysql update error: %v", err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}

	return nil
}

func (g *goods) DeleteInTxn(ctx context.Context, txn *gorm.DB, ID uint64) error {
	err := txn.Where("goods_id = ?", ID).Delete(&do.GoodsImageModel{}).Error
	if err != nil {
		log.Errorf("mysql delete good Images error: %v", err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	err = txn.Where("id = ?", ID).Delete(&do.GoodsDO{}).Error
	if err != nil {
		log.Errorf("mysql delete good  error: %v", err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}

	return nil
}

func (g *goods) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.GoodsDOList, error) {
	//实现gorm查询
	ret := &do.GoodsDOList{}

	//分页
	limit := opts.GetLimit()
	offset := opts.GetOffset()

	query := g.db.Preload("Category").Preload("Brands")
	//排序
	for _, value := range orderby {
		query = query.Order(value)
	}

	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		log.Errorf("mysql query error: %v", d.Error)
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

func (g *goods) Get(ctx context.Context, ID uint64) (*do.GoodsDO, error) {
	good := &do.GoodsDO{}
	err := g.db.Preload("Category").Preload("Brands").First(good, ID).Error
	if err != nil {
		log.Errorf("mysql query error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrGoodsNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return good, nil
}

func (g *goods) ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error) {
	//实现gorm查询
	ret := &do.GoodsDOList{}

	//排序
	query := g.db.Preload("Category").Preload("Brands")
	for _, value := range orderby {
		query = query.Order(value)
	}

	d := query.Where("id in ?", ids).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		log.Errorf("mysql query error: %v", d.Error)
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

func (g *goods) Create(ctx context.Context, goods *v1.GoodsInfo) error {
	tx := g.db.Create(goods)
	if tx.Error != nil {
		log.Errorf("mysql create error: %v", tx.Error)
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Update(ctx context.Context, goods *v1.GoodsInfo) error {
	tx := g.db.Save(goods)
	if tx.Error != nil {
		log.Errorf("mysql update error: %v", tx.Error)
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Delete(ctx context.Context, ID uint64) error {
	err := g.db.Where("id = ?", ID).Delete(&do.GoodsDO{}).Error
	if err != nil {
		log.Errorf("mysql delete error: %v", err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return err
}

var _ v1.GoodsStore = &goods{}
