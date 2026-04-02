package main

import (
	"fmt"
	"internal/abi"
	"unsafe"
)

// 定义一个简单的结构体
type MyStruct struct {
	Name string
	Age  int
}

func main() {
	// 示例1: 将 int 类型赋值给 interface{}
	var i interface{}
	num := 42
	i = num

	fmt.Printf("--- 示例1: int 赋值给 interface{} ---")
	printEface(i)

	// 示例2: 将 MyStruct 类型赋值给 interface{}
	var j interface{}
	ms := MyStruct{Name: "Alice", Age: 30}
	j = ms

	fmt.Printf("--- 示例2: MyStruct 赋值给 interface{} ---")
	printEface(j)

	// 示例3: 将指针类型赋值给 interface{}
	var k interface{}
	ptrNum := &num // 指向 num 的指针
	k = ptrNum

	fmt.Printf("--- 示例3: *int 赋值给 interface{} ---")
	printEface(k)
}

// 辅助函数，用于打印 eface 的底层结构
func printEface(itf interface{}) {
	// 将 interface{} 变量强制转换为 eface 结构体指针
	// 注意：这是 unsafe 操作，仅用于学习底层原理
	e := (*eface)(unsafe.Pointer(&itf))

	// 打印 _type 字段
	fmt.Printf("eface._type (type descriptor pointer): %p", e._type)
	// 尝试通过 _type 获取类型名称 (需要 reflect 包)
	fmt.Printf("  -> Dynamic Type Name: %s", e._type.String()) // 依赖 reflect.Type.String()
	// 打印 _type 结构体的部分字段
	fmt.Printf("  -> Type.Size_: %d bytes", e._type.Size_)
	fmt.Printf("  -> Type.Kind_: %s (enum %d)", e._type.Kind_.String(), e._type.Kind_)

	// 打印 data 字段
	fmt.Printf("eface.data (value pointer): %p", e.data)

	// 尝试解引用 data 字段，获取实际值
	// 注意：需要知道原始类型才能正确解引用
	switch e._type.String() { // 根据类型名判断，更严谨应使用 Type.Kind() 或 Type.Name()
	case "int":
		val := *(*int)(e.data)
		fmt.Printf("  -> Dynamic Value (int): %d", val)
	case "main.MyStruct":
		val := *(*MyStruct)(e.data)
		fmt.Printf("  -> Dynamic Value (MyStruct): %+v", val)
	case "*int":
		val := *(*int)(e.data) // data 指向的是 int 的地址，所以解引用得到 int
		fmt.Printf("  -> Dynamic Value (*int, pointed to value): %d", val)
	default:
		fmt.Printf("  -> Dynamic Value (unknown type, raw pointer data): %v", e.data)
	}
}

// 这是 Go 运行时内部的 _type 和 Kind 定义，为方便分析而复制到此处
// 实际在 src/runtime/runtime2.go 和 src/reflect/type.go
type _type struct {
	Size_       uintptr
	PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash        uint32  // hash of type; avoids computation in hash tables
	TFlag       TFlag   // extra type information flags
	Align_      uint8   // alignment of variable with this type
	FieldAlign_ uint8   // alignment of struct field with this type
	Kind_       Kind    // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal     func(unsafe.Pointer, unsafe.Pointer) bool
	GCData    *byte
	Str       NameOff // string form
	PtrToThis TypeOff // type for pointer to this type, may be zero
}

type TFlag uint8
type NameOff int32
type TypeOff int32

// Kind represents the specific kind of type that a Type represents.
// The zero value is not a valid Kind.
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

// Kind.String() 方法 (简化版，实际在 reflect 包中)
func (k Kind) String() string {
	if k < Kind(len(kindNames)) {
		return kindNames[k]
	}
	return "Kind(" + fmt.Sprint(k) + ")"
}

var kindNames = []string{
	Invalid:       "invalid",
	Bool:          "bool",
	Int:           "int",
	Int8:          "int8",
	Int16:         "int16",
	Int32:         "int32",
	Int64:         "int64",
	Uint:          "uint",
	Uint8:         "uint8",
	Uint16:        "uint16",
	Uint32:        "uint32",
	Uint64:        "uint64",
	Uintptr:       "uintptr",
	Float32:       "float32",
	Float64:       "float64",
	Complex64:     "complex64",
	Complex128:    "complex128",
	Array:         "array",
	Chan:          "chan",
	Func:          "func",
	Interface:     "interface",
	Map:           "map",
	Ptr:           "ptr",
	Slice:         "slice",
	String:        "string",
	Struct:        "struct",
	UnsafePointer: "unsafe.Pointer",
}

// 模拟 _type 的 String() 方法，实际在 reflect.Type 中实现
func (t *_type) String() string {
	// 在实际运行时，这个 Str 字段是一个偏移量，指向类型名称字符串
	// 这里为了简化演示，我们直接用 Kind() 结合其他信息
	// 真实的 Type.String() 是一个复杂的查找过程
	if t == nil {
		return "<nil>"
	}
	// 更严谨的做法是解析 t.Str (NameOff)
	// 但这需要访问运行时内部的字符串表，超出了本示例的范围
	// 这里直接使用 Kind() 作为近似表示
	return t.Kind_.String() // 这样会打印 "int", "struct", "ptr"
}

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}
type itab struct {
	Inter *abi.InterfaceType
	Type  *_type
	Hash  uint32     // copy of Type.Hash. Used for type switches.
	Fun   [1]uintptr // variable sized. fun[0]==0 means Type does not implement Inter.
}
