package v1

import (
	"Advanced_Shop/app/order/srv/internal/domain/do"
	"Advanced_Shop/app/order/srv/internal/domain/dto"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"context"
	"gorm.io/gorm"
)

type OrderStore interface {
	Get(ctx context.Context, detail dto.OrderDetailRequest) (*dto.OrderInfoResponse, error)

	List(ctx context.Context, userID uint64, meta metav1.ListMeta, orderby []string) (*do.OrderInfoDOList, error)

	Create(ctx context.Context, txn *gorm.DB, order *dto.OrderInfoResponse) error

	UpdateStatus(ctx context.Context, orderSn string, status string) (int64, error)
}
