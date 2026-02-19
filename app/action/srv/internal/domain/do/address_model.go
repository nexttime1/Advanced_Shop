package do

import "Advanced_Shop/app/pkg/gorm"

type AddressDO struct {
	gorm.Model
	UserId       int32  `gorm:"type:int;index"`
	Province     string `gorm:"type:varchar(10)"`
	City         string `gorm:"type:varchar(10)"`
	District     string `gorm:"type:varchar(20)"`
	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SignerMobile string `gorm:"type:varchar(11)"`
}

func (AddressDO) TableName() string {
	return "addresses"
}
