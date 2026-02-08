package v1

import (
	"context"
	metav1 "xshop/pkg/common/meta/v1"

	"testing"
	"xshop/app/user/srv/data/v1/mock"
)

func TestUserList(t *testing.T) {
	userSrv := NewUserService(mock.NewUsers())
	userSrv.List(context.Background(), metav1.ListMeta{})
}
