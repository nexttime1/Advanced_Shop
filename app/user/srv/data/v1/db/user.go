package db

import (
	"Advanced_Shop/app/pkg/code"
	dv1 "Advanced_Shop/app/user/srv/data/v1"
	code2 "Advanced_Shop/gnova/code"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) dv1.UserStore {
	return &users{db: db}
}

func (u *users) GetByMobile(ctx context.Context, mobile string) (*dv1.UserDO, error) {
	user := dv1.UserDO{}
	err := u.db.Where("mobile = ?", mobile).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

func (u *users) GetByID(ctx context.Context, id uint64) (*dv1.UserDO, error) {
	user := dv1.UserDO{}
	err := u.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

func (u *users) Create(ctx context.Context, user *dv1.UserDO) error {
	tx := u.db.Create(user)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (u *users) Update(ctx context.Context, user *dv1.UserDO) error {
	tx := u.db.Save(user)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func newUsers(db *gorm.DB) *users {
	return &users{db: db}
}

var _ dv1.UserStore = &users{}

func (u *users) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*dv1.UserDOList, error) {
	ret := &dv1.UserDOList{}

	var limit, offset int
	if opts.PageSize == 0 {
		limit = 10
	} else {
		limit = opts.PageSize
	}

	if opts.Page > 0 {
		offset = (opts.Page - 1) * limit
	}

	query := u.db
	for _, value := range orderby {
		query = query.Order(value)
	}

	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}
