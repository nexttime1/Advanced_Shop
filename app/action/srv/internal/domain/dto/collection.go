package dto

import "Advanced_Shop/app/action/srv/internal/domain/do"

type CollectionDTO struct {
	do.UserCollectionDO
}

type CollectionDTOList struct {
	TotalCount int              `json:"total_count,omitempty"`
	Items      []*CollectionDTO `json:"data"`
}
