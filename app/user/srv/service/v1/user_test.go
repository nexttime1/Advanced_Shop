package v1

import (
	"context"
	metav1 "mxshop/pkg/common/meta/v1"

	"mxshop/app/user/srv/data/v1/mock"
	"testing"
)

func TestUserList(t *testing.T) {
	userSrv := NewUserService(mock.NewUsers())
	userSrv.List(context.Background(), metav1.ListMeta{})
}
