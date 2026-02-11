package gorm

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        int32          `gorm:"primarykey" structs:"-"`
	CreatedAt time.Time      `gorm:"column:add_time" structs:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" structs:"-"`
	DeletedAt gorm.DeletedAt `structs:"-"`
	idDeleted bool           `structs:"-"`
}
