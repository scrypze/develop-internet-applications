package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "develop-internet-applications/docs"
	"develop-internet-applications/internal/handler"
	"develop-internet-applications/internal/repository"
	"develop-internet-applications/internal/service"
	"develop-internet-applications/pkg"
	"develop-internet-applications/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config      *config.Config
	Router      *gin.Engine
	Handler     *handler.Handler
	RedisClient *pkg.RedisClient
}

func NewApp() *Application {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresCofigString := config.ConfDB()
	db, err := pkg.NewPostgresDB(postgresCofigString)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	config.MigrateDB()
	// if err := pkg.SyncSequences(db); err != nil {
	// 	logrus.Fatalf("failed to sync sequences: %v", err)
	// }

	minioClient, err := pkg.NewMinioClient(conf)
	if err != nil {
		logrus.Fatalf("failed to create minio client: %v", err)
	}

	ctx := context.Background()
	redisClient, err := pkg.NewRedisClient(ctx, conf.Redis)
	if err != nil {
		logrus.Fatalf("failed to create redis client: %v", err)
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository, minioClient, redisClient, conf)
	handler := handler.NewHandler(service)
	router := gin.Default()

	return &Application{
		Config:      conf,
		Router:      router,
		Handler:     handler,
		RedisClient: redisClient,
	}
}

func (a *Application) RunApp() {

	a.Handler.RegisterHandler(a.Router)
	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	server := &http.Server{
		Addr:    serverAddress,
		Handler: a.Router,
	}

	// Настройка HTTPS если протокол https и указаны сертификаты
	useHTTPS := false
	if a.Config.ServiceProtocol == "https" && a.Config.TLSCertFile != "" && a.Config.TLSKeyFile != "" {
		// Проверяем существование файлов сертификатов
		if _, err := os.Stat(a.Config.TLSCertFile); os.IsNotExist(err) {
			logrus.Warnf("TLS certificate file not found: %s, falling back to HTTP", a.Config.TLSCertFile)
		} else if _, err := os.Stat(a.Config.TLSKeyFile); os.IsNotExist(err) {
			logrus.Warnf("TLS key file not found: %s, falling back to HTTP", a.Config.TLSKeyFile)
		} else {
			// Загружаем сертификаты
			cert, err := tls.LoadX509KeyPair(a.Config.TLSCertFile, a.Config.TLSKeyFile)
			if err != nil {
				logrus.Fatalf("Failed to load TLS certificates: %v", err)
			}

			server.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
			useHTTPS = true
		}
	}

	// Запускаем сервер
	go func() {
		if useHTTPS {
			logrus.Printf("Service started with HTTPS on %s", serverAddress)
			if err := server.ListenAndServeTLS(a.Config.TLSCertFile, a.Config.TLSKeyFile); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("Server error: %v", err)
			}
		} else {
			logrus.Printf("Service started with HTTP on %s", serverAddress)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("Server error: %v", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown error: %v", err)
	}

	if err := a.RedisClient.Close(); err != nil {
		logrus.Errorf("Redis close error: %v", err)
	}

	logrus.Info("Server down")
}
