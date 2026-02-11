package do

import (
	bgorm "Advanced_Shop/app/pkg/gorm"
)

type BrandsDO struct {
	bgorm.Model `structs:"-"`
	Name        string `gorm:"type:varchar(20);not null" structs:"name"`
	Logo        string `gorm:"type:varchar(200);default:'';not null"  structs:"logo"`
}

func (BrandsDO) TableName() string {
	return "brands"
}

type BrandsDOList struct {
	TotalCount int64       `json:"totalCount,omitempty"`
	Items      []*BrandsDO `json:"items"`
}
