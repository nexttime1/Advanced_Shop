package main

import (
	"fmt"
)

//func stringVar[T mystruct ](t T) string {
//	t.String()
//}

func main() {

	s := []int{0, 1, 2, 3, 4}
	s = append(s[:2], s[3:]...)
	fmt.Println("s: %v, len: %d, cap: %d", s, len(s), cap(s))
	v := s[4]
	// 是否会数组访问越界
	fmt.Println(v)

}

func changeSlice(s1 []int) {
	s1 = append(s1, 10)
}
