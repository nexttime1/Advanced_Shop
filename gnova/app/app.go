package app

import (
	"Advanced_Shop/gnova/registry"
	gs "Advanced_Shop/gnova/server"
	"Advanced_Shop/pkg/log"
	"context"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	opts options

	lk       sync.Mutex
	instance *registry.ServiceInstance

	cancel func()
}

func New(opts ...Option) *App {
	o := options{
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
		stopTimeout:      10 * time.Second,
	}

	if id, err := uuid.NewUUID(); err == nil {
		// 服务注册的 id
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &App{
		opts: o,
	}
}

// Run 启动整个服务
func (a *App) Run() error {
	//注册的信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}

	//这个变量可能被其他的goroutine访问到
	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()

	//现在启动了两个server，一个是 restserver 一个是 rpcserver
	/*
		这两个 server 要必须同时启动成功
		如果有一个启动失败，那么我们就要停止另外一个 server
		如果启动了多个， 如果其中一个启动失败，其他的应该被取消
			如果剩余的server的状态：
				1. 还没有开始调用start
					stop
				2. start进行中
					调用进行中的cancel
				3. start已经完成
					调用stop
		如果我们的服务启动了然后这个时候用户立马进行了访问
	*/

	var servers []gs.Server
	if a.opts.restServer != nil {
		servers = append(servers, a.opts.restServer)
	}
	if a.opts.rpcServer != nil {
		servers = append(servers, a.opts.rpcServer)
	}
	// 主 context
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range servers {
		//启动server
		//在启动一个goroutine 去监听是否有err产生
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() //wait for stop signal
			//不可能无休止的等待 stop
			sctx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(sctx)
		})

		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			log.Info("start server")
			return srv.Start(ctx)
		})
	}

	wg.Wait()

	//注册服务
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		err := a.opts.registrar.Register(rctx, instance)
		if err != nil {
			log.Errorf("register service error: %s", err)
			return err
		}
	}

	//监听退出信息
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return a.Stop()
		}
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

/*
http basic 认证
cache： 1. redis 2. memcache 3. local cache
jwt
*/

// Stop 停止服务
func (a *App) Stop() error {
	a.lk.Lock()
	instance := a.instance
	a.lk.Unlock()

	log.Info("start deregister service")
	if a.opts.registrar != nil && instance != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer rcancel()
		if err := a.opts.registrar.Deregister(rctx, instance); err != nil {
			log.Errorf("deregister service error: %s", err)
			return err
		}
	}

	if a.cancel != nil {
		a.cancel()
	}

	return nil
}

// 创建服务注册结构体
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	//从rpcserver， restserver去主动获取这些信息
	if a.opts.rpcServer != nil {
		if a.opts.rpcServer.Endpoint() != nil {
			endpoints = append(endpoints, a.opts.rpcServer.Endpoint().String())
		} else {
			u := &url.URL{
				Scheme: "grpc",
				Host:   a.opts.rpcServer.Address(),
			}
			endpoints = append(endpoints, u.String())
		}
	}

	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
	}, nil
}
