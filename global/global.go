package global

import (
	"fmt"
	"youke/global/config"
	cos2 "youke/global/cos"
	"youke/global/database"
	"youke/global/logger"

	"github.com/spf13/pflag"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/gorm"
)

type GlobalData struct {
	Config *config.Config
	Logger *logger.Logger
	Db     *gorm.DB
	Cos    *cos.Client
}

var Global = &GlobalData{}

func InitGlobal() {
	fmt.Printf("服务初始化....")
	//解析配置文件
	var configPath = pflag.StringP("config", "c", "config.yml", "配置文件路径")
	var err error
	Global.Config, err = config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}
	//初始化日志
	Global.Logger, err = logger.NewLogger(Global.Config.Logger)
	if err != nil {
		panic(err)
	}
	//初始化MySQL
	Global.Db, err = database.NewDB(Global.Config.Mysql)
	if err != nil {
		panic(err)
	}
	//初始化cos
	Global.Cos, err = cos2.NewCosClient(*Global.Config.Cos)
	if err != nil {
		panic(err)
	}
	fmt.Printf("初始化完成....")
}
