package db

import (
	"Advanced_Shop/app/pkg/code"
	code2 "Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"context"

	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"gorm.io/gorm"
)

type categorys struct {
	db *gorm.DB
}

func newCategories(factory *mysqlFactory) *categorys {
	categories := &categorys{
		db: factory.db,
	}
	return categories
}

func (c *categorys) Get(ctx context.Context, ID uint64) (*do.CategoryDO, error) {
	category := &do.CategoryDO{}

	err := c.db.Preload("SubCategory").Preload("SubCategory.SubCategory").First(category, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrCategoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return category, nil
}

func (c *categorys) ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error) {
	ret := &do.CategoryDOList{}
	query := c.db

	for _, value := range orderby {
		query = query.Order(value)
	}

	d := query.Where("level = 1").Preload("SubCategory.SubCategory").Find(&ret.Items)
	return ret, d.Error
}

func (c *categorys) Create(ctx context.Context, category *do.CategoryDO) error {
	tx := c.db.Create(category)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (c *categorys) Update(ctx context.Context, category *do.CategoryDO) error {
	tx := c.db.Save(category)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (c *categorys) Delete(ctx context.Context, ID uint64) error {
	return c.db.Where("id = ?", ID).Delete(&do.GoodsDO{}).Error
}

var _ v1.CategoryStore = &categorys{}
