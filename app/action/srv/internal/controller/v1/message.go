package v1

import (
	pb "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"Advanced_Shop/app/action/srv/internal/domain/dto"
	"context"
)

// MessageList 获取留言列表
func (o *actionServer) MessageList(ctx context.Context, request *pb.MessageRequest) (*pb.MessageListResponse, error) {
	// 调用业务层
	dtoList, err := o.srv.Message().MessageList(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	// DTO转换为Proto响应
	response := &pb.MessageListResponse{
		Total: int32(dtoList.TotalCount),
		Data:  make([]*pb.MessageResponse, 0, len(dtoList.Items)),
	}

	for _, dtoItem := range dtoList.Items {
		response.Data = append(response.Data, &pb.MessageResponse{
			Id:          dtoItem.ID,
			UserId:      dtoItem.UserId,
			MessageType: int32(dtoItem.MessageType),
			Subject:     dtoItem.Subject,
			Message:     dtoItem.Message,
			File:        dtoItem.File,
		})
	}

	return response, nil
}

// CreateMessage 创建留言
func (o *actionServer) CreateMessage(ctx context.Context, request *pb.MessageRequest) (*pb.MessageResponse, error) {
	// Proto转换为DTO
	messageDTO := &dto.LeavingMessageDTO{
		LeavingMessageDO: do.LeavingMessageDO{
			UserId:      request.UserId,
			MessageType: do.MessageType(request.MessageType),
			Subject:     request.Subject,
			Message:     request.Message,
			File:        request.File,
		},
	}

	// 调用业务层
	createdDTO, err := o.srv.Message().CreateMessage(ctx, messageDTO)
	if err != nil {
		return nil, err
	}

	// DTO转换为Proto响应
	return &pb.MessageResponse{
		Id: createdDTO.ID,
	}, nil
}
