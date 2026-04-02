package f1

import (
	"fmt"
	"unsafe"
)

// ------------------------------
// Go 底层 interface{} 结构模拟（简化版）
// 仅用于学习原理，生产环境禁止滥用
// ------------------------------

// eface 空接口底层结构
type eface struct {
	_type *_type         // 类型元数据指针
	data  unsafe.Pointer // 数据指针
}

// _type 类型元数据结构体
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
	str        NameOff
	ptrToThis  TypeOff
}

// 辅助类型定义
type tflag uint8
type NameOff int32
type TypeOff int32

// Kind 类型枚举（模拟 reflect.Kind）
type Kind uint8

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

// String 输出类型名称
func (k Kind) String() string {
	switch k {
	case Int:
		return "int"
	case String:
		return "string"
	case Struct:
		return "struct"
	case Ptr:
		return "ptr"
	default:
		return fmt.Sprintf("Kind(%d)", k)
	}
}

// MyStruct 测试用结构体
type MyStruct struct {
	Name string
	Age  int
	Tags []string
}

func main() {
	fmt.Println("--- Go interface{} (eface) 底层原理分析 ---")

	// 示例1：int 类型
	fmt.Println("\n=== 示例 1: int 赋值给 interface{} ===")
	var i interface{}
	num := 123
	i = num
	printEfaceDetails("num(int)", i)

	// 示例2：string 类型
	fmt.Println("\n=== 示例 2: string 赋值给 interface{} ===")
	var j interface{}
	str := "Hello, Go interfaces!"
	j = str
	printEfaceDetails("str(string)", j)

	// 示例3：*int 指针类型
	fmt.Println("\n=== 示例 3: *int 赋值给 interface{} ===")
	var k interface{}
	ptrNum := &num
	k = ptrNum
	printEfaceDetails("ptrNum(*int)", k)

	// 示例4：结构体类型
	fmt.Println("\n=== 示例 4: MyStruct 赋值给 interface{} ===")
	var l interface{}
	myS := MyStruct{Name: "Alice", Age: 30, Tags: []string{"developer", "go"}}
	l = myS
	printEfaceDetails("myS(MyStruct)", l)

	// 示例5：结构体指针类型
	fmt.Println("\n=== 示例 5: *MyStruct 赋值给 interface{} ===")
	var m interface{}
	ptrMyS := &myS
	m = ptrMyS
	printEfaceDetails("ptrMyS(*MyStruct)", m)
}

// printEfaceDetails 打印空接口底层结构信息
func printEfaceDetails(varName string, itf interface{}) {
	fmt.Printf("处理变量: %s\n", varName)

	// 强转成底层 eface 结构体（unsafe 操作）
	e := (*eface)(unsafe.Pointer(&itf))

	fmt.Printf("  eface 结构体地址: %p\n", unsafe.Pointer(&itf))
	fmt.Printf("  eface._type: %p\n", e._type)
	fmt.Printf("  eface.data: %p\n", e.data)

	if e._type == nil {
		fmt.Println("  接口为 nil，无类型与值")
		return
	}

	// 打印类型元数据
	fmt.Printf("  --- _type 内部信息 ---\n")
	fmt.Printf("    size: %d bytes\n", e._type.size)
	fmt.Printf("    ptrdata: %d bytes\n", e._type.ptrdata)
	fmt.Printf("    hash: 0x%x\n", e._type.hash)
	fmt.Printf("    tflag: %d\n", e._type.tflag)
	fmt.Printf("    align: %d bytes\n", e._type.align)
	fmt.Printf("    fieldAlign: %d bytes\n", e._type.fieldAlign)
	fmt.Printf("    kind: %s (%d)\n", Kind(e._type.kind).String(), e._type.kind)

	// 解析实际值
	fmt.Printf("  --- 动态值解析 ---\n")
	kind := Kind(e._type.kind)
	switch kind {
	case Int:
		val := *(*int)(e.data)
		fmt.Printf("    int 值: %d\n", val)

	case String:
		val := *(*string)(e.data)
		fmt.Printf("    string 值: %q\n", val)

	case Ptr:
		valPtr := *(*unsafe.Pointer)(e.data)
		switch e._type.kind {
		case uint8(Int):
			val := *(*int)(valPtr)
			fmt.Printf("    *int 指向值: %d，指针地址: %p\n", val, valPtr)
		default:
			val := *(*MyStruct)(valPtr)
			fmt.Printf("    *MyStruct 指向值: %+v，指针地址: %p\n", val, valPtr)
		}

	case Struct:
		val := *(*MyStruct)(e.data)
		fmt.Printf("    struct 值: %+v\n", val)

	default:
		fmt.Printf("    暂不支持解析该类型: %s\n", kind.String())
	}

	fmt.Println("--------------------------------------------------")
}

// getTypeNameOffset 模拟函数，仅用于编译通过
func getTypeNameOffset(_ string) NameOff {
	return 0
}
