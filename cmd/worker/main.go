package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adexcell/image-processor/internal/config"
	"github.com/adexcell/image-processor/internal/models"
	"github.com/adexcell/image-processor/internal/processor"
	"github.com/adexcell/image-processor/internal/storage"
	"github.com/segmentio/kafka-go"
	kafkav2 "github.com/wb-go/wbf/kafka/kafka-v2"
	"github.com/wb-go/wbf/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, _ := logger.InitLogger(logger.SlogEngine, cfg.App.Name, cfg.App.Env)
	log.Info("Worker starting...")

	s, err := storage.NewLocalStorage(cfg.Storage.UploadsDir, cfg.Storage.ProcessedDir)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	p := processor.NewImageProcessor()

	consumer := kafkav2.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Group, log)
	proc, err := kafkav2.NewProcessor(consumer, nil, log, kafkav2.MaxAttempts(3))
	if err != nil {
		log.Error("failed to create processor", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	proc.Start(ctx, func(ctx context.Context, msg kafka.Message) error {
		var task models.TaskMessage
		if err := json.Unmarshal(msg.Value, &task); err != nil {
			return fmt.Errorf("failed to unmarshal task: %w", err)
		}

		log.Info("Processing task", "id", task.ID, "file", task.Filename)

		status, _ := s.GetStatus(task.ID)
		status.Status = models.StatusProcessing
		_ = s.SetStatus(task.ID, status)

		srcPath := s.GetPath(task.ID, task.Filename, false)
		srcFile, err := os.Open(srcPath)
		if err != nil {
			status.Status = models.StatusFailed
			status.Error = err.Error()
			_ = s.SetStatus(task.ID, status)
			return err
		}
		defer srcFile.Close()

		out, err := p.Process(srcFile, task.ProcessOps)
		if err != nil {
			status.Status = models.StatusFailed
			status.Error = err.Error()
			_ = s.SetStatus(task.ID, status)
			return err
		}

		_, err = s.Save(task.ID, task.Filename, out, true)
		if err != nil {
			status.Status = models.StatusFailed
			status.Error = err.Error()
			_ = s.SetStatus(task.ID, status)
			return err
		}

		status.Status = models.StatusCompleted
		_ = s.SetStatus(task.ID, status)

		log.Info("Task completed", "id", task.ID)
		return nil
	})

	log.Info("Worker is running. Press Ctrl+C to stop.")
	<-ctx.Done()
	log.Info("Worker shutting down...")
}
