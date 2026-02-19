package dto

import "Advanced_Shop/app/action/srv/internal/domain/do"

type AddressDTO struct {
	do.AddressDO
}

type AddressDTOList struct {
	TotalCount int           `json:"total_count,omitempty"`
	Items      []*AddressDTO `json:"data"`
}
