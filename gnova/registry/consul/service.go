package consul

import (
	"sync"
	"sync/atomic"

	"Advanced_Shop/gnova/registry"
)

type serviceSet struct {
	serviceName string
	watcher     map[*watcher]struct{}
	services    *atomic.Value
	lock        sync.RWMutex
}

func (s *serviceSet) broadcast(ss []*registry.ServiceInstance) {
	//原子操作， 保证线程安全， 把ss（最新的服务实例列表）原子性地存到s.services里
	s.services.Store(ss)
	s.lock.RLock()
	defer s.lock.RUnlock()
	// 遍历所有监听者，逐个发通知  实现 非阻塞发送 → 能发就发，发不了就跳过  发不了说明本来里面有 那没必要发 上次通知都没收
	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default: // 由于chan 容量1 所以 满了就啥都不干
		}
	}
}
