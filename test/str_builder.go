package main

import (
	"fmt"
)

func main() {
	//strings.Builder{}
	s := "hello world"
	modify(s)
}

func modify(s string) {
	// s 是拷贝，但只拷贝了 header（16字节）
	// 底层字节数组没有被拷贝！
	fmt.Println(s)
}
