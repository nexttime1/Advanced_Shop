package v1

import (
	pb "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/xshop/api/internal/data"
	"context"
)

type MessageSrv interface {
	MessageList(context.Context, *pb.MessageRequest) (*pb.MessageListResponse, error)
	CreateMessage(context.Context, *pb.MessageRequest) (*pb.MessageResponse, error)
}

type messageService struct {
	data data.DataFactory
}

func NewMessageService(data data.DataFactory) MessageSrv {
	return &messageService{
		data: data,
	}
}

func (ms *messageService) MessageList(ctx context.Context, request *pb.MessageRequest) (*pb.MessageListResponse, error) {
	return ms.data.Message().MessageList(ctx, request)
}

func (ms *messageService) CreateMessage(ctx context.Context, request *pb.MessageRequest) (*pb.MessageResponse, error) {
	return ms.data.Message().CreateMessage(ctx, request)
}

var _ MessageSrv = (*messageService)(nil)
