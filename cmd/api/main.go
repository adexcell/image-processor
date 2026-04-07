package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/adexcell/image-processor/internal/api"
	"github.com/adexcell/image-processor/internal/config"
	"github.com/adexcell/image-processor/internal/storage"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, _ := logger.InitLogger(logger.SlogEngine, cfg.App.Name, cfg.App.Env)

	s, err := storage.NewLocalStorage(cfg.Storage.UploadsDir, cfg.Storage.ProcessedDir)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer producer.Close()

	h := api.NewHandler(cfg, s, producer, log)

	r := ginext.New("debug")
	r.Use(ginext.Logger(), ginext.Recovery())

	// Static files for UI
	r.StaticFile("/", "./ui/index.html")
	r.Static("/ui", "./ui")

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/upload", h.Upload)
		apiGroup.GET("/image/:id", h.GetImage)
		apiGroup.GET("/status/:id", h.GetStatus)
		apiGroup.DELETE("/image/:id", h.Delete)
		apiGroup.GET("/list", h.List)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown:", "error", err)
	}
}
