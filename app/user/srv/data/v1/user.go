package v1

import (
	bgorm "Advanced_Shop/app/pkg/grom"
	"context"
	"time"

	metav1 "Advanced_Shop/pkg/common/meta/v1"
)

type UserDO struct {
	bgorm.Model `structs:"-"`
	Mobile      string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null" structs:"-"`
	Password    string     `gorm:"type:varchar(100);not null" structs:"password"`
	NickName    string     `gorm:"type:varchar(100);"  structs:"nick_name"`
	Birthday    *time.Time `gorm:"type:datetime" structs:"birthday"`
	Gender      string     `gorm:"column:gender;default:male;type:varchar(6)"  structs:"gender"`
	Role        int        `gorm:"column: role;default 2"  structs:"role"` // 1管理员  2 普通用户
}

type UserDOList struct {
	TotalCount int64     `json:"totalCount,omitempty"` //总数
	Items      []*UserDO `json:"data"`                 //数据
}

type UserStore interface {
	// List 用户列表
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDOList, error)

	// GetByMobile 通过手机号码查询用户
	GetByMobile(ctx context.Context, mobile string) (*UserDO, error)

	// GetByID 通过用户ID查询用户
	GetByID(ctx context.Context, id uint64) (*UserDO, error)

	// Create 创建用户
	Create(ctx context.Context, user *UserDO) error

	// Update 更新用户
	Update(ctx context.Context, user *UserDO) error
}
