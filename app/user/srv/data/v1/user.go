package v1

import (
	"context"
	"gorm.io/gorm"
	"time"

	metav1 "Advanced_Shop/pkg/common/meta/v1"
)

type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

/*
1. 密文 2. 密文不可反解
 1. 对称加密
 2. 非对称加密
 3. md5 信息摘要算法
    密码如果不可以反解，用户找回密码
*/
type UserDO struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
}

func (u *UserDO) TableName() string {
	return "user"
}

type UserDOList struct {
	TotalCount int64     `json:"totalCount,omitempty"` //总数
	Items      []*UserDO `json:"data"`                 //数据
}

type UserStore interface {
	/*
		有数据访问的方法，一定要有error
		参数中最好有ctx
	*/
	//用户列表
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDOList, error)

	//通过手机号码查询用户
	GetByMobile(ctx context.Context, mobile string) (*UserDO, error)

	//通过用户ID查询用户
	GetByID(ctx context.Context, id uint64) (*UserDO, error)

	//创建用户
	Create(ctx context.Context, user *UserDO) error

	//更新用户
	Update(ctx context.Context, user *UserDO) error
}
