package driver

import (
	"Advanced_Shop/gnova/registry/etcd"
	"context"
	"fmt"
	etcdAPI "go.etcd.io/etcd/client/v3"

	// Kratos的consul/etcd注册中心适配包
	"Advanced_Shop/gnova/registry" // Kratos注册中心抽象接口
	"Advanced_Shop/gnova/registry/consul"
	// Kratos的gRPC resolver（核心：解析discovery://协议）
	_ "Advanced_Shop/gnova/server/rpcserver/resolver/direct"
	"Advanced_Shop/gnova/server/rpcserver/resolver/discovery"
	consulAPI "github.com/hashicorp/consul/api" // Consul原生客户端
	"google.golang.org/grpc/resolver"           // gRPC地址解析器接口
	"net/url"                                   // URL解析工具
	"strings"                                   // 字符串处理
)

// 常量定义：驱动名和协议名
const (
	DriverName    = "dtm-driver-kratos" // 自定义驱动的名称（唯一标识）
	DefaultScheme = "discovery"         // Kratos默认的服务发现协议
	EtcdScheme    = "etcd"              // Etcd协议
	ConsulScheme  = "consul"            // Consul协议
)

// kratosDriver：实现DTM Driver接口的结构体
// 注：结构体本身无字段，因为核心逻辑依赖Kratos的注册中心，不需要存储状态
type xtmDriver struct {
	consulAddr   string // Consul地址（比如127.0.0.1:8500）
	consulScheme string // Consul协议（http/https）
}

func NewXtmDriver(consulAddr, consulScheme string) *xtmDriver {
	return &xtmDriver{
		consulAddr:   consulAddr,
		consulScheme: consulScheme,
	}
}

// GetName 返回驱动名称（DTM驱动接口要求）
// 作用：DTM通过这个方法识别驱动，对应dtmdriver.Register()的key
func (k *xtmDriver) GetName() string {
	return DriverName // 返回常量"dtm-driver-kratos"
}

// RegisterAddrResolver 注册地址解析器（DTM驱动接口要求）
// 作用：让DTM能识别并解析自定义协议（比如discovery://）
// 这个实现是空的，因为Kratos的resolver在RegisterService中注册了
func (k *xtmDriver) RegisterAddrResolver() {
	// 空实现原因：Kratos的resolver在RegisterService中通过resolver.Register()注册，无需在这里重复注册
}

// RegisterService 注册服务实例到注册中心（DTM驱动接口要求）
// 作用：
// 1. 解析target（比如"discovery:///xshop-inventory-srv"）；
// 2. 根据协议（etcd/consul/discovery）创建Kratos注册中心客户端；
// 3. 注册Kratos的resolver，让DTM能解析discovery://协议；
// 4. 把服务实例注册到注册中心（可选，如果你需要DTM自身注册服务）。
// 参数：
// - target：服务地址（比如"discovery:///xshop-inventory-srv"）；
// - endpoint：服务端点（比如"grpc://127.0.0.1:8081"）。
func (k *xtmDriver) RegisterService(target string, endpoint string) error {
	// 1. 如果target为空，直接返回（无需注册）
	if target == "" {
		return nil
	}

	// 2. 解析target为URL对象（比如把"discovery:///xshop-inventory-srv"解析成URL结构体）
	u, err := url.Parse(target)
	if err != nil {
		return err // 解析失败返回错误
	}

	// 3. 根据URL的Scheme（协议）分支处理
	switch u.Scheme {
	// 处理discovery://或etcd://协议
	case DefaultScheme: // discovery://
		fallthrough // 穿透到EtcdScheme逻辑（因为discovery默认基于etcd）

	// 处理consul://协议
	case ConsulScheme: // consul://
		// 1. 创建Consul客户端（和你的NewRegistrar逻辑完全一致）
		c := consulAPI.DefaultConfig()
		c.Address = k.consulAddr  // 用你的Consul地址
		c.Scheme = k.consulScheme // 用你的Consul协议
		cli, err := consulAPI.NewClient(c)
		if err != nil {
			return err
		}

		// 封装成Kratos的Consul注册器（开启健康检查，和你一致）
		consulRegistry := consul.New(cli, consul.WithHealthCheck(true))

		//注册我们的discovery resolver
		// 作用：让gRPC（DTM底层用gRPC）能识别discovery://协议
		// 这一步后，DTM就能从Consul中拉取你已注册的服务地址
		resolver.Register(discovery.NewBuilder(consulRegistry, discovery.WithInsecure(true)))

		// 服务已经注册到Consul了，不需要DTM再注册
		// 直接返回nil，只完成resolver注册
		return nil

	case EtcdScheme: // etcd://
		// 3.1 构建Kratos的服务实例对象
		registerInstance := &registry.ServiceInstance{
			Name:      strings.TrimPrefix(u.Path, "/"), // 服务名（比如xshop-inventory-srv）
			Endpoints: strings.Split(endpoint, ","),    // 服务端点（拆分多个地址）
		}
		// 3.2 创建Etcd原生客户端（地址从URL.Host获取，比如"127.0.0.1:2379"）
		client, err := etcdAPI.New(etcdAPI.Config{
			Endpoints: strings.Split(u.Host, ","), // 支持多个etcd节点（逗号分隔）
		})
		if err != nil {
			return err // Etcd客户端创建失败
		}
		// 3.3 创建Kratos的Etcd注册中心
		registry := etcd.New(client)
		// 3.4 核心：注册Kratos的discovery resolver到gRPC
		// 作用：让gRPC（DTM底层用gRPC调用）能解析discovery://协议
		resolver.Register(discovery.NewBuilder(registry, discovery.WithInsecure(true)))
		// 3.5 （可选）把服务实例注册到Etcd（如果需要DTM自身注册服务）
		return registry.Register(context.Background(), registerInstance)
	// 未知协议：返回错误
	default:
		return fmt.Errorf("unknown scheme: %s", u.Scheme)
	}
}

// ParseServerMethod 解析URI，拆分“服务地址”和“方法名”（DTM驱动接口要求）
// 作用：DTM调用服务时，需要把完整URI拆分成“服务地址”和“方法名”，比如：
// 输入："discovery:///xshop-inventory-srv/Inventory/Sell"
// 输出：server="discovery:///xshop-inventory-srv", method="/Inventory/Sell"
func (k *xtmDriver) ParseServerMethod(uri string) (server string, method string, err error) {
	// 情况1：URI不含"//"（比如"127.0.0.1:8081/Inventory/Sell"）
	if !strings.Contains(uri, "//") {
		sep := strings.IndexByte(uri, '/') // 找第一个"/"的位置
		if sep == -1 {                     // 没有"/"，格式错误
			return "", "", fmt.Errorf("bad url: '%s'. no '/' found", uri)
		}
		// 拆分：server=127.0.0.1:8081，method=/Inventory/Sell
		return uri[:sep], uri[sep:], nil
	}

	// 情况2：URI含"//"（比如"discovery:///xshop-inventory-srv/Inventory/Sell"）
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", nil // 解析失败返回空（DTM会处理错误）
	}
	// 找Path中第一个"/"的位置（比如Path="/Inventory/Sell"，index=0）
	index := strings.IndexByte(u.Path[1:], '/') + 1
	// 拆分：
	// server = scheme://host + Path[:index] → discovery:///xshop-inventory-srv
	// method = Path[index:] → /Inventory/Sell
	return u.Scheme + "://" + u.Host + u.Path[:index], u.Path[index:], nil
}
