package main

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseServerMethod(uri string) (server string, method string, err error) {
	// 情况1：URI不含"//"（比如"127.0.0.1:8081/Inventory/Sell"）
	if !strings.Contains(uri, "//") {
		sep := strings.IndexByte(uri, '/') // 找第一个"/"的位置
		if sep == -1 {                     // 没有"/"，格式错误
			return "", "", fmt.Errorf("bad url: '%s'. no '/' found", uri)
		}
		// 拆分：server=127.0.0.1:8081，method=/Inventory/Sell
		return uri[:sep], uri[sep:], nil
	}

	// 情况2：URI含"//"（比如"discovery:///xshop-inventory-srv/Inventory/Sell"）
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", nil // 解析失败返回空（DTM会处理错误）
	}
	// 找Path中第一个"/"的位置（比如Path="/Inventory/Sell"，index=0）
	index := strings.IndexByte(u.Path[1:], '/') + 1
	// 拆分：
	// server = scheme://host + Path[:index] → discovery:///xshop-inventory-srv
	// method = Path[index:] → /Inventory/Sell
	return u.Scheme + "://" + u.Host + u.Path[:index], u.Path[index:], nil
}

func main() {
	path1 := "discovery:///xshop-inventory-srv/Inventory/Sell"
	server, m, err := ParseServerMethod(path1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n %s\n", server, m)

}
