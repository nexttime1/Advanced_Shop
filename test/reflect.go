package main

import (
	"fmt"
	"reflect"
)

// 定义一个结构体 Coder
type Coder struct {
	Name string
}

// 为 *Coder 类型实现 fmt.Stringer 接口（接口要求实现 String() 方法）
func (c Coder) String() string {
	return c.Name + " 1111"
}

func main() {
	// 1. 创建一个 *Coder 类型的实例
	coder := Coder{Name: "xtm"}

	// 2. 获取 coder 的反射类型对象 reflect.Type
	typ := reflect.TypeOf(coder)
	fmt.Println("TypeOf coder:", typ)
	// 5. 打印 *Coder 类型的 Kind（底层类型）
	fmt.Println("kind of coder:", typ.Kind())
	// 3. 获取 coder 的反射值对象 reflect.Value
	val := reflect.ValueOf(coder)
	// 7. 打印 *Coder 的反射值对象
	fmt.Println("ValueOf coder:", val)
	// 4. 获取 fmt.Stringer 接口的反射类型（核心技巧）
	typeOfStringer := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	// 8. 判断 *Coder 类型是否实现了 fmt.Stringer 接口
	fmt.Println("implements stringer:", typ.Implements(typeOfStringer))

	var b fmt.Stringer
	b = coder
	// 🔥 修复：不要用 TypeOf(b).Elem()
	fmt.Println("TypeOf(b) 类型：", reflect.TypeOf(b))        // fmt.Stringer (接口)
	fmt.Println("TypeOf(b) 种类：", reflect.TypeOf(b).Kind()) // interface

	// ✅ 正确获取接口内部的动态类型（Coder）
	a := reflect.ValueOf(b).Type()
	fmt.Println("接口动态类型 a：", a)  // main.Coder
	fmt.Println("a 的种类：", a.Kind()) // struct
}
