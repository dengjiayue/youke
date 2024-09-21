package cos

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type CosConfig struct {
	CosUrl    string `yaml:"CosUrl"`
	SecretID  string `yaml:"SecretID"`
	SecretKey string `yaml:"SecretKey"`
}

// 建立cos连接
// 连接腾讯云cos
func NewCosClient(config CosConfig) (*cos.Client, error) {
	u, _ := url.Parse(config.CosUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})

	// 验证连接
	_, _, err := client.Bucket.Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// 定义文件上传功能的接口
func UploadFile(client *cos.Client, file io.Reader, fileSize int64, objectKey string) (string, error) {
	const chunkSize = 5 * 1024 * 1024 // 分块大小：5MB

	// 根据文件大小判断是否使用分块上传
	if fileSize <= chunkSize {
		// 如果文件较小，使用简单上传
		return simpleUpload(client, file, objectKey)
	} else {
		// 如果文件较大，使用分块上传
		return multipartUpload(client, file, objectKey)
	}
}

// 简单上传方法
func simpleUpload(client *cos.Client, file io.Reader, objectKey string) (string, error) {
	_, err := client.Object.Put(context.Background(), objectKey, file, nil)
	if err != nil {
		return "", err
	}

	// 返回下载 URL
	url := client.Object.GetObjectURL(objectKey).String()
	return url, nil
}

// 分块上传方法
func multipartUpload(client *cos.Client, file io.Reader, objectKey string) (string, error) {
	const chunkSize = 3 * 1024 * 1024 // 每块 3MB
	var (
		partNumber     = 1
		completedParts []cos.Object
	)

	// 初始化分块上传
	initResp, _, err := client.Object.InitiateMultipartUpload(context.Background(), objectKey, nil)
	if err != nil {
		return "", err
	}
	uploadID := initResp.UploadID

	buffer := make([]byte, chunkSize) // 创建固定大小的缓冲区用于分块上传
	for {
		n, err := file.Read(buffer) // 逐块读取文件
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break // 读取完毕，退出循环
		}

		// 上传当前分块
		partResp, err := client.Object.UploadPart(
			context.Background(),
			objectKey,
			uploadID,
			partNumber,
			io.NopCloser(io.MultiReader(bytes.NewReader(buffer[:n]))),
			nil,
		)
		if err != nil {
			return "", err
		}

		// 记录已上传分块的信息
		completedParts = append(completedParts, cos.Object{
			PartNumber: partNumber,
			ETag:       partResp.Header.Get("ETag"),
		})
		partNumber++
	}

	// 完成分块上传
	completeOpt := &cos.CompleteMultipartUploadOptions{Parts: completedParts}
	_, _, err = client.Object.CompleteMultipartUpload(context.Background(), objectKey, uploadID, completeOpt)
	if err != nil {
		return "", err
	}

	// 返回下载 URL
	url := client.Object.GetObjectURL(objectKey).String()
	return url, nil
}
