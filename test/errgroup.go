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
		return errors.New("task1 error") // 业务错误：作为取消原因
	})

	eg.Go(func() error {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("doing task2")
			case <-ctx.Done():
				fmt.Println("task2 canceled")
				//  新增：调用 context.Cause()，获取取消的具体原因
				cause := context.Cause(ctx)
				fmt.Printf("task2 被取消的原因：%v\n", cause)
				fmt.Printf("ctx.err：%v\n", ctx.Err())
				return ctx.Err()
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
				// 🌟 新增：调用 context.Cause()，获取取消的具体原因
				cause := context.Cause(ctx)
				fmt.Printf("task3 被取消的原因：%v\n", cause)
				fmt.Printf("ctx.err：%v\n", ctx.Err())
				return ctx.Err()
			}
		}
	})

	err := eg.Wait()
	if err != nil {
		fmt.Println("task failed")
		// 🌟 可选新增：在 Wait() 后，也可以调用 context.Cause() 获取原因
		fmt.Printf("整体任务被取消的原因：%v\n", err)
	} else {
		fmt.Println("task success")
	}
}
