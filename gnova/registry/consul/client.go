package consul

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"Advanced_Shop/gnova/registry"
	"Advanced_Shop/pkg/log"

	"github.com/hashicorp/consul/api"
)

// Client is consul client config
type Client struct {
	cli    *api.Client //  Consul 官方 API 客户端实例
	ctx    context.Context
	cancel context.CancelFunc // 上下文取消函数

	// resolve service entry endpoints
	resolver ServiceResolver // 	服务解析器：将 Consul 原生ServiceEntry转换为项目自定义ServiceInstance
	// healthcheck time interval in seconds
	healthcheckInterval int // 	健康检查间隔（秒）：TCP 健康检查 / TTL 心跳的时间间隔，默认 10 秒
	// heartbeat enable heartbeat
	heartbeat bool // 是否启用 TTL 心跳保活：启用后定时更新 TTL 状态，防止服务被 Consul 标记为异常
	// deregisterCriticalServiceAfter time interval in seconds
	deregisterCriticalServiceAfter int // 服务异常后自动注销时间（秒）：默认 600 秒，Consul 核心配置
	// serviceChecks  user custom checks
	serviceChecks api.AgentServiceChecks // 用户自定义健康检查配置：会追加到默认检查（TCP/TTL）之后
}

// NewClient creates consul client
// 接收外部传入的 Consul 官方客户端cl   这样变得就更自由了 可以自己指定参数
func NewClient(cli *api.Client) *Client {
	c := &Client{
		cli:                            cli,
		resolver:                       defaultResolver, // 默认解析器
		healthcheckInterval:            10,              // 默认健康检查间隔10秒
		heartbeat:                      true,            // 默认启用心跳
		deregisterCriticalServiceAfter: 600,             // 默认异常后10分钟注销
	}
	c.ctx, c.cancel = context.WithCancel(context.Background()) // 创建可取消的上下文
	return c
}

func defaultResolver(_ context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance {
	services := make([]*registry.ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		// 1. 解析版本号：从Tags中提取"version=xxx"格式的标签
		var version string
		for _, tag := range entry.Service.Tags {
			ss := strings.SplitN(tag, "=", 2)
			if len(ss) == 2 && ss[0] == "version" {
				version = ss[1]
			}
		}
		// 2. 解析服务端点（Endpoints）
		endpoints := make([]string, 0)
		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			endpoints = append(endpoints, addr.Address)
		}
		if len(endpoints) == 0 && entry.Service.Address != "" && entry.Service.Port != 0 {
			endpoints = append(endpoints, fmt.Sprintf("http://%s:%d", entry.Service.Address, entry.Service.Port))
		}
		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: endpoints,
		})
	}

	return services
}

// ServiceResolver is used to resolve service endpoints
type ServiceResolver func(ctx context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance

// Service get services from consul
func (c *Client) Service(ctx context.Context, service string, index uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: index,
		WaitTime:  time.Second * 55,
	}
	opts = opts.WithContext(ctx)
	entries, meta, err := c.cli.Health().Service(service, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}
	return c.resolver(ctx, entries), meta.LastIndex, nil
}

// Register register service instance to consul
func (c *Client) Register(_ context.Context, svc *registry.ServiceInstance, enableHealthCheck bool) error {
	addresses := make(map[string]api.ServiceAddress, len(svc.Endpoints))
	checkAddresses := make([]string, 0, len(svc.Endpoints))
	for _, endpoint := range svc.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr := raw.Hostname()
		port, _ := strconv.ParseUint(raw.Port(), 10, 16)

		checkAddresses = append(checkAddresses, net.JoinHostPort(addr, strconv.FormatUint(port, 10)))
		addresses[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}
	}
	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		Meta:            svc.Metadata,
		Tags:            []string{fmt.Sprintf("version=%s", svc.Version)},
		TaggedAddresses: addresses,
	}
	if len(checkAddresses) > 0 {
		host, portRaw, _ := net.SplitHostPort(checkAddresses[0])
		port, _ := strconv.ParseInt(portRaw, 10, 32)
		asr.Address = host
		asr.Port = int(port)
	}
	if enableHealthCheck {
		for _, address := range checkAddresses {
			asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
				TCP:                            address,
				Interval:                       fmt.Sprintf("%ds", c.healthcheckInterval),
				DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.deregisterCriticalServiceAfter),
				Timeout:                        "5s",
			})
		}
	}
	if c.heartbeat {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			CheckID:                        "service:" + svc.ID,
			TTL:                            fmt.Sprintf("%ds", c.healthcheckInterval*2),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.deregisterCriticalServiceAfter),
		})
	}

	// custom checks
	asr.Checks = append(asr.Checks, c.serviceChecks...)

	err := c.cli.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	if c.heartbeat {
		go func() {
			time.Sleep(time.Second)
			err = c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
			if err != nil {
				log.Errorf("[Consul]update ttl heartbeat to consul failed!err:=%v", err)
			}
			ticker := time.NewTicker(time.Second * time.Duration(c.healthcheckInterval))
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err = c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
					if err != nil {
						log.Errorf("[Consul]update ttl heartbeat to consul failed!err:=%v", err)
					}
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
	return nil
}

// Deregister deregister service by service ID
func (c *Client) Deregister(_ context.Context, serviceID string) error {
	c.cancel()
	return c.cli.Agent().ServiceDeregister(serviceID)
}
