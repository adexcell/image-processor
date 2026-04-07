package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/adexcell/image-processor/internal/config"
	"github.com/adexcell/image-processor/internal/models"
	"github.com/adexcell/image-processor/internal/storage"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/logger"
)

type Handler struct {
	cfg      *config.Config
	storage  storage.Storage
	producer *kafka.Producer
	log      logger.Logger
}

func NewHandler(cfg *config.Config, s storage.Storage, p *kafka.Producer, log logger.Logger) *Handler {
	return &Handler{
		cfg:      cfg,
		storage:  s,
		producer: p,
		log:      log,
	}
}

func (h *Handler) Upload(c *ginext.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "no file uploaded"})
		return
	}

	id := uuid.New().String()
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to open file"})
		return
	}
	defer f.Close()

	_, err = h.storage.Save(id, file.Filename, f, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to save file"})
		return
	}

	status := models.ImageInfo{
		ID:        id,
		Filename:  file.Filename,
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}
	_ = h.storage.SetStatus(id, status)

	// Send to Kafka
	msg := models.TaskMessage{
		ID:         id,
		Filename:   file.Filename,
		CreatedAt:  time.Now(),
		ProcessOps: []string{"resize", "thumbnail"},
	}
	msgBytes, _ := json.Marshal(msg)

	err = h.producer.Send(c.Request.Context(), []byte(id), msgBytes)
	if err != nil {
		h.log.Error("failed to send message to kafka", "error", err)
		// We could delete the file here, but it's better to keep it and maybe retry
	}

	c.JSON(http.StatusAccepted, status)
}

func (h *Handler) GetStatus(c *ginext.Context) {
	id := c.Param("id")
	status, err := h.storage.GetStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "image not found"})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *Handler) GetImage(c *ginext.Context) {
	id := c.Param("id")
	status, err := h.storage.GetStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "image not found"})
		return
	}

	path := h.storage.GetPath(id, status.Filename, true)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return original if processed not ready? No, task says "get processed image"
		c.JSON(http.StatusProcessing, ginext.H{"status": "processing"})
		return
	}

	c.File(path)
}

func (h *Handler) Delete(c *ginext.Context) {
	id := c.Param("id")
	err := h.storage.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to delete"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) List(c *ginext.Context) {
	// Simple implementation: list directories in uploads
	files, err := os.ReadDir(h.cfg.Storage.UploadsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to list"})
		return
	}

	var result []models.ImageInfo
	for _, f := range files {
		if f.IsDir() {
			status, err := h.storage.GetStatus(f.Name())
			if err == nil {
				result = append(result, status)
			}
		}
	}
	c.JSON(http.StatusOK, result)
}
