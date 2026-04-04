package main

import "fmt"

func B() (b int) {
	b = 2
	defer func() {
		println("defer1")
	}()
	defer func() { // 这个 defer 离 panic 最近，所以先执行
		if r := recover(); r != nil { // recover 在 defer 中被调用
			println("defer2: recovered with", r) // 会捕获 panic(1) 的值
		}
	}()
	panic(1) // 这里发生 panic
	return 1 // 这行代码不会被执行
}
func main() {
	n := 1
	fmt.Println(n)
	fmt.Println(B())

}
