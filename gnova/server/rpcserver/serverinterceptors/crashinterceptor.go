package serverinterceptors

import (
	"context"
	"google.golang.org/grpc"
	"runtime/debug"

	"Advanced_Shop/pkg/log"
)

func StreamCrashInterceptor(svr interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v\n \n %s", r, debug.Stack())
	})

	return handler(svr, stream)
}

func UnaryCrashInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v\n \n %s", r, debug.Stack())
	})

	return handler(ctx, req)
}

func handleCrash(solver func(interface{})) {
	// recover() 只有在 defer 修饰的函数内部调用，才能生效，作用是捕获当前 goroutine 的 panic，返回 panic 的内容
	if r := recover(); r != nil {
		solver(r)
	}
}

/*
第1步：程序执行到 defer handleCrash(...) 时，先解析传入的参数 —— 这个参数是一个【匿名函数】
       func(r interface{}) { log.Errorf("%+v\n \n %s", r, debug.Stack()) }
       这个匿名函数此时只是「定义好了、传进去了」，并没有执行！

第2步：defer 把「handleCrash函数 + 传入的匿名函数参数」一起「注册」起来，等待执行（函数退出/panic时触发）。

第3步：程序执行业务代码 handler(...)，触发 panic，此时程序立刻终止业务逻辑，跳转到执行 defer 注册的逻辑。

第4步：执行 handleCrash 函数，进入函数体：
       if r := recover(); r != nil { ... }
       这行代码执行：recover() 捕获到 panic的内容（比如"空指针异常"），赋值给r，r≠nil，条件成立。

第5步：执行 solver(r) —— 也就是执行第1步传进去的那个匿名函数，把捕获到的r传进去，最终执行 log.Errorf 打印日志+堆栈。

*/
