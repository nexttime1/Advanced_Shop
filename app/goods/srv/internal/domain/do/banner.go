package do

import (
	bgorm "Advanced_Shop/app/pkg/gorm"
)

type BannerDO struct {
	bgorm.Model `structs:"-"`
	Image       string `gorm:"type:varchar(200);not null;comment: 图片的url" structs:"image"`
	Url         string `gorm:"type:varchar(200);not null;comment:跳转的详情"  structs:"url"`
	Index       int32  `gorm:"type:int;default:1;not null"  structs:"index"`
}

func (BannerDO) TableName() string {
	return "banner_models"
}

type BannerList struct {
	TotalCount int64       `json:"totalCount,omitempty"`
	Items      []*BannerDO `json:"items"`
}
