package main

import (
	"fmt"
	"youke/global/config"
)

func main() {
	data, err := config.LoadConfig("config.yml")
	if err != nil {
		fmt.Printf("err=%#v\n", err)
		return
	}
	fmt.Printf("config=%#v\n", data.Logger)
}
