package v1

import (
	v1 "Advanced_Shop/app/action/srv/internal/data/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"Advanced_Shop/app/action/srv/internal/domain/dto"
	"Advanced_Shop/pkg/log"
	"context"
)

// MessageSrv 留言业务逻辑层接口
type MessageSrv interface {
	// MessageList 根据用户ID获取留言列表
	MessageList(ctx context.Context, userID int32) (*dto.LeavingMessageDTOList, error)

	// CreateMessage 创建留言
	CreateMessage(ctx context.Context, messageDTO *dto.LeavingMessageDTO) (*dto.LeavingMessageDTO, error)
}

type messageService struct {
	data v1.DataFactory
}

func newMessage(srv *serviceFactory) MessageSrv {
	return &messageService{
		data: srv.data,
	}
}

// MessageList 根据用户ID获取留言列表
func (s *messageService) MessageList(ctx context.Context, userID int32) (*dto.LeavingMessageDTOList, error) {
	// 调用数据层获取DO列表和总数
	messageDOs, count, err := s.data.Messages().ListByUserID(ctx, userID)
	if err != nil {
		log.Errorf("MessageList failed: %v", err)
		return nil, err
	}

	// DO转换为DTO
	dtoList := &dto.LeavingMessageDTOList{
		TotalCount: int(count),
		Items:      make([]*dto.LeavingMessageDTO, 0, len(messageDOs)),
	}

	for _, doItem := range messageDOs {
		dtoItem := &dto.LeavingMessageDTO{
			LeavingMessageDO: *doItem,
		}
		dtoList.Items = append(dtoList.Items, dtoItem)
	}

	return dtoList, nil
}

// CreateMessage 创建留言
func (s *messageService) CreateMessage(ctx context.Context, messageDTO *dto.LeavingMessageDTO) (*dto.LeavingMessageDTO, error) {
	// DTO转换为DO
	messageDO := &do.LeavingMessageDO{
		UserId:      messageDTO.UserId,
		MessageType: messageDTO.MessageType,
		Subject:     messageDTO.Subject,
		Message:     messageDTO.Message,
		File:        messageDTO.File,
	}

	// 调用数据层创建留言
	err := s.data.Messages().Create(ctx, messageDO)
	if err != nil {
		log.Errorf("CreateMessage failed: %v", err)
		return nil, err
	}

	// 设置创建后的ID并返回
	messageDTO.ID = messageDO.ID
	return messageDTO, nil
}

// 确保实现了接口
var _ MessageSrv = &messageService{}
