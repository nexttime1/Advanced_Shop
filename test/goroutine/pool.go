package goroutine

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrPoolClosed = errors.New("pool closed")
	ErrQueueFull  = errors.New("task queue full")
	ErrNilTask    = errors.New("nil task")
)

type RejectPolicy int

const (
	RejectReturnError RejectPolicy = iota
	RejectBlock
	RejectBlockWithTimeout
	RejectCallerRuns
	RejectDiscard
)

/*
CPU 密集型任务
加密；压缩；图片处理；大量计算   建议和 CPU 核数一致
 IO 密集型任务
查数据库；发送 http；读写文件    不是看 CPU，而是看瓶颈资源。


*/

type Options struct {
	WorkerNum int
	QueueSize int // 队列大小  要根据 峰值流量；任务执行耗时；内存；下游承载能力。

	SubmitTimeout time.Duration
	TaskTimeout   time.Duration // 带超时等待 等待一段时间 还是没有空位 直接返回

	RejectPolicy RejectPolicy

	RecoverPanic bool
	PanicHandler func(any)

	ErrorHandler func(err error)

	WorkerName string
}

type Pool struct {
	opts Options

	taskCh chan Task

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	closed atomic.Bool

	submitted atomic.Int64
	running   atomic.Int64
	completed atomic.Int64
	failed    atomic.Int64
	rejected  atomic.Int64
	panicked  atomic.Int64
}

type Task func(ctx context.Context) error

func NewPool(opts Options) *Pool {

	// 合法判断
	if opts.WorkerNum <= 0 {
		opts.WorkerNum = 1
	}
	if opts.QueueSize < 0 {
		opts.QueueSize = 0
	}

	if opts.RejectPolicy == RejectBlockWithTimeout && opts.SubmitTimeout <= 0 {
		opts.SubmitTimeout = time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool{
		opts:   opts,
		taskCh: make(chan Task, opts.QueueSize),
		ctx:    ctx,
		cancel: cancel,
	}

	p.start()

	return p
}

func (p *Pool) start() {
	for i := 0; i < p.opts.WorkerNum; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return

		case task, ok := <-p.taskCh:
			if !ok {
				return
			}
			p.runTask(task)
		}
	}
}

func (p *Pool) runTask(task Task) {
	p.running.Add(1)
	defer p.running.Add(-1)

	start := time.Now()
	_ = start

	defer func() {
		if v := recover(); v != nil {
			p.panicked.Add(1)

			if p.opts.PanicHandler != nil {
				p.opts.PanicHandler(v)
			}
		}
	}()

	taskCtx := p.ctx
	var cancel context.CancelFunc

	if p.opts.TaskTimeout > 0 {
		taskCtx, cancel = context.WithTimeout(p.ctx, p.opts.TaskTimeout)
		defer cancel()
	}

	err := task(taskCtx)
	if err != nil {
		p.failed.Add(1)
		if p.opts.ErrorHandler != nil {
			p.opts.ErrorHandler(err)
		}
		return
	}

	p.completed.Add(1)
}

func (p *Pool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return ErrNilTask
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if p.closed.Load() {
		p.rejected.Add(1)
		return ErrPoolClosed
	}

	switch p.opts.RejectPolicy {
	case RejectReturnError:
		return p.submitNonBlock(ctx, task)

	case RejectBlock:
		return p.submitBlock(ctx, task)

	case RejectBlockWithTimeout:
		return p.submitBlockWithTimeout(ctx, task)

	case RejectCallerRuns:
		return p.submitCallerRuns(ctx, task)

	case RejectDiscard:
		return p.submitDiscard(ctx, task)

	default:
		return p.submitNonBlock(ctx, task)
	}
}

func (p *Pool) submitNonBlock(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		p.rejected.Add(1)
		return ctx.Err()

	case <-p.ctx.Done():
		p.rejected.Add(1)
		return ErrPoolClosed

	case p.taskCh <- task:
		p.submitted.Add(1)
		return nil

	default:
		p.rejected.Add(1)
		return ErrQueueFull
	}
}

func (p *Pool) submitBlock(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		p.rejected.Add(1)
		return ctx.Err()

	case <-p.ctx.Done():
		p.rejected.Add(1)
		return ErrPoolClosed

	case p.taskCh <- task:
		p.submitted.Add(1)
		return nil
	}
}

func (p *Pool) submitBlockWithTimeout(ctx context.Context, task Task) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, p.opts.SubmitTimeout)
	defer cancel()

	select {
	case <-timeoutCtx.Done():
		p.rejected.Add(1)
		return timeoutCtx.Err()

	case <-p.ctx.Done():
		p.rejected.Add(1)
		return ErrPoolClosed

	case p.taskCh <- task:
		p.submitted.Add(1)
		return nil
	}
}

func (p *Pool) submitCallerRuns(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		p.rejected.Add(1)
		return ctx.Err()

	case <-p.ctx.Done():
		p.rejected.Add(1)
		return ErrPoolClosed

	case p.taskCh <- task:
		p.submitted.Add(1)
		return nil

	default:
		p.submitted.Add(1)
		p.runTask(task)
		return nil
	}
}

func (p *Pool) submitDiscard(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		p.rejected.Add(1)
		return ctx.Err()

	case <-p.ctx.Done():
		p.rejected.Add(1)
		return ErrPoolClosed

	case p.taskCh <- task:
		p.submitted.Add(1)
		return nil

	default:
		p.rejected.Add(1)
		return nil
	}
}

func (p *Pool) Stop() {
	if p.closed.CompareAndSwap(false, true) {
		p.cancel()
		p.wg.Wait()
	}
}
