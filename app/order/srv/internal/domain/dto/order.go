package dto

import "Advanced_Shop/app/order/srv/internal/domain/do"

type OrderDTO struct {
	do.OrderInfoDO
}

type OrderDTOList struct {
	TotalCount int64       `json:"totalCount,omitempty"`
	Items      []*OrderDTO `json:"data"`
}

type OrderInfoResponse struct {
	do.OrderInfoDO
	GoodIds    []int32
	OrderGoods []*do.OrderGoodsModel
}

type OrderDetailRequest struct {
	UserID  int32
	OrderID int32
}
