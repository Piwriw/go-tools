package main

import (
	"log"

	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/minio/minio-go/v7"
)

var (
	minioClient *minio.Client
)

func initMinio() {
	// MinIO配置
	endpoint := "113.45.181.99:9000"  // MinIO的地址
	accessKeyID := "Joohwan"          // 你的 Access Key
	secretAccessKey := "Joohwan2020." // 你的 Secret Key
	useSSL := false                   // 是否使用SSL

	// 初始化 MinIO 客户端
	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

}
