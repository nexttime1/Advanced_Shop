package mock

import (
	"context"

	dv1 "mxshop/app/user/srv/data/v1"
	metav1 "mxshop/pkg/common/meta/v1"
)

type users struct {
	users []*dv1.UserDO
}

func NewUsers() *users {
	return &users{}
}

func (u *users) List(ctx context.Context, opts metav1.ListMeta) (*dv1.UserDOList, error) {
	users := []*dv1.UserDO{}
	return &dv1.UserDOList{
		TotalCount: 1,
		Items:      users,
	}, nil
}
