package do

import (
	bgorm "Advanced_Shop/app/pkg/gorm"
)

type GoodsCategoryBrandDO struct {
	bgorm.Model

	CategoryID int32 `gorm:"type:int;not null;comment:分类ID（逻辑外键，关联category_models.id）;index:idx_category_brand,unique"`
	//禁用物理外键约束
	Category *CategoryDO `gorm:"foreignKey:CategoryID;references:ID;constraint:<-:false,foreignKey:no action"`
	BrandsID int32       `gorm:"type:int;not null;comment:品牌ID（逻辑外键，关联brands.id）;index:idx_category_brand,unique"`
	// 禁用物理外键约束
	Brands *BrandsDO `gorm:"foreignKey:BrandsID;references:ID;constraint:<-:false,foreignKey:no action"`
}

func (GoodsCategoryBrandDO) TableName() string {
	return "brand_category_models"
}

type GoodsCategoryBrandList struct {
	TotalCount int64                   `json:"totalCount,omitempty"`
	Items      []*GoodsCategoryBrandDO `json:"items"`
}
