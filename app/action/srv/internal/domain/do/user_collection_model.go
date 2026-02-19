package do

import "Advanced_Shop/app/pkg/gorm"

type UserCollectionDO struct {
	gorm.Model
	UserId int32 `gorm:"type:int;index:idx_user_goods,unique"`
	GoodId int32 `gorm:"type:int;index:idx_user_goods,unique"`
}

func (UserCollectionDO) TableName() string {
	return "user_collection_models"
}
