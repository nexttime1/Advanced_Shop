package common

type PageInfo struct {
	Limit int32  `form:"limit"`
	Page  int32  `form:"page"`
	Key   string `form:"key"`
	Sort  string `form:"sort"` //前端可以覆盖
}
