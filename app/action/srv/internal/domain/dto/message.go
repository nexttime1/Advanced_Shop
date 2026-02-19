package dto

import "Advanced_Shop/app/action/srv/internal/domain/do"

type LeavingMessageDTO struct {
	do.LeavingMessageDO
}

type LeavingMessageDTOList struct {
	TotalCount int                  `json:"total_count,omitempty"`
	Items      []*LeavingMessageDTO `json:"data"`
}
