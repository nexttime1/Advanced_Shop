package srv

import (
	upb "Advanced_Shop/api/user/v1"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/core/trace"
	"Advanced_Shop/gnova/server/rpcserver"
	"fmt"
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/alibaba/sentinel-golang/pkg/adapters/grpc"
	"github.com/alibaba/sentinel-golang/pkg/datasource/nacos"
)

func NewNacosDataSource(opts *options.NacosOptions) (*nacos.NacosDataSource, error) {
	//nacos server地址
	sc := []constant.ServerConfig{
		{
			ContextPath: "/nacos",
			Port:        opts.Port,
			IpAddr:      opts.Host,
		},
	}

	//nacos client 相关参数配置,具体配置可参考github.com/nacos-group/nacos-sdk-go
	cc := constant.ClientConfig{
		NamespaceId: opts.Namespace,
		TimeoutMs:   5000,
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return nil, err
	}

	//注册流控规则Handler
	h := datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser)
	//创建NacosDataSource数据源
	nds, err := nacos.NewNacosDataSource(client, opts.Group, opts.DataId, h)
	if err != nil {
		return nil, err
	}
	return nds, nil
}

func NewUserRPCServer(telemetry *options.TelemetryOptions, serverOpts *options.ServerOptions, userver upb.UserServer, dataNacos *nacos.NacosDataSource) (*rpcserver.Server, error) {
	//初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		telemetry.Name,
		telemetry.Endpoint,
		telemetry.Sampler,
		telemetry.Batcher,
	})

	rpcAddr := fmt.Sprintf("%s:%d", serverOpts.Host, serverOpts.Port)

	var opts []rpcserver.ServerOption
	opts = append(opts, rpcserver.WithAddress(rpcAddr))
	if serverOpts.EnableLimit {
		opts = append(opts, rpcserver.WithUnaryInterceptor(grpc.NewUnaryServerInterceptor()))
		//我去初始化nacos
		err := dataNacos.Initialize()
		if err != nil {
			return nil, err
		}
	}
	urpcServer := rpcserver.NewServer(opts...)

	upb.RegisterUserServer(urpcServer.Server, userver)

	//r := gin.Default()
	//upb.RegisterUserServerHTTPServer(userver, r)
	//r.Run(":8075")
	return urpcServer, nil
}
