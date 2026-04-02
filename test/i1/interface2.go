package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ===================== Go 运行时底层结构体定义（模拟源码） =====================
// 对应 runtime/type.go 中的基础类型结构体
type tflag uint8
type nameOff int32
type typeOff int32

// _type 是所有类型的基础描述结构体
type _type struct {
	size       uintptr
	ptrdata    uintptr
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	equal      func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata     *byte
	str        nameOff
	ptrToThis  typeOff
}

// Kind 类型枚举（模拟）
type Kind uint

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)

func (k Kind) String() string {
	switch k {
	case Interface:
		return "interface"
	case Struct:
		return "struct"
	default:
		return "unknown"
	}
}

// 接口方法结构体
type imethod struct {
	name nameOff
	ityp typeOff
}

// 接口类型元信息
type interfaceType struct {
	_type   _type
	pkgPath name
	methods []imethod
}

// 空接口结构体 (eface)：无方法的 interface{} 底层结构
type eface struct {
	_type *_type
	data  unsafe.Pointer
}

// 方法表结构体 (itab)：带方法接口的核心映射表
type itab struct {
	inter *interfaceType
	_type *_type
	hash  uint32
	_     [4]byte
	fun   [1]uintptr
}

// 带方法接口结构体 (iface)：有方法的接口底层结构
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

// 占位符结构体
type name struct{}

// ===================== 业务代码：接口与实现 =====================
// Stringer 自定义带方法接口
type Stringer interface {
	String() string
	ID() int
}

// MyData 实现 Stringer 接口的结构体
type MyData struct {
	Value string
	IDVal int
}

// String 实现 Stringer 接口方法
func (md MyData) String() string {
	return fmt.Sprintf("MyData{Value: %q, ID: %d}", md.Value, md.IDVal)
}

// ID 实现 Stringer 接口方法
func (md MyData) ID() int {
	return md.IDVal
}

// OnlyStringer 子集接口（用于接口间转换演示）
type OnlyStringer interface {
	String() string
	ID() int
}

// ===================== 核心演示函数 =====================
func main() {
	fmt.Println("--- Go 带方法接口 (iface) 底层原理分析 ---")

	// 1. 具体类型赋值给带方法接口 (iface 首次创建)
	fmt.Println("\n=== 1. MyData 类型赋值给 Stringer 接口 ===")
	var s Stringer
	data1 := MyData{Value: "hello", IDVal: 100}
	s = data1
	printIfaceDetails("s (Stringer, data1)", s)

	// 2. 再次赋值相同动态类型 (验证 itab 缓存)
	fmt.Println("\n=== 2. 再次赋值相同动态类型 (验证 itab 缓存) ===")
	var s2 Stringer
	data2 := MyData{Value: "world", IDVal: 200}
	s2 = data2
	printIfaceDetails("s2 (Stringer, data2)", s2)
	fmt.Printf("  itab for s (%p) and s2 (%p) are %s\n",
		(*iface)(unsafe.Pointer(&s)).tab,
		(*iface)(unsafe.Pointer(&s2)).tab,
		ternary(
			(*iface)(unsafe.Pointer(&s)).tab == (*iface)(unsafe.Pointer(&s2)).tab,
			"THE SAME (itab cached)",
			"DIFFERENT (something is wrong or not cached)",
		),
	)

	// 3. 接口方法调用
	fmt.Println("\n=== 3. 接口方法调用 (s.String() 和 s.ID()) ===")
	fmt.Printf("  s.String() => %s\n", s.String())
	fmt.Printf("  s.ID() => %d\n", s.ID())
	explainMethodCall(s)

	// 4. 类型断言
	fmt.Println("\n=== 4. 接口类型断言 ===")
	// 断言成功
	if md, ok := s.(MyData); ok {
		fmt.Printf("  类型断言成功: s 实际上是 MyData 类型，值为 %+v\n", md)
		explainTypeAssertion(s, reflect.TypeOf(MyData{}))
	} else {
		fmt.Println("  类型断言失败: s 不是 MyData 类型")
	}
	// 断言失败
	if _, ok := s.(fmt.Stringer); ok {
		fmt.Printf("  类型断言成功: s 实际上是 fmt.Stringer 类型\n")
	} else {
		fmt.Printf("  类型断言失败: s 不是 fmt.Stringer 类型\n")
		explainTypeAssertion(s, reflect.TypeOf((*fmt.Stringer)(nil)).Elem())
	}

	// 5. 空接口转带方法接口
	fmt.Println("\n=== 5. interface{} 到 Stringer 的转换 ===")
	var emptyI interface{}
	emptyI = MyData{Value: "converted", IDVal: 300}
	if convS, ok := emptyI.(Stringer); ok {
		fmt.Printf("  interface{} -> Stringer 转换成功: %s\n", convS.String())
		explainEfaceToIfaceConversion(emptyI, convS)
	} else {
		fmt.Println("  interface{} -> Stringer 转换失败")
	}

	// 6. 接口间赋值（子集接口）
	fmt.Println("\n=== 6. 接口到接口的赋值 (子集接口) ===")
	var os OnlyStringer
	os = s
	printIfaceDetails("os (OnlyStringer, s)", os)
	fmt.Printf("  itab for s (%p) and os (%p) are %s\n",
		(*iface)(unsafe.Pointer(&s)).tab,
		(*iface)(unsafe.Pointer(&os)).tab,
		ternary(
			(*iface)(unsafe.Pointer(&s)).tab == (*iface)(unsafe.Pointer(&os)).tab,
			"THE SAME (iface reuse)",
			"DIFFERENT (new itab generated)",
		),
	)
	explainIfaceToIfaceConversion(s, os)
}

// printIfaceDetails 打印接口底层结构体详情
func printIfaceDetails(varName string, itf Stringer) {
	fmt.Printf("处理变量: %s\n", varName)
	i := (*iface)(unsafe.Pointer(&itf))

	fmt.Printf("  iface 结构体地址: %p\n", unsafe.Pointer(&itf))
	fmt.Printf("  iface.tab (itab 指针): %p\n", i.tab)
	fmt.Printf("  iface.data (动态值指针): %p\n", i.data)

	if i.tab == nil {
		fmt.Println("  iface.tab 为 nil，接口为 nil 值。")
		return
	}

	// 打印 itab 内部字段
	fmt.Printf("  --- itab 内部字段 (地址: %p) --- \n", i.tab)
	fmt.Printf("    itab.inter (interfaceType 指针): %p\n", i.tab.inter)
	fmt.Printf("    itab._type (_type 指针): %p\n", i.tab._type)
	fmt.Printf("    itab.hash: 0x%x\n", i.tab.hash)

	// 打印接口元信息
	if i.tab.inter != nil {
		fmt.Printf("    --- itab.inter (interfaceType %p) 内部字段 --- \n", i.tab.inter)
		fmt.Printf("      interfaceType._type.size: %d bytes\n", i.tab.inter._type.size)
		fmt.Printf("      interfaceType._type.kind: %s (原始值: %d)\n", Kind(i.tab.inter._type.kind).String(), i.tab.inter._type.kind)
		fmt.Printf("      interfaceType.methods: %d methods\n", len(i.tab.inter.methods))
	}

	// 打印动态类型信息
	if i.tab._type != nil {
		fmt.Printf("    --- itab._type (_type %p) 内部字段 --- \n", i.tab._type)
		fmt.Printf("      _type.size: %d bytes\n", i.tab._type.size)
		fmt.Printf("      _type.ptrdata: %d bytes\n", i.tab._type.ptrdata)
		fmt.Printf("      _type.hash: 0x%x\n", i.tab._type.hash)
		fmt.Printf("      _type.kind: %s (原始值: %d)\n", Kind(i.tab._type.kind).String(), i.tab._type.kind)
	}

	// 打印方法指针数组
	fmt.Printf("    itab.fun (方法函数指针):\n")
	if i.tab.fun[0] != 0 {
		fmt.Printf("      fun[0] (String() 方法指针): %p\n", i.tab.fun[0])
	}
	if len(i.tab.inter.methods) > 1 {
		secondFuncPtr := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&i.tab.fun[0])) + unsafe.Sizeof(uintptr(0))))
		fmt.Printf("      fun[1] (ID() 方法指针): %p\n", secondFuncPtr)
	}

	// 解析动态值
	fmt.Printf("  --- 动态值解析 --- \n")
	val := *(*MyData)(i.data)
	fmt.Printf("    动态值 (MyData): %+v\n", val)
	fmt.Println("--------------------------------------------------")
}

// explainMethodCall 解释接口方法调用流程
func explainMethodCall(itf Stringer) {
	fmt.Println("\n-- 解释接口方法调用 (s.String()) --")
	i := (*iface)(unsafe.Pointer(&itf))
	fmt.Printf("  1. 从 iface (%p) 获取 itab (%p) 和 data (%p)\n", unsafe.Pointer(&itf), i.tab, i.data)
	fmt.Printf("  2. 从 itab 获取接口方法定义与索引\n")
	fmt.Printf("  3. 通过索引从 itab.fun 获取方法指针 (%p)\n", i.tab.fun[0])
	fmt.Printf("  4. 将 data (%p) 作为接收者调用方法\n", i.data)
	fmt.Printf("  结果: 执行 MyData.String()\n")
}

// explainTypeAssertion 解释类型断言原理
func explainTypeAssertion(itf Stringer, targetType reflect.Type) {
	fmt.Printf("\n-- 解释类型断言 (s.(%s)) --\n", targetType.String())
	i := (*iface)(unsafe.Pointer(&itf))
	fmt.Printf("  1. 从 iface 获取 itab (%p)\n", i.tab)
	fmt.Printf("  2. 对比动态类型与目标类型的 _type 指针\n")
	fmt.Printf("  3. 指针一致则断言成功，否则失败\n")
	fmt.Printf("  4. 成功则提取 data (%p) 中的值\n", i.data)
}

// explainEfaceToIfaceConversion 解释空接口转带方法接口
func explainEfaceToIfaceConversion(efaceVar interface{}, ifaceVar Stringer) {
	fmt.Println("\n-- 解释 interface{} 到 Stringer 的转换 --")
	e := (*eface)(unsafe.Pointer(&efaceVar))
	i := (*iface)(unsafe.Pointer(&ifaceVar))
	fmt.Printf("  1. 源 eface: _type=%p, data=%p\n", e._type, e.data)
	fmt.Printf("  2. 运行时创建 itab(Stringer, 动态类型)\n")
	fmt.Printf("  3. 构建新 iface: tab=%p, data=%p\n", i.tab, i.data)
}

// explainIfaceToIfaceConversion 解释接口间转换
func explainIfaceToIfaceConversion(srcIface Stringer, destIface OnlyStringer) {
	fmt.Println("\n-- 解释 Stringer 到 OnlyStringer 的赋值 --")
	src := (*iface)(unsafe.Pointer(&srcIface))
	dest := (*iface)(unsafe.Pointer(&destIface))
	fmt.Printf("  1. 源 iface: itab=%p, data=%p\n", src.tab, src.data)
	fmt.Printf("  2. 为目标接口创建新 itab\n")
	fmt.Printf("  3. 目标 iface: itab=%p, data=%p\n", dest.tab, dest.data)
	fmt.Println("  注意：接口不同会生成新 itab，data 指针复用")
}

// ternary 三元运算符辅助函数
func ternary(cond bool, trueVal, falseVal string) string {
	if cond {
		return trueVal
	}
	return falseVal
}
