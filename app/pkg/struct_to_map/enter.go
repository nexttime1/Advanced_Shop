package struct_to_map

import (
	"github.com/fatih/structs"
	"reflect"
)

func StructToMap(obj interface{}) map[string]interface{} {
	m := structs.Map(obj)
	data := DeleteEmpty(m)
	return data
}

func DeleteEmpty(m map[string]interface{}) map[string]interface{} {
	var data = make(map[string]interface{}, 0)
	for key, v := range m {
		switch val := v.(type) {
		case string:
			if val != "" {
				data[key] = val
			}
		case int:
			if val != 0 {
				data[key] = val
			}
		case uint:
			if val != 0 {
				data[key] = val
			}
		case int32:
			if val != 0 {
				data[key] = val
			}
		case int64:
			if val != 0 {
				data[key] = val
			}
		case float32:
			if val != 0 {
				data[key] = val
			}
		case float64:
			if val != 0 {
				data[key] = val
			}
		case []string:
			if val != nil {
				data[key] = val
			}

		//case enum.Array:
		//	if len(val) != 0 {
		//		data[key] = val
		//	}
		default:
			// 看看指针为空么 为空取消赋值
			value := reflect.ValueOf(v)
			if value.Kind() == reflect.Ptr {
				// 如果是指针且指向为空，则跳过
				if value.IsNil() {
					continue
				}
			}
			data[key] = v
		}

	}
	return data
}
