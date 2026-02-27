package main

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseServerMethod(uri string) (server string, method string, err error) {
	// 保留Kratos原版：处理直连地址（如127.0.0.1:8081/Inventory/Sell）
	fmt.Println("ParseServerMethod 调用  url: ", uri)
	if !strings.Contains(uri, "//") {
		sep := strings.IndexByte(uri, '/')
		if sep == -1 {
			return "", "", fmt.Errorf("bad url: '%s'. no '/' found", uri)
		}
		return uri[:sep], uri[sep:], nil
	}

	// 步骤1：修复Kratos bug1 - 解析失败返回具体错误（而非nil）
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", fmt.Errorf("parse consul discovery uri %s failed: %v", uri, err)
	}

	// 步骤2：修复Kratos bug2 - 正确拆分Path（适配Consul resolver格式）
	// 核心：把Kratos的index计算逻辑替换为更鲁棒的拆分方式，避免多斜杠
	cleanPath := strings.TrimPrefix(u.Path, "/") // 去掉Path开头的/
	pathParts := strings.SplitN(cleanPath, "/", 2)
	if len(pathParts) < 1 {
		return "", "", fmt.Errorf("consul discovery url %s missing service name", uri)
	}
	if len(pathParts) < 2 {
		return "", "", fmt.Errorf("consul discovery url %s missing method (e.g. /Inventory/Sell)", uri)
	}

	server = fmt.Sprintf("%s:///%s", u.Scheme, pathParts[0])

	// method 保持 /Inventory/Sell
	method = "/" + pathParts[1]
	return server, method, nil

}

func main() {
	path1 := "discovery:///xshop-inventory-srv/Inventory/Sell"
	server, m, err := ParseServerMethod(path1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s \n %s\n", server, m)
	// discovery:///xshop-inventory-srv
	//   /Inventory/Sell

}
