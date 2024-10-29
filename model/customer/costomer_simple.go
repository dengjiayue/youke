package model_customer

import (
	"fmt"

	"gorm.io/gorm"
)

type CostomerSimple struct {
	Id           int64  `gorm:"primaryKey;autoIncrement;column:id;comment:主键,自增;omitempty"`
	Name         string `gorm:"index:name_phone_idx;not null;column:name;type:varchar(20);comment:姓名;omitempty"`
	PhoneNumber  string `gorm:"index:name_phone_idx;not null;column:phone_number;type:varchar(15);comment:电话号码;omitempty"`
	IdcardNumber string `gorm:"uniqueIndex:idcard_number_uindex;not null;column:idcard_number;type:varchar(18);comment:身份证号;omitempty"`
}

// 自定义表名
func (CostomerSimple) TableName() string {
	return "customer" // 自定义表名
}

func ClumsOfSimble(db *gorm.DB) *gorm.DB {
	return db.Select("id", "name", "phone_number", "idcard_number")
}

func SelectCostomerSimple(db *gorm.DB, req *CostomerSimple) ([]CostomerSimple, int, error) {
	condition := ClumsOfSimble(db)
	ok := false
	if req.Name != "" {
		condition = condition.Where("name = ?", req.Name)
		ok = true
	}
	if req.PhoneNumber != "" && len(req.PhoneNumber) > 3 {
		condition = condition.Where("phone_number LIKE ?", "%"+req.PhoneNumber)
		ok = true
	}
	if req.IdcardNumber != "" && len(req.IdcardNumber) > 3 {
		condition = condition.Where("idcard_number LIKE ?", "%"+req.IdcardNumber)
		ok = true
	}
	if !ok {
		return nil, 0, fmt.Errorf("查询条件有误,名字,身份证,电话号码至少需要一个条件,并且身份证与电话号码至少需要4位尾数")
	}
	var resp []CostomerSimple
	err := condition.Find(&resp).Error
	if err != nil {
		return nil, 0, err
	}
	return resp, len(resp), nil
}
