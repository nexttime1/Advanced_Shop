package clientinterceptors

import (
	"Advanced_Shop/pkg/errors"
	"context"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor gRPC 一元调用客户端拦截器，自动解析错误
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			// 将 gRPC 错误转换回自定义错误
			return errors.FromGrpcError(err)
		}
		return nil
	}
}

// StreamClientInterceptor gRPC 流式调用客户端拦截器，自动解析错误
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			// 将 gRPC 错误转换回自定义错误
			return cs, errors.FromGrpcError(err)
		}
		return cs, nil
	}
}
