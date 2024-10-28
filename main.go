package main

import (
	"fmt"
	"youke/global"
	"youke/model"
	"youke/src/router"
)

func Init() {
	//初始化全局变量
	global.InitGlobal()
	//初始化/迁移表结构
	model.InitDataBaseModel(global.Global.Db)
}

func main() {
	Init()

	fmt.Printf("--------------------------\n-------server start-------\n--------------------------")

	router.Run()
}
