package models

import "time"

type TaskMessage struct {
	ID         string    `json:"id"`
	Filename   string    `json:"filename"`
	CreatedAt  time.Time `json:"created_at"`
	ProcessOps []string  `json:"process_ops"` // e.g. ["resize", "thumbnail", "watermark"]
}

type ImageStatus string

const (
	StatusPending    ImageStatus = "pending"
	StatusProcessing ImageStatus = "processing"
	StatusCompleted  ImageStatus = "completed"
	StatusFailed     ImageStatus = "failed"
)

type ImageInfo struct {
	ID        string      `json:"id"`
	Filename  string      `json:"filename"`
	Status    ImageStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	Error     string      `json:"error,omitempty"`
}
