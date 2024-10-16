package operation

import (
	"context"
	"os"
	"testing"
)

func TestUploadObject(t *testing.T) {
	initMinio()
	filePath := "../minio.png"
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if err = UploadObject(minioClient, "go-minio-bucket", "minio-1.png", data); err != nil {
		t.Fatalf("failed to upload file: %v", err)
	}
}

func TestShareLink(t *testing.T) {
	initMinio()
	// 定义存储桶(bucket)名称
	bucketName := "go-minio-bucket"
	objectName := "minio-1.png"
	url, err := ShareLink(minioClient, context.TODO(), bucketName, objectName, nil)
	if err != nil {
		t.Fatalf("failed to upload file: %v", err)
	}
	t.Logf("object url:%v", url)
}

func TestDownloadFile(t *testing.T) {
	initMinio()
	bucketName := "go-minio-bucket"

	// 使用生成的链接下载文件
	downloadFilePath := "downloaded-minio.png"
	objectName := "minio-1.png"
	url, err := ShareLink(minioClient, context.TODO(), bucketName, objectName, nil)
	if err != nil {
		t.Fatalf("failed to upload file: %v", err)
	}

	if err = DownloadFileFromURL(url.String(), downloadFilePath); err != nil {
		t.Fatalf("failed to download file: %v", err)
	}
	t.Log("File downloaded successfully to:", downloadFilePath)
}
