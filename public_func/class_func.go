package public_func

import (
	jsoniter "github.com/json-iterator/go"
)

// 结构体快速映射
// to : 传递指针变量
func StructToStruct(val, to interface{}) error {
	data, err := jsoniter.Marshal(val)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(data, to)
	if err != nil {
		return err
	}

	return nil
}
