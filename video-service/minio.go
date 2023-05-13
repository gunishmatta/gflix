package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"time"
)

func UploadToMinio(file multipart.File, header *multipart.FileHeader) (string, error) {
	endpoint := "localhost:9000"
	accessKeyID := "AWWDpPvVFabkJDWJ"
	secretAccessKey := "6sLvaULz338lrEEEJNzCYFkO8t4a84MT"
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	bucketName := "videos"
	fileName := fmt.Sprintf("%d%s", time.Now().Unix(), filepath.Ext(header.Filename))

	_, err = minioClient.PutObject(context.Background(), bucketName, fileName, file, header.Size, minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")})
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	u := url.URL{
		Scheme: "http",
		Host:   "localhost:9000",
		Path:   fmt.Sprintf("/%s/%s", fileName, bucketName),
	}
	return u.String(), nil
}
