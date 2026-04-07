package storage

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/adexcell/image-processor/internal/models"
)

type Storage interface {
	Save(id string, filename string, r io.Reader, processed bool) (string, error)
	GetPath(id string, filename string, processed bool) string
	Delete(id string) error
	Exists(id string, filename string, processed bool) bool
	SetStatus(id string, status models.ImageInfo) error
	GetStatus(id string) (models.ImageInfo, error)
}

type LocalStorage struct {
	uploadsDir   string
	processedDir string
}

func NewLocalStorage(uploadsDir, processedDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(processedDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{
		uploadsDir:   uploadsDir,
		processedDir: processedDir,
	}, nil
}

func (s *LocalStorage) Save(id string, filename string, r io.Reader, processed bool) (string, error) {
	dir := s.uploadsDir
	if processed {
		dir = s.processedDir
	}

	itemDir := filepath.Join(dir, id)
	if err := os.MkdirAll(itemDir, 0755); err != nil {
		return "", err
	}

	dstPath := filepath.Join(itemDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, r); err != nil {
		return "", err
	}

	return dstPath, nil
}

func (s *LocalStorage) GetPath(id string, filename string, processed bool) string {
	dir := s.uploadsDir
	if processed {
		dir = s.processedDir
	}
	return filepath.Join(dir, id, filename)
}

func (s *LocalStorage) Delete(id string) error {
	_ = os.RemoveAll(filepath.Join(s.uploadsDir, id))
	_ = os.RemoveAll(filepath.Join(s.processedDir, id))
	return nil
}

func (s *LocalStorage) Exists(id string, filename string, processed bool) bool {
	path := s.GetPath(id, filename, processed)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (s *LocalStorage) SetStatus(id string, info models.ImageInfo) error {
	path := filepath.Join(s.uploadsDir, id, "status.json")
	_ = os.MkdirAll(filepath.Dir(path), 0755)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(info)
}

func (s *LocalStorage) GetStatus(id string) (models.ImageInfo, error) {
	path := filepath.Join(s.uploadsDir, id, "status.json")
	f, err := os.Open(path)
	if err != nil {
		return models.ImageInfo{}, err
	}
	defer f.Close()

	var info models.ImageInfo
	err = json.NewDecoder(f).Decode(&info)
	return info, err
}
