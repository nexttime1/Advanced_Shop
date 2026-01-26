package rpcserver

import (
	"Advanced_Shop/gnova/server/rpcserver/clientinterceptors"
	"Advanced_Shop/gnova/server/rpcserver/resolver/discovery"
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/grpclog"

	grpcinsecure "google.golang.org/grpc/credentials/insecure"

	"Advanced_Shop/gnova/registry"
	"Advanced_Shop/pkg/log"
	"google.golang.org/grpc"
	"time"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint string
	timeout  time.Duration
	// discovery接口
	discovery    registry.Discovery
	unaryInts    []grpc.UnaryClientInterceptor
	streamInts   []grpc.StreamClientInterceptor
	rpcOpts      []grpc.DialOption
	balancerName string
	/* 如果使用 withLog去传的话  这个要写成接口 不能是具体的log 因为 log里面有锁 复制锁有问题因为接口底层是
	interface {
	    tab  *itab   // 类型信息
	    data unsafe.Pointer // 指向真实对象
	}
	所以是 复制指针  这样也是可以都用一把锁的   如果选择非接口 具体的结构体 那就这里传指针
	*/
	log           log.LogHelper
	enableTracing bool
	enableMetrics bool
}

func WithEnableTracing(enable bool) ClientOption {
	return func(o *clientOptions) {
		o.enableTracing = enable
	}
}

// WithEndpoint 设置地址
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithClientTimeout 设置超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithDiscovery 设置服务发现
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// WithClientUnaryInterceptor 设置拦截器
func WithClientUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryInts = in
	}
}

// WithClientStreamInterceptor 设置stream拦截器
func WithClientStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamInts = in
	}
}

// WithClientOptions 设置grpc的dial选项
func WithClientOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.rpcOpts = opts
	}
}

// WithBalancerName 设置负载均衡器
func WithBalancerName(name string) ClientOption {
	return func(o *clientOptions) {
		o.balancerName = name
	}
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout:       2000 * time.Millisecond,
		balancerName:  "round_robin",
		enableTracing: true,
	}

	for _, o := range opts {
		o(&options)
	}

	//TODO 客户端默认拦截器
	ints := []grpc.UnaryClientInterceptor{
		clientinterceptors.TimeoutInterceptor(options.timeout),
	}
	if options.enableTracing {
		ints = append(ints, otelgrpc.UnaryClientInterceptor())
	}

	if options.enableMetrics {
		ints = append(ints, clientinterceptors.PrometheusInterceptor())
	}

	streamInts := []grpc.StreamClientInterceptor{}

	if len(options.unaryInts) > 0 {
		ints = append(ints, options.unaryInts...)
	}
	if len(options.streamInts) > 0 {
		streamInts = append(streamInts, options.streamInts...)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "` + options.balancerName + `"}`),
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(streamInts...),
	}

	//  服务发现的选项
	if options.discovery != nil {
		// 有大用
		grpcOpts = append(grpcOpts, grpc.WithResolvers( // 填一个 实现过 Build  和 Scheme 两个方法的 Builder
			// 自己实现接口的实例
			discovery.NewBuilder(
				// 这里要传一个 没有初始化 也不能初始化写死  正好 调用 dial 会传一个
				options.discovery,
				discovery.WithInsecure(insecure), // 安全吗

			),
		))
	}

	if insecure {
		// 是否 安全拨号
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}

	if len(options.rpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.rpcOpts...)

	}

	//client, err := grpc.NewClient(options.endpoint, grpcOpts...)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// 第二步：调用Connect（有返回值）建立连接
	//conn, err := client.Connect(ctx)
	//if err != nil {
	//	return nil, err
	//}
	// DialOption 也就是  我们传入的grpcOpts  本质上是 grpc  NewClient 的参数“（options）
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
