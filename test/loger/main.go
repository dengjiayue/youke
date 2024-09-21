package main

import (
	"fmt"
	"youke/global/logger"
)

func main() {
	l, err := logger.NewLogger(&logger.LoggerConfig{OutPath: "logger_test"})
	if err != nil {
		fmt.Printf("err=%#v\n", err)
		return
	}
	l.Error("错误")
	l.Warning("警告")
	l.Info("信息")
}
