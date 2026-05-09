package serverinterceptors

import (
	"Advanced_Shop/pkg/errors"
	"context"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor gRPC 一元调用服务端拦截器，自动转换错误
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			// 将自定义错误转换为 gRPC 错误
			return resp, errors.ToGrpcError(err)
		}
		return resp, nil
	}
}

// StreamServerInterceptor gRPC 流式调用服务端拦截器，自动转换错误
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err != nil {
			// 将自定义错误转换为 gRPC 错误
			return errors.ToGrpcError(err)
		}
		return nil
	}
}
