package consul

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"Advanced_Shop/gnova/registry"

	"github.com/hashicorp/consul/api"
)

var (
	_ registry.Registrar = &Registry{}
	_ registry.Discovery = &Registry{}
)

// Option is consul registry option.
type Option func(*Registry)

// WithHealthCheck with registry health check option.
func WithHealthCheck(enable bool) Option {
	return func(o *Registry) {
		o.enableHealthCheck = enable
	}
}

// WithHeartbeat enable or disable heartbeat
func WithHeartbeat(enable bool) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.heartbeat = enable
		}
	}
}

// WithServiceResolver with endpoint function option.
func WithServiceResolver(fn ServiceResolver) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.resolver = fn
		}
	}
}

// WithHealthCheckInterval with healthcheck interval in seconds.
func WithHealthCheckInterval(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.healthcheckInterval = interval
		}
	}
}

// WithDeregisterCriticalServiceAfter with deregister-critical-service-after in seconds.
func WithDeregisterCriticalServiceAfter(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.deregisterCriticalServiceAfter = interval
		}
	}
}

// WithServiceCheck with service checks
func WithServiceCheck(checks ...*api.AgentServiceCheck) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.serviceChecks = checks
		}
	}
}

// Config is consul registry config
type Config struct {
	*api.Config
}

// Registry is consul registry
type Registry struct { // 实例化 实现接口
	cli               *Client
	enableHealthCheck bool
	registry          map[string]*serviceSet
	lock              sync.RWMutex
}

// New creates consul registry
func New(apiClient *api.Client, opts ...Option) *Registry {
	r := &Registry{
		cli:               NewClient(apiClient),
		registry:          make(map[string]*serviceSet),
		enableHealthCheck: true,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Register register service
func (r *Registry) Register(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Register(ctx, svc, r.enableHealthCheck)
}

// Deregister deregister service
func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

// GetService return service by name
func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	set := r.registry[name]

	getRemote := func() []*registry.ServiceInstance {
		services, _, err := r.cli.Service(ctx, name, 0, true)
		if err == nil && len(services) > 0 {
			return services
		}
		return nil
	}

	if set == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}
		return nil, fmt.Errorf("service %s not resolved in registry", name)
	}
	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if ss == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}
		return nil, fmt.Errorf("service %s not found in registry", name)
	}
	return ss, nil
}

// ListServices return service list.
func (r *Registry) ListServices() (allServices map[string][]*registry.ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	allServices = make(map[string][]*registry.ServiceInstance)
	for name, set := range r.registry {
		var services []*registry.ServiceInstance
		ss, _ := set.services.Load().([]*registry.ServiceInstance)
		if ss == nil {
			continue
		}
		services = append(services, ss...)
		allServices[name] = services
	}
	return
}

// Watch resolve service by name
func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	set, ok := r.registry[name]
	// 刚开始肯定没有 所有就创建放进去 如果服务名一样 复用 set
	if !ok {
		set = &serviceSet{
			// 管理该服务名下的所有活跃监听者，空结构体不占用内存，仅用于标记监听者的 “存在性”（后续服务实例变化时，遍历该 map 通知所有监听者）
			watcher:     make(map[*watcher]struct{}),
			services:    &atomic.Value{}, // 原子值类型 用于 读多 写少  服务实例列表
			serviceName: name,
		}
		r.registry[name] = set
	}

	// 初始化watcher
	w := &watcher{ // 代表 一个订阅者
		event: make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.set = set
	set.lock.Lock()
	// 把这一个订阅者 放进map中
	set.watcher[w] = struct{}{}
	set.lock.Unlock()
	// 原子读取存储的服务实例列表
	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if len(ss) > 0 { // 非首次监听  里面有  我直接让你去取 那边阻塞了 next方法
		w.event <- struct{}{}
	}
	// 也就是 刚初始化 还没有 serviceSet
	if !ok {
		err := r.resolve(set)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}

// 刚开始执行的函数
func (r *Registry) resolve(ss *serviceSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 首次查询 Consul 健康服务实例   idx 是 Consul 的一致性索引（Consistent Index）
	//idx 是 Consul 服务端数据版本的标识 —— 每次服务实例变化（新增 / 删除 / 健康状态变更），该索引会递增
	services, idx, err := r.cli.Service(ctx, ss.serviceName, 0, true)
	// 返回：Consul 返回的服务实例列表（包含实例地址、端口、元数据、健康状态等）
	cancel()
	if err != nil {
		return err
	} else if len(services) > 0 {
		// 广播 更新 s.services 这个  服务实例列表
		ss.broadcast(services)
	}
	go func() {
		// 创建一个定时器 1秒的
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
			tmpService, tmpIdx, err := r.cli.Service(ctx, ss.serviceName, idx, true)
			cancel()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if len(tmpService) != 0 && tmpIdx != idx {
				services = tmpService
				ss.broadcast(services)
			}
			idx = tmpIdx
		}
	}()

	return nil
}
