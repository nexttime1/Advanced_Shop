package db

import (
	"Advanced_Shop/app/pkg/code"
	code2 "Advanced_Shop/gnova/code"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/errors"
	"context"

	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
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

func (g *goods) CreateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	tx := txn.Create(goods)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	tx := txn.Save(goods)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) DeleteInTxn(ctx context.Context, txn *gorm.DB, ID uint64) error {
	return txn.Where("id = ?", ID).Delete(&do.GoodsDO{}).Error
}

//func NewGoods(db *gorm.DB) *goods {
//	return &goods{
//		db: db,
//	}
//}

func (g *goods) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.GoodsDOList, error) {
	//实现gorm查询
	ret := &do.GoodsDOList{}

	//分页
	var limit, offset int
	if opts.PageSize == 0 {
		limit = 10
	} else {
		limit = opts.PageSize
	}

	if opts.Page > 0 {
		offset = (opts.Page - 1) * limit
	}

	//排序
	query := g.db.Preload("Category").Preload("Brands")
	for _, value := range orderby {
		//坑
		query = query.Order(value)
	}

	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

func (g *goods) Get(ctx context.Context, ID uint64) (*do.GoodsDO, error) {
	good := &do.GoodsDO{}
	err := g.db.Preload("Category").Preload("Brands").First(good, ID).Error
	if err != nil {
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
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

func (g *goods) Create(ctx context.Context, goods *do.GoodsDO) error {
	tx := g.db.Create(goods)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Update(ctx context.Context, goods *do.GoodsDO) error {
	tx := g.db.Save(goods)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Delete(ctx context.Context, ID uint64) error {
	return g.db.Where("id = ?", ID).Delete(&do.GoodsDO{}).Error
}

var _ v1.GoodsStore = &goods{}
