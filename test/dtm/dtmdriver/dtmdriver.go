package dtmdriver

import "fmt"

// 1. 全局变量定义
// drivers：存储所有已注册的驱动，key=驱动名，value=驱动实例
var (
	drivers = map[string]Driver{}
	// current：当前正在使用的驱动（全局唯一）
	current Driver
)

// Register 驱动注册函数（驱动开发者调用）
// 作用：把自定义驱动注册到DTM的驱动列表中
// 参数：driver - 实现了DTM Driver接口的自定义驱动实例
func Register(driver Driver) {
	// 以驱动名作为key，存入全局map
	drivers[driver.GetName()] = driver
}

// Use 驱动使用函数（业务开发者调用）
// 作用：指定DTM使用哪个驱动（比如指定使用"dtm-driver-kratos"驱动）
// 参数：name - 驱动名（对应driver.GetName()的返回值）
func Use(name string) error {
	// 从全局map中获取指定名称的驱动
	v := drivers[name]
	// 校验1：驱动未注册则返回错误
	if v == nil {
		return fmt.Errorf("no dtm driver with name: %s has been registered", name)
	}
	// 校验2：禁止重复使用不同驱动（避免冲突）
	if current != nil && v != current {
		return fmt.Errorf("use has been called previously with name: %s different from: %s", current.GetName(), name)
	}
	// 校验3：如果是首次使用该驱动，初始化驱动的地址解析器
	if current == nil {
		current = v
		// 调用驱动的RegisterAddrResolver()，注册地址解析逻辑（比如让DTM识别discovery://协议）
		v.RegisterAddrResolver()
	}
	return nil
}

// GetDriver 获取取当前使用的驱动（DTM内部调用）
// 作用：DTM执行服务调用时，获取当前激活的驱动，用它解析地址/注册服务
func GetDriver() Driver {
	// 如果未指定驱动，使用默认驱动（DTM内置的default驱动）
	if current == nil {
		return drivers["default"]
	}
	// 否则返回当前激活的驱动
	return current
}
