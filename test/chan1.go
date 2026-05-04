package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int) // 无缓冲

	// 启动两个读取协程
	go func() {
		fmt.Println("协程1 开始等待...")
		val := <-ch
		fmt.Println("协程1 收到:", val)
	}()

	go func() {
		fmt.Println("协程2 开始等待...")
		val := <-ch
		fmt.Println("协程2 收到:", val)
	}()

	time.Sleep(time.Second) // 让两个协程都进入阻塞状态
	fmt.Println("准备发送第一个数据")

	ch <- 100 // 只有一个协程能收到
	time.Sleep(time.Second)

	fmt.Println("准备发送第二个数据")
	ch <- 200 // 另一个协程收到
	time.Sleep(time.Second)
}
