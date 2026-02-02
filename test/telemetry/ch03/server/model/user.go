package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        int32          `gorm:"primarykey" structs:"-"`
	CreatedAt time.Time      `gorm:"column:add_time" structs:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" structs:"-"`
	DeletedAt gorm.DeletedAt `structs:"-"`
	idDeleted bool           `structs:"-"`
}

type UserModels struct {
	Model    `structs:"-"`
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null" structs:"-"`
	Password string     `gorm:"type:varchar(100);not null" structs:"password"`
	NickName string     `gorm:"type:varchar(100);"  structs:"nick_name"`
	Birthday *time.Time `gorm:"type:datetime" structs:"birthday"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6)"  structs:"gender"`
	Role     int        `gorm:"column: role;default 2"  structs:"role"` // 1管理员  2 普通用户
}
