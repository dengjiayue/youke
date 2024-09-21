package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	Dns     string `yaml:"Dns"`
	MaxPool int    `yaml:"MaxPool"`
}

func NewDB(sql *Mysql) (*gorm.DB, error) {
	//查看仓库名
	fmt.Printf("链接数据库:%v\n", sql)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(sql.Dns), &gorm.Config{CreateBatchSize: sql.MaxPool})
	if err != nil {
		return nil, err
	}
	return db, nil
}
