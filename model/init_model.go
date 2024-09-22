package model

import (
	model_customer "youke/model/customer"
	model_order "youke/model/order"

	"gorm.io/gorm"
)

func InitDataBaseModel(db *gorm.DB) {
	err := model_customer.CreateTable(db)
	if err != nil {
		panic(err)
	}
	err = model_order.CreateTable(db)
	if err != nil {
		panic(err)
	}
	// ...
}
