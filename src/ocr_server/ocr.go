package ocr_server

import (
	"encoding/base64"
	"fmt"
	"time"
	"youke/global"
	public_db_func "youke/model/public"

	jsoniter "github.com/json-iterator/go"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

const Config = `{"CropPortrait":true}`

func getConfig() *string {
	c := Config
	return &c
}

type OcrIdCardResp struct {
	HeadImg []byte
	Name    string `json:"Name"`
	Age     int
	Address string `json:"Address"`
	IdNum   string `json:"IdNum"`
}

func IdCardOCR(img string) (*OcrIdCardResp, error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		global.Global.Config.Ocr.SecretID,
		global.Global.Config.Ocr.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ocr.NewClient(credential, "ap-guangzhou", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ocr.NewIDCardOCRRequest()
	request.ImageBase64 = &img
	request.Config = getConfig()

	// 返回的resp是一个IDCardOCRResponse的实例，与请求对象对应
	response, err := client.IDCardOCR(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned;err:%v", err)
	}
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("返回数据为空")
	}

	var headImg []byte
	{
		headBase64 := ""

		if response.Response.AdvancedInfo != nil || len(*response.Response.AdvancedInfo) > 10 {
			headBase64, err = ExtractPortrait(*response.Response.AdvancedInfo)
			if err != nil {
				global.Global.Logger.Warning(err)
			} else if len(headBase64) > 0 {
				headImg, err = ConvertBase64Tobytes(headBase64)
				if err != nil {
					global.Global.Logger.Warning(err)
				}
			}
		}

	}

	rsp := new(OcrIdCardResp)
	err = public_db_func.StructToStruct(response.Response, rsp)
	if err != nil {
		return nil, err
	}
	rsp.HeadImg = headImg
	if response.Response.Birth != nil && len(*response.Response.Birth) > 0 {
		rsp.Age, err = CalculateAge(*response.Response.Birth)
		if err != nil {
			global.Global.Logger.Warning(err)
		}
	}

	return rsp, nil

}

// ConvertBase64ToReader 将 Base64 编码的字符串解码为二进制数据
func ConvertBase64Tobytes(base64Str string) ([]byte, error) {
	// 解码 Base64 字符串为二进制数据
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type PortraitResponse struct {
	Portrait string `json:"Portrait"`
}

// ExtractPortrait 简化版方法，提取 "Portrait" 字段
func ExtractPortrait(jsonStr string) (string, error) {
	var response PortraitResponse

	// 仅解码 "Portrait" 字段
	err := jsoniter.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return "", err
	}

	return response.Portrait, nil
}

// CalculateAge 根据生日字符串计算年龄
// 时间格式: 2002/6/16
func CalculateAge(birthdayStr string) (int, error) {
	// 定义生日的日期格式
	layout := "2006/1/2"

	birthday, err := time.Parse(layout, birthdayStr)
	if err != nil {
		return 0, err
	}
	now := time.Now()

	age := now.Year() - birthday.Year()

	if now.YearDay() < birthday.YearDay() {
		age--
	}

	return age, nil
}
