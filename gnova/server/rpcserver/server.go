package rpcserver

import (
	"context"

	srvintc "Advanced_Shop/gnova/server/rpcserver/serverinterceptors"
	"Advanced_Shop/pkg/host"
	"Advanced_Shop/pkg/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/url"
	"time"

	apimd "Advanced_Shop/api/metadata"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server

	address    string
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	grpcOpts   []grpc.ServerOption
	lis        net.Listener
	timeout    time.Duration

	health   *health.Server
	metadata *apimd.Server
	endpoint *url.URL

	enableMetrics bool
}

func (s *Server) Endpoint() *url.URL {
	return s.endpoint
}

func (s *Server) Address() string {
	return s.address
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address: ":0",
		health:  health.NewServer(),
		//timeout: 1 * time.Second,
	}

	for _, o := range opts {
		o(srv)
	}

	//TODO 我们现在希望用户不设置拦截器的情况下，我们会自动默认加上一些必须的拦截器， crash，tracing
	unaryInts := []grpc.UnaryServerInterceptor{
		srvintc.UnaryCrashInterceptor,
		otelgrpc.UnaryServerInterceptor(),
	}

	if srv.enableMetrics {
		unaryInts = append(unaryInts, srvintc.UnaryPrometheusInterceptor)
	}

	if srv.timeout > 0 {
		unaryInts = append(unaryInts, srvintc.UnaryTimeoutInterceptor(srv.timeout))
	}

	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}

	//把我们传入的拦截器转换成grpc的ServerOption  在它内部也是我们一样的方法  type ServerOption func(o *Server) 放到它那边
	grpcOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(unaryInts...)}

	//把用户自己传入的grpc.ServerOption放在一起
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)

	//注册metadata的Server  创建一个元数据（Metadata）服务实例 这个实例会自动调用 srv.Server.GetServiceInfo() 方法
	srv.metadata = apimd.NewServer(srv.Server) //load() 函数会负责把服务名、方法名、文件描述符集合（FileDescriptorSet） 这些数据缓存起来，供后续查询使用

	//解析address
	err := srv.listenAndEndpoint()
	if err != nil {
		panic(err)
	}
	// 细节 去看 grpc注册.md
	// 健康检查注册 gRPC 服务端 内部注册一个服务   Consul Client 来调用它
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	// 所有想对外提供的 RPC 能力（健康检查、元数据查询、写的业务接口、支付 / 用户 / 订单服务），
	// 都必须通过 RegisterXXXServer() 这个动作，把「服务的定义 + 服务的实现」挂载到 gRPC 服务端。
	// 不管是健康检查、Metadata、还是你自己写的业务服务，底层逻辑完全一模一样 因为都是通过protoc + protoc-gen-go-grpc  生成
	apimd.RegisterMetadataServer(srv.Server, srv.metadata)

	reflection.Register(srv.Server)
	//可以支持用户直接通过grpc的一个接口查看当前支持的所有的rpc服务

	return srv
}

func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

func WithMetrics(metric bool) ServerOption {
	return func(s *Server) {
		s.enableMetrics = metric
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithLis(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = in
	}
}

func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = in
	}
}

func WithOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

// 完成ip和端口的提取
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	return nil
}

// Start 启动grpc的服务
func (s *Server) Start(ctx context.Context) error {
	log.Infof("[grpc] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	//设置服务的状态为not_serving，防止接收新的请求过来
	s.health.Shutdown()
	s.GracefulStop()
	log.Infof("[grpc] server stopped")
	return nil
}
