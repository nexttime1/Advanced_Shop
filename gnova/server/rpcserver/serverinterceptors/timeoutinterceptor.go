package serverinterceptors

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryTimeoutInterceptor returns a func that sets timeout to incoming unary requests.  闭包
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var resp interface{}
		var err error
		//强制编译器和CPU不重排指令
		//确保之前的写操作在解锁前完成并对其他核可见
		//确保锁定方看到最新的数据
		var lock sync.Mutex

		done := make(chan struct{}) //只用来“发信号”，不传数据
		// create channel with buffer size 1 to avoid goroutine leak 防止泄露
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// attach call stack to avoid missing in different goroutine
					//Sprintf 例子
					/*
						runtime error: index out of range

						goroutine 10 [running]:
						main.foo()
						    /app/main.go:12
						...
					*/
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			lock.Lock()
			defer lock.Unlock()
			resp, err = handler(ctx, req)
			close(done)
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-done:
			lock.Lock()
			defer lock.Unlock()
			return resp, err
		case <-ctx.Done():
			err := ctx.Err()
			if err == context.Canceled {
				//我们之前说过我们把error统一了， grpc的error我们也可以统一, 自己完成 TODO
				err = status.Error(codes.Canceled, err.Error())
			} else if err == context.DeadlineExceeded {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, err
		}
	}
}
