package public_db_func

import (
	"regexp"
	"strconv"
	"time"
	"youke/global"

	jsoniter "github.com/json-iterator/go"

	"gorm.io/gorm"
)

func IsTime(t string) bool {
	_, err := time.Parse(time.DateTime, t)
	if err == nil {
		return true
	}
	_, err = time.Parse("2006-01-02 15:04", t)
	if err == nil {
		return true
	}
	_, err = time.Parse("2006-01-02 15", t)
	if err == nil {
		return true
	}
	_, err = time.Parse(time.DateOnly, t)
	return err == nil
}

func RangeOfTime(db *gorm.DB, clum string, Range []string) *gorm.DB {
	switch len(Range) {
	case 0:
		break
	case 1:
		if IsTime(Range[0]) {
			db.Where(clum+" >= ?", Range[0])
		}
	case 2:
		if IsTime(Range[0]) {
			db.Where(clum+" >= ?", Range[0])
		}
		if IsTime(Range[1]) {
			db.Where(clum+" <= ?", Range[1])
		}

	default:
		global.Global.Logger.Warning("错误的时间范围长度,时间范围长度必须小于等于2")
	}
	return db
}

// 分页查询
// 如果你不想跳过任何数据,请从第0页开始查询
func PageCondition(db *gorm.DB, page, pageSize int) *gorm.DB {
	offset := page * pageSize
	return db.Offset(offset).Limit(pageSize)
}

// true 倒叙; false 正序
func OrderIsDesc(db *gorm.DB, clum string, desc bool) *gorm.DB {
	if desc {
		return db.Order(clum + " DESC")
	}
	return db.Order(clum + " ASC")
}

func CheckPhoneNumber(phone string) bool {
	// 定义中国大陆手机号码的正则表达式
	mobilePattern := `^1[3-9]\d{9}$`
	// 定义中国大陆固定电话号码的正则表达式 (区号+电话号码)
	// landlinePattern := `^(\d{3,4}-)?\d{7,8}$`

	// 编译正则表达式
	mobileRegex := regexp.MustCompile(mobilePattern)
	// landlineRegex := regexp.MustCompile(landlinePattern)

	// 检查是否匹配手机
	return mobileRegex.MatchString(phone)
}

// 校验是否为有效的中国大陆身份证号
func CheckIDCard(id string) bool {
	// 15位身份证号的正则表达式
	// pattern15 := `^\d{15}$`
	// 18位身份证号的正则表达式
	pattern18 := `^\d{17}[\dXx]$`

	// 编译正则表达式
	// regex15 := regexp.MustCompile(pattern15)
	regex18 := regexp.MustCompile(pattern18)

	// 判断是否符合15位或18位身份证号格式
	// if regex15.MatchString(id) {
	// 	return true
	// }
	// if regex18.MatchString(id) {
	// 	return check18IDCardChecksum(id)
	// }

	return regex18.MatchString(id)
}

// 校验18位身份证的校验码
func check18IDCardChecksum(id string) bool {
	// 身份证号前17位的加权因子
	weightFactors := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 18位身份证号的校验码对应表
	checksumTable := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	// 计算前17位的加权和
	sum := 0
	for i := 0; i < 17; i++ {
		num, err := strconv.Atoi(string(id[i]))
		if err != nil {
			return false
		}
		sum += num * weightFactors[i]
	}

	// 取模11，得到校验码的索引
	checksumIndex := sum % 11

	// 比较计算得到的校验码和身份证的第18位
	return id[17] == checksumTable[checksumIndex] || (id[17] == 'x' && checksumTable[checksumIndex] == 'X')
}

// 结构体快速映射
// to : 传递指针变量
func StructToStruct(val, to interface{}) error {
	data, err := jsoniter.Marshal(val)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(data, to)
	if err != nil {
		return err
	}

	return nil
}
