package main

import "fmt"

type mystruct struct {
	msg string
}

func (m mystruct) String() string {
	fmt.Println(m.msg)
	return m.msg
}

//func stringVar[T mystruct ](t T) string {
//	t.String()
//}

func main() {
	var a []int
	fmt.Println(len(a), cap(a))
	a = append(a, 1)
	fmt.Println(len(a), cap(a))
	a = append(a, []int{1, 2, 3, 4}...)
	fmt.Println(len(a), cap(a))

}
