# errgroup 

专门用于简化一组 goroutine 的生命周期管理和错误收集，解决了普通`WaitGroup`无法便捷收集错误、无法快速取消其他 goroutine 的痛点。

## errgroup vs 原生 WaitGroup

|        特性         |       sync.WaitGroup        |             errgroup              |
| :-----------------: | :-------------------------: | :-------------------------------: |
| 等待 goroutine 结束 |     ✅ 需要手动 Add/Done     |     ✅ 自动管理，无需手动调用      |
|      收集错误       |  ❌ 需要自己用 channel 实现  |    ✅ 自动收集第一个非 nil 错误    |
| 取消其他 goroutine  | ❌ 需要自己结合 Context 实现 |     ✅ 内置 Context，自动取消      |
|     使用复杂度      |    低（但扩展功能复杂）     | 中（一站式解决多 goroutine 问题） |

## 示例

下面代码是 通过开启三个 goroutine 并发执行，能够感知错误，并且便捷收集第一个错误

`context.Cause(ctx)` 这个函数是专门去

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		fmt.Println("doing task1")
		time.Sleep(5 * time.Second)
		return errors.New("task1 error") // 业务错误：作为取消原因 这边取消 其他 goroutine  就会走 <-ctx.Done()
	})

	eg.Go(func() error {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("doing task2")
			case <-ctx.Done():
				fmt.Println("task2 canceled")
				// 可以调用 context.Cause()，获取取消的具体原因 
				cause := context.Cause(ctx)
				fmt.Printf("task2 被取消的原因：%v\n", cause)
				return ctx.Err()  //这时候 return 的错误 是第二次或者第三次了 我们只记录第一次的
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("doing task3")
			case <-ctx.Done():
				fmt.Println("task3 canceled")
				// 可以调用 context.Cause()，获取取消的具体原因
				cause := context.Cause(ctx)
				fmt.Printf("task3 被取消的原因：%v\n", cause)
				return ctx.Err()
			}
		}
	})

	err := eg.Wait()
	if err != nil {
		fmt.Println("task failed")
		// 在 Wait() 后, 返回第一次错误
		fmt.Printf("整体任务被取消的原因：%v\n", err)
	} else {
		fmt.Println("task success")
	}
}

```

## errgroup 源码解析

### 结构体

```GO
// 空结构体在 Go 中不占用字节内存，是标准用于通知的
type token struct{}

type Group struct {
	cancel func(error)  //  取消下级关联 Context 的函数

	wg sync.WaitGroup  // 管理 goroutine 的等待，替代手动 Add/Done

	sem chan token  // 带缓冲的通道 限制最大能有几个 goroutine 并行 用 SetLimit 函数去指明 不用这个函数那这个字段就不用管

	errOnce sync.Once  // 原子语句  无论多少 goroutine 写入 都只执行一次 
	err     error
}

```

### 入口函数

新版本使用了 `WithCancelCause`  换掉了传统的 `context.WithCancel` 

对于  `context.WithCancelCause` 与  `context.WithCancel` 的 源码 放在文章最后，深入理解 `context`  可阅读	

```GO
func WithContext(ctx context.Context) (*Group, context.Context) {
    //返回 子Context 和 可以取消的函数  这里覆盖了主 ctx 无所谓
	ctx, cancel := context.WithCancelCause(ctx)
 	// 返回 errgroup 内置 结构体 和 context
	return &Group{cancel: cancel}, ctx
}

```

### `(g *Group) SetLimit(n int)` 函数

用于控制同时活跃的 goroutine 数量，实现并发限制，小于0，就没有限制，也就不用 `g.sem` 字段

必须是在初始化的时候去调用 `SetLimit`  函数，必须保证当前无活跃 goroutine，否则直接 panic 为了安全

```GO
func (g *Group) SetLimit(n int) {
	if n < 0 {
		g.sem = nil
		return
	}
	if active := len(g.sem); active != 0 {
		panic(fmt.Errorf("errgroup: modify limit while %v goroutines in the group are still active", active))
	}
    // 启动 容量为 n 的 cannel 
	g.sem = make(chan token, n)
}

```

### `(g *Group) done()` 函数

本质对 `WaitGroup` 封装了一个 对 goroutine 数量限制 不设置的话 就相当于 `WaitGroup.Done()` 

```GO
func (g *Group) done() {
	if g.sem != nil {  // 如果设置了goroutine数量限制
		<-g.sem       // 释放信号量，从sem通道取出一个token，允许新的goroutine启动
	}
	g.wg.Done()       // 通知 WaitGroup：一个goroutine执行完成
}
```

### `(g *Group) Wait()` 函数

同样也是多封装了一下 调用 `cancel(g.err) `方法 去记录错误（如果有的话） 内部具体细节可以看最后源码

```GO
func (g *Group) Wait() error {
	g.wg.Wait()          // 阻塞，直到所有goroutine调用了wg.Done()
	if g.cancel != nil { // 如果绑定了Context（通过WithContext创建）
		g.cancel(g.err)  // 取消Context，传递错误作为取消原因
	}
	return g.err         // 返回第一个非nil错误（无错误则返回nil）
}
```

### 核心 `(g *Group) Go` 函数

这里并没有捕获 `panic` ，如果捕捉的话可能导致 `panic` 时机延迟，调试困难，`panic` 栈会被转为普通值，无法被监控工具捕获，而且还可能导致死锁

比如 两个` goroutine `被创建 第一个发生 `panic `第二个在阻塞等待，`Wait` 需要等待两个任务结束才返回，死锁

```GO
func (g *Group) Go(f func() error) {
	// 若设置了goroutine数量限制，先获取信号量 如果已经满了 阻塞直到有空闲token
	if g.sem != nil {
		g.sem <- token{} // 发送空token到sem通道，占用一个并发名额
	}

	// 增加WaitGroup计数
	g.wg.Add(1)
	
	// 启动新goroutine执行任务
	go func() {
		defer g.done() // 最后执行无论错误还是正常

		// 执行用户传入的函数，获取错误
		if err := f(); err != nil {
			// 仅第一次执行  保证只存第一个错误
			g.errOnce.Do(func() {
				g.err = err          // 存储第一个错误
				if g.cancel != nil { // 若绑定了Context，取消所有关联 goroutine 大部分都会有的 因为我们在入口传了
					g.cancel(g.err)
				}
			})
		}
	}()
}

// 这个函数和 Go 的唯一区别就是 它不阻塞 一旦数量满了 直接推出
func (g *Group) TryGo(f func() error) bool {
    	if g.sem != nil {
		select {
		case g.sem <- token{}: // 成功获取token
		default:               // 无空闲token  并发数达上限，直接返回false
			return false
		}
	}
    ...  //逻辑一样
}

```

## Context 部分源码

### `context.WithCancelCause` 与  `context.WithCancel` 的 源码 区别

简单理解就是在返回 `cancel`函数的时候多了一个参数而已

```GO
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := withCancel(parent)
	return c, func() { c.cancel(true, Canceled, nil) }
}

func WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc) {
	c := withCancel(parent)
	return c, func(cause error) { c.cancel(true, Canceled, cause) }
}

```

### 先介绍核心结构体 `cancelCtx`  和 `cancelCtx.cancel` 方法

`context.WithCancel` 和 `context.WithCancelCause` 以及 `context.WithTimeout` 都有它的影子

```go
// 核心结构体 cancelCtx
type cancelCtx struct {
	Context          //       嵌入父Context，继承父Context的所有方法（Done/Err/Value等）  该字段是个接口
	mu       sync.Mutex            // 保护以下字段的并发读写安全
	done     atomic.Value       // 存储chan struct{}，负责通知外部阻塞函数 查看 Done方法 和 Cancel 方法 之后就明白
	children map[canceler]struct{} // 存储当前Context的子canceler，取消时会遍历取消所有子Context
	err      atomic.Value          // 存储取消时的错误（如context.Canceled），原子操作避免锁竞争
	cause    error                 // WithCancelCause的核心 ：存储取消原因，仅在第一次取消时赋值
}

// context包内部的接口，定义了「可取消」的行为
type canceler interface {
	cancel(removeFromParent bool, err error, cause error)
	Done() <-chan struct{}
}
```

### 核心方法 `cancelCtx.cancel`

这个函数 是写方法，使用频率是极低的 可能就是手动 cancel 或者 父 context 进行 cancel，为了防止并发问题也就是同时进行 `c.err.Load()` 发现没有，然后都执行了 `c.err.Store(err)` 是有问题的 我们只记录第一次 而且后面 `close(d)` 会报错，这里必须用锁，也不会影响性能

读方法 `Err()、Done()` 方法 常见 `for + select` 可能每个 goroutine 每毫秒跑几万次，这种就不能加锁，太影响性能 这个时候就体现出来 `atomic.Value ` 这个字段的好处，防止一个 goroutine 在 `cancel() `里写 `c.err`，另一个 goroutine 同时在 Err() 里读 `c.err`

```GO
func (c *cancelCtx) cancel(removeFromParent bool, err, cause error) {
    // Context 设计原则：取消必须关联错误
    if err == nil {
        panic("context: internal error: missing cancel error")
    }
    
    // 旧版本 并没有 WithCancelCause 也就没有这个字段 兼容 旧逻辑
    if cause == nil {
        cause = err
    }

    // 上面已经解释
    c.mu.Lock()
    // 看一下 里面有没有数据
    if c.err.Load() != nil {  //原子地把值拿出来
        c.mu.Unlock()
        return  // 说明有数据 已经写过了 只要写过了 就一定是在这个函数中写的 说明运行过一次 直接return
    }

    // 到这说明第一次进 cancel 方法 写入错误下一个 goroutine 再来调用 直接 return
    c.err.Store(err)   // 原子操作：高频读取时无需加锁，性能更优
    c.cause = cause    // 保证了仅第一次赋值有效  因为第二次直接 return 

    
    // 关闭取消信号通道，用于通知外部 select { case <-ctx.Done(): } 消除阻塞
    d, _ := c.done.Load().(chan struct{})
    if d == nil {
        // 懒创建 下面会详细讲
        c.done.Store(closedchan)
    } else {
        // 若 done 通道已创建，关闭通道 消除阻塞
        close(d)
    }

    // 递归取消所有子 Context
    for child := range c.children {
        // 持有父锁时获取子锁会有嵌套锁风险，但 Context 取消是低频操作，可接受 
        child.cancel(false, err, cause)
    }
    // 清空子列表，释放内存，避免泄漏
    c.children = nil

    //解锁
    c.mu.Unlock()

    // 基本上都是 True  默认移除 避免内存泄漏
    if removeFromParent {
        // 从父 Context 移除自身
        removeChild(c.Context, c)
    }
}

```

### 懒创建

我们在写代码的时候，很少去写 `ctx.Done()`  方法，非常常见的两种，要么不用，要么只用  `cancel`

```go
// 第一种 不用
func work(ctx context.Context) {
    time.Sleep(10 * time.Second)
}

// 第二种 只 cancel
ctx, cancel := context.WithCancel(parent)
defer cancel()

doSomething(ctx) // ctx 只是被传来传去
```

正常逻辑是 建完 `ctx` 要使用 `Done()`方法

它进行了两次判断 其实可以进行一次 在刚开始进来的时候加锁就可以直接判断是否有值，可以一样的效果，但是这个函数读极多，就写一次，所以大多数情况下都直接返回，根本不需要锁，如果按刚开始就加锁，性能直线下降

```GO
func (c *cancelCtx) Done() <-chan struct{} {  //返回一个只能接收的 channel
	d := c.done.Load() //第一次读取
	if d != nil {
        //之前已经调用过 Done() 
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.done.Load() // 第二次重新读取
	if d == nil {
        // 第一次 初始化并存入
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}
```

解释  ` Done()`  是 极高频函数

go语言非常常见的 阻塞循环，每次执行 `ctx.Done()` 函数 返回的没有数据的 channel 只能走 default 逻辑，default 啥也不干 然后再回到 select 语句中 在执行 `ctx.Done()` 函数 直到 调用 `concel` 函数 ，可能每毫秒跑几万次，只第一次进行了修改，进行了两次判断，其他的都是直接返回，所以不能在刚开始进行加锁

```GO
for {
    select {
    case <-ctx.Done():
        return
    default:
    }
}
```

### 深入 `withCancel` 函数



我们都知道 `parent` 参数一般情况就是我们自己传的 `ctx` 可能是 `context.WithCancel` 或者 `context.Background`

在标准项目或者大型项目中，在主要内部结构体中会定义 `cancel` 字段用于 优雅退出

如果 `parent` 是  `context.Background`  没啥意义，

```go
func withCancel(parent Context) *cancelCtx {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	// 1. 创建空的cancelCtx实例
	c := &cancelCtx{}
	// 2. 建立当前 cancelCtx 和父 Context 的关联，作用是父取消时，子也会被取消
	c.propagateCancel(parent, c)
	return c
}
```

### `propagateCancel` 函数

```go
func (c *cancelCtx) propagateCancel(parent Context, child canceler) {
	// 挂载到字段 
	c.Context = parent

	// 查看 父Context是否可取消
	done := parent.Done()
	if done == nil {
		return 		//  如果是 context.Background 无需关联，父永远不会取消，子只能手动取消
	}

	// 检查父Context是否已经取消  非阻塞
	select {
	case <-done:
		// 父已取消：立即取消子Context，继承父的错误和原因  这里的child 也就是我们刚创建的  空的cancelCtx实例
		child.cancel(false, parent.Err(), Cause(parent))
		return
	default:  // 父没取消  正常情况，在代码健壮的情况下 不太可能刚创建就取消了
	}

	// 断言 将父Context转为 *cancelCtx  几乎所有的Context 底层都有它的影子 这是最常见的
	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock() // 加锁保证并发安全
		if err := p.err.Load(); err != nil {
			// 父已取消：立即取消刚创建的 cancelCtx实例
			child.cancel(false, err.(error), p.cause)
		} else {
			// 父未取消：将子Context加入父的children列表 当父取消时会遍历取消子
			if p.children == nil {  //第一次需要创建
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
		return   // 大多数到这返回了
	}

	// 特殊情况：父Context实现了afterFuncer接口（比如WithDeadline/WithTimeout）
	if a, ok := parent.(afterFuncer); ok {
		c.mu.Lock()
		// 注册回调：父取消时，触发子取消
		stop := a.AfterFunc(func() {
			child.cancel(false, parent.Err(), Cause(parent))
		})
		// 包装父Context，记录stop函数（子取消时可停止回调）
		c.Context = stopCtx{
			Context: parent,
			stop:    stop,
		}
		c.mu.Unlock()
		return
	}

	// 兜底方案：父Context不是cancelCtx也不是afterFuncer，但可取消（Done()≠nil） 自己实现的 Context
	goroutines.Add(1)
	go func() {
		select {
		case <-parent.Done():
			// 父取消时，取消子Context
			child.cancel(false, parent.Err(), Cause(parent))
		case <-child.Done():
			// 子先取消，无需处理
		}
	}()
}
```

### 核心函数 `context.Cause` 

`Cause` 是用户获取 Context 取消原因的 **唯一入口**， `errgroup` 在 `wait` 判断时，可以拿到第一个错误，你也可以调用这个函数拿到错误，阅读源码我们知道，`errgroup` 内部主动 `cancel(err)` 传入了错误，并存储在了 `cause` 当然可以用该函数拿出

后续可以通过 `cause := context.Cause(ctx)`  拿到错误 

```go
func Cause(c Context) error {
    // 每个 Context 都有 Value(key any) any 这个函数 
	if cc, ok := c.Value(&cancelCtxKey).(*cancelCtx); ok {
		cc.mu.Lock()
		cause := cc.cause
		cc.mu.Unlock()
		if cause != nil {
			return cause
		}
	}

	return c.Err()
}

```

### `Value` 函数

其实就是递归向上找，我们这种情况直接返回，需要想上找的情况是自己存的业务 `key` ，后面会举例

```GO
func (c *cancelCtx) Value(key any) any {
	if key == &cancelCtxKey {
		return c
	}
	return value(c.Context, key)
}

// switch 所有的 已知的 context 找不到继续 向上找
func value(c Context, key any) any {
	for {
		switch ctx := c.(type) {
		case *valueCtx:
			if key == ctx.key {
				return ctx.val
			}
			c = ctx.Context
		case *cancelCtx:
			if key == &cancelCtxKey {
				return c
			}
			c = ctx.Context
		case withoutCancelCtx:
			if key == &cancelCtxKey {
				return nil
			}
			c = ctx.c
		case *timerCtx:
			if key == &cancelCtxKey {
				return &ctx.cancelCtx
			}
			c = ctx.Context
		case backgroundCtx, todoCtx:  // 空 emptyCtx
			return nil
		default:
			return c.Value(key)
		}
	}
}
```

### 自定义业务 key

```GO
// 自定义业务 key
type userKey struct{}

func main() {
    // 根 context，存入业务数据
    ctx1 := context.WithValue(context.Background(), userKey{}, "张三")
    
    // 在它基础上创建 cancelCtx
    ctx2, cancel := context.WithCancel(ctx1)
    
    // 再套一层 cancelCtx
    ctx3, cancel2 := context.WithCancel(ctx2)
    
    // 现在的结构：
    // Background (空)
    //    └── ctx1 (valueCtx, 存了 userKey="张三")
    //           └── ctx2 (cancelCtx)
    //                  └── ctx3 (cancelCtx)  ← 我们拿着这个
    
    // 查询业务数据：ctx3 本身没有，会递归向上找
    user := ctx3.Value(userKey{})
    fmt.Println("找到用户:", user)  // 输出: 找到用户: 张三
    
    // 即使 cancel 了，依然能找到
    cancel2()
    cancel()
    
    user2 := ctx3.Value(userKey{})
    fmt.Println("cancel后:", user2)  // 输出: cancel后: 张三
}
```



## 总结

`errgroup` 是对  `sync.WaitGroup ` 的增强封装，可以对 goroutine 数量限制 ，并且不需要显式的写`WaitGroup..Add(1)`，内部自动管理，防止漏写，出错自动 `cancel`

如果自己相加业务逻辑 直接可以复制代码，按自己的需求添加字段或者函数

对于 `context` 源代码，要搞清楚当 `cancel` 后切断的是父对子的连接，但子有个字段专门存放父，这个并没有切断，才使得我们可以调用 `context.Cause` 方法一直向上找