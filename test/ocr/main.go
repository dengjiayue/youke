package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"youke/src/ocr_server"
)

func ConvertImageToBase64(imagePath string) (string, error) {
	// 打开文件
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 读取文件内容
	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// 将文件数据编码为Base64
	base64String := base64.StdEncoding.EncodeToString(imageData)
	return base64String, nil
}

func main() {
	b, err := ConvertImageToBase64("1.jpg")
	if err != nil {
		fmt.Printf("err1=%#v\n", err)
		return
	}
	fmt.Printf("%s\n", b)
	rsp, err := ocr_server.IdCardOCR(b)
	if err != nil {
		fmt.Printf("err2=%#v\n", err)
		return
	}
	fmt.Printf("rsp=%#v\n", rsp)
}
