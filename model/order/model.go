package model_order

import (
	"time"
	public_db_func "youke/model/public"

	"gorm.io/gorm"
)

type Model struct {
	Id           int64     `gorm:"primaryKey;autoIncrement;column:id;comment:主键,自增"`
	CustomerId   int64     `gorm:"not null;column:customer_id;comment:顾客ID"`
	CustomerName int64     `gorm:"not null;column:customer_name;type:varchar(10);comment:顾客名字"`
	PhoneNumber  string    `gorm:"not null;column:phone_number;type:varchar(15);comment:电话号码"`
	RoomNumber   string    `gorm:"not null;column:room_number;type:varchar(10);comment:房间号"`
	Price        int       `gorm:"not null;column:price;comment:价格"`
	Ymd          time.Time `gorm:"index:ymd_index;not null;column:ymd;comment:创建日期(不含时分秒)"`
	CreatedAt    time.Time `gorm:"not null;column:created_at;autoCreateTime;comment:创建时间"`
}

// 自定义表名
func (Model) TableName() string {
	return "order" // 自定义表名
}

// 自动建表
func CreateTable(db *gorm.DB) error {
	// 自动迁移，创建表并加注释
	err := db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='订单登记记录表'").AutoMigrate(&Model{})
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) Full(tx *gorm.DB) {
	m.Ymd = time.Now().Truncate(24 * time.Hour) // 只保留年月日，去除时分秒
}

// 按日期搜索订单
func SelectByYmd(db *gorm.DB, t string) ([]Model, error) {
	ymd, err := time.Parse(time.DateOnly, t)
	if err != nil {
		return nil, err
	}
	var result []Model
	condition := public_db_func.OrderIsDesc(db, "created_at", true)
	err = condition.Where("ymd = ?", ymd).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, err
}

func (m *Model) Create(db *gorm.DB) error {
	return db.Create(m).Error
}
