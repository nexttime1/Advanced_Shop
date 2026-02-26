package do

import (
	"Advanced_Shop/app/pkg/gorm"
)

type CategoryDO struct {
	gorm.Model
	Name             string `gorm:"type:varchar(20);not null" json:"name,omitempty"`
	ParentCategoryID int32  `gorm:"comment:父分类ID（逻辑外键，关联自身ID）" json:"parent_category_id,omitempty"`
	// 用 constraint 禁用物理外键约束   保留是因为 查询方便
	SubCategory []*CategoryDO `gorm:"foreignKey:ParentCategoryID;references:ID;constraint:<-:false,foreignKey:no action" json:"sub_category,omitempty"`
	Level       int32         `gorm:"type:int;not null;default:1" json:"level,omitempty"`
	IsTab       bool          `gorm:"default:false;not null" json:"is_tab,omitempty"`
}

func (CategoryDO) TableName() string {
	return "category_models"
}

type CategoryDOList struct {
	TotalCount int64 `json:"totalCount,omitempty"`
	JsonData   string
	Items      []*CategoryDO `json:"items"`
}
