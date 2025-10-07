package pkg

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient() (*MinioClient, error) {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio124", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &MinioClient{client: client}, nil
}

func (m *MinioClient) UploadImage(star string, file io.Reader, fileSize int64, filename string) (string, error) {
	objectName := fmt.Sprintf("%s%s", star, ".png")

	_, err := m.client.PutObject(context.Background(), "stars", objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: "image/png",
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:9000/stars/%s", objectName), nil
}

func (m *MinioClient) DeleteImage(starID int) error {
	objectName := fmt.Sprintf("%d.png", starID)

	err := m.client.RemoveObject(context.Background(), "stars", objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) GetImageName(starID int) string {
	return fmt.Sprintf("%d.png", starID)
}
