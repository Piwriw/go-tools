package operation

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
)

const (
	defaultShareTime = 5 * time.Minute
)

// ShareLink 生成一个可分享的下载链接
// 如果不传expires参数，将使用默认过期时间。
func ShareLink(minioClient *minio.Client, ctx context.Context, bucketName, objectName string, expires *time.Duration) (*url.URL, error) {
	// 设置默认过期时间为 5 分钟
	if expires == nil {
		sharTime := defaultShareTime
		expires = &sharTime // 使用默认值
	}

	// 生成下载链接
	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, *expires, url.Values{})
	if err != nil {
		return nil, err
	}
	return presignedURL, nil
}

// UploadObject - 上传对象到 MinIO
func UploadObject(minioClient *minio.Client, bucketName, objectName string, fileData []byte) error {

	fileSize := int64(len(fileData))
	if fileSize < 0 {
		return errors.New("file size is zero")
	}

	// 上传字节数组
	_, err := minioClient.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(fileData), fileSize, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DownloadFileFromURL 使用预签名的 URL 下载文件
func DownloadFileFromURL(url, filePath string) error {
	// 发送 GET 请求
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将响应内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
