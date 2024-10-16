package operation

import (
	"context"

	"github.com/minio/minio-go/v7"
)

// CheckBucket 检查是否存在存储桶
func CheckBucket(client *minio.Client, ctx context.Context, bucketName string) (err error) {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
