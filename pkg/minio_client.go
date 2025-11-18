package pkg

import (
	"context"
	"fmt"
	"io"

	"develop-internet-applications/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
	config *config.Config
}

func NewMinioClient(cfg *config.Config) (*MinioClient, error) {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio124", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &MinioClient{
		client: client,
		config: cfg,
	}, nil
}

func (m *MinioClient) UploadImage(starID int, file io.Reader, fileSize int64, filename string) (string, error) {
	objectName := fmt.Sprintf("%d.png", starID)

	_, err := m.client.PutObject(context.Background(), "stars", objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: "image/png",
	})
	if err != nil {
		return "", err
	}

	// Формируем URL используя конфигурацию
	protocol := m.config.ServiceProtocol
	if protocol == "" {
		protocol = "https"
	}
	publicHost := m.config.ServicePublicHost
	if publicHost == "" {
		publicHost = "172.20.10.4"
	}
	port := m.config.ServicePort
	if port == 0 {
		port = 8080
	}

	return fmt.Sprintf("%s://%s:%d/api/stars/%d/image", protocol, publicHost, port, starID), nil
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
