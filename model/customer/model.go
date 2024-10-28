package model_customer

import (
	"fmt"
	"time"
	public_db_func "youke/model/public"
	"youke/public_func"

	"gorm.io/gorm"
)

type Model struct {
	Id           *int64  `gorm:"primaryKey;autoIncrement;column:id;comment:主键,自增;omitempty"`
	Name         *string `gorm:"index:name_idx;not null;column:name;type:varchar(20);comment:姓名;omitempty"`
	PhoneNumber  *string `gorm:"index:phone_number_idx;not null;column:phone_number;type:varchar(15);comment:电话号码;omitempty"`
	FaceImg      *string `gorm:"not null;column:face_img;type:varchar(255);comment:头像;omitempty"`
	IdcardImg    *string `gorm:"not null;column:idcard_img;type:varchar(255);comment:身份证;omitempty"`
	IdcardNumber *string `gorm:"uniqueIndex:idcard_number_uindex;not null;column:idcard_number;type:varchar(20);comment:身份证号;omitempty"`
	Address      *string `gorm:"not null;column:address;type:varchar(255);comment:地址;omitempty"`
	// Age          *int       `gorm:"not null;column:age;comment:年龄;omitempty"`
	GuardianId *int64     `gorm:"not null;column:guardian_id;comment:监护人Id;omitempty"`
	UpdatedAt  *time.Time `gorm:"not null;column:updated_at;autoUpdateTime;comment:更新时间;omitempty"`
	CreatedAt  *time.Time `gorm:"not null;column:created_at;autoCreateTime;comment:创建时间;omitempty"`
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
	if !public_db_func.CheckPhoneNumber(*m.PhoneNumber) {
		return fmt.Errorf("电话号码有误,请检查")
	}
	if !public_db_func.CheckIDCard(*m.IdcardNumber) {
		return fmt.Errorf("身份证号码有误,请重新拍照识别")
	}
	return db.Create(m).Error
}

// model中必须包含id,并且不为0
func (m *Model) UpdateById(db *gorm.DB) error {
	if m.Id == nil || *m.Id == 0 {
		return fmt.Errorf("")
	}
	return db.Where("id = ?", m.Id).Updates(m).Error
}

// 检查用户是否存在
func IsExcitOfPhoneNumber(db *gorm.DB, phone string) (bool, error) {
	var n int64
	err := db.Model(&Model{}).Where("phone_number = ?", phone).Count(&n).Error
	if err != nil || n <= 0 {
		return false, err
	}
	return true, nil
}

func IsExcit(db *gorm.DB) (bool, error) {
	var n int64
	err := db.Model(&Model{}).Count(&n).Error
	if err != nil || n <= 0 {
		return false, err
	}
	return true, nil
}

// 无记录:创建;有记录更新
func (m *Model) CreateOrUpdateByPhonenumber(db *gorm.DB) error {
	if m == nil || m.PhoneNumber == nil || len(*m.PhoneNumber) == 0 || !public_func.CheckPhoneNumber(*m.PhoneNumber) || m.IdcardNumber != nil || len(*m.IdcardNumber) > 0 || public_func.CheckIDCard(*m.IdcardNumber) {
		return fmt.Errorf("请求参数有误,请检查电话号码是否完整")
	}
	ok, err := IsExcitOfPhoneNumber(db, *m.PhoneNumber)
	if err != nil {
		return err
	}
	if ok {
		//存在:更新
		err = db.Model(m).Where("phone_number = ?", *m.PhoneNumber).Updates(m).Error
		// if err != nil {
		// 	return err
		// }
		// err := db.Model(&model_order.Model).Where("customer_id = ?", m.Id)
		return err
	} else {
		err = m.Create(db)
		return err
	}
}

func (m *Model) CreateOrUpdateByIdcardNumber(db *gorm.DB) error {
	if m == nil || m.PhoneNumber == nil || len(*m.PhoneNumber) == 0 || !public_func.CheckPhoneNumber(*m.PhoneNumber) || m.IdcardNumber != nil || len(*m.IdcardNumber) > 0 || public_func.CheckIDCard(*m.IdcardNumber) {
		return fmt.Errorf("请求参数有误,请检查电话号码是否完整")
	}

	var data Model
	err := db.Where("idcard_number = ?", *m.IdcardNumber).Select("id").First(&data).Error
	if err != nil {
		return err
	}

	if data.Id != nil && *data.Id != 0 {
		m.Id = data.Id
		//存在:更新
		err = db.Model(m).Where("idcard_number = ?", *m.IdcardNumber).Updates(m).Error
		// if err != nil {
		// 	return err
		// }
		// err := db.Model(&model_order.Model).Where("customer_id = ?", m.Id)
		return err
	} else {
		err = m.Create(db)
		return err
	}
}
