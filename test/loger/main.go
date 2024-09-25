package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type StructA struct {
	Name  *string
	Age   *int
	Email *string
}

type StructB struct {
	Name  string
	Age   int
	Email string
}

// func main() {
// 	name := "Alice"
// 	age := 30
// 	email := "alice@example.com"
// 	a := StructA{
// 		Name:  &name,
// 		Age:   &age,
// 		Email: &email,
// 	}

// 	b := new(StructB)

// 	// 调用 StructToStruct 函数
// 	err := StructToStruct(a, &b)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	// 输出转换后的 StructB
// 	fmt.Printf("StructB: %+v\n", b)
// }

func main() {
	// l, err := logger.NewLogger(&logger.LoggerConfig{OutPath: "logger_test"})
	// if err != nil {
	// 	fmt.Printf("err=%#v\n", err)
	// 	return
	// }
	// l.Error("错误")
	// l.Warning("警告")
	// l.Info("信息")
	a, b, c := 1, true, "123"
	fmt.Println(a, b, c)
	v1 := &A{
		Name: a,
		// IsOk:  b,
		Email: c,
	}
	data, err := jsoniter.MarshalToString(v1)
	fmt.Printf("data=%#v\n", data)
	v2 := new(B)
	err = StructToStruct(v1, v2)
	fmt.Printf("data=%#v\n;err=%v\n", *v2, err)
}

type A struct {
	Name  int    `json:"name"`
	IsOk  bool   `json:"isOk"`
	Email string `json:"email"`
}

type B struct {
	Name  *int    `json:"name"`
	IsOk  *bool   `json:"isOk"`
	Email *string `json:"email"`
}

func StructToStruct(val, to interface{}) error {

	// 序列化源结构体
	data, err := jsoniter.Marshal(val)
	if err != nil {
		return err
	}
	fmt.Printf("data=%s\n", data)

	// 直接将反序列化目标传递到 'to'，而不是 '&to'
	err = jsoniter.Unmarshal(data, to) // 这里去掉 & 符号
	if err != nil {
		return err
	}

	return nil
}
