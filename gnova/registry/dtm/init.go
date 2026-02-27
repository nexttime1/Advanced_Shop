package driver

import (
	"fmt"
	"github.com/dtm-labs/dtmdriver"
)

func init() {
	// dtm 服务注册
	Address := "192.168.163.132:8500"
	Scheme := "http"
	dtmDriver := NewXtmDriver(Address, Scheme)
	dtmdriver.Register(dtmDriver)
	// 调用的Use方法激活驱动
	err := dtmdriver.Use(DriverName)
	if err != nil {
		fmt.Println("激活DTM驱动失败: %v", err)
		panic(err)
	}
	fmt.Println("激活DTM驱动成功")
}
