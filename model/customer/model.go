package model_customer

import (
	"fmt"
	"time"
	public_db_func "youke/model/public"

	"gorm.io/gorm"
)

type Model struct {
	Id           int64     `gorm:"primaryKey;autoIncrement;column:id;comment:主键,自增"`
	Name         string    `gorm:"index:name_phone_idx;not null;column:name;type:varchar(20);comment:姓名"`
	PhoneNumber  string    `gorm:"index:name_phone_idx;not null;column:phone_number;type:varchar(15);comment:电话号码"`
	Birthday     time.Time `gorm:"not null;column:birthday;type:date;comment:生日"`
	FaceImg      string    `gorm:"not null;column:face_img;type:varchar(255);comment:头像"`
	IdCardImg    string    `gorm:"not null;column:id_card_img;type:varchar(255);comment:身份证"`
	IdcardNumber string    `gorm:"uniqueIndex:idcard_number_uindex;not null;column:idcard_number;type:varchar(20);comment:身份证号"`
	Address      string    `gorm:"not null;column:address;type:varchar(255);comment:地址"`
	Age          int       `gorm:"not null;column:age;comment:年龄"`
	GuardianId   int64     `gorm:"not null;column:guardian_id;comment:监护人Id"`
	UpdatedAt    time.Time `gorm:"not null;column:updated_at;autoUpdateTime;comment:更新时间"`
	CreatedAt    time.Time `gorm:"not null;column:created_at;autoCreateTime;comment:创建时间"`
}

// 自定义表名
func (Model) TableName() string {
	return "customer" // 自定义表名
}

// 自动建表
func CreateTable(db *gorm.DB) error {
	// 自动迁移，创建表并加注释
	err := db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='客户信息表'").AutoMigrate(&Model{})
	if err != nil {
		return err
	}
	return nil
}

func SelectById(db *gorm.DB, id int64) (*Model, error) {
	resp := new(Model)
	err := db.Where("id = ?", id).First(&resp).Error
	return resp, err
}

func (m *Model) Create(db *gorm.DB) error {
	if !public_db_func.CheckPhoneNumber(m.PhoneNumber) {
		return fmt.Errorf("电话号码有误,请检查")
	}
	if !public_db_func.CheckIDCard(m.IdcardNumber) {
		return fmt.Errorf("身份证号码有误,请重新拍照识别")
	}
	return db.Create(m).Error
}
