package processor

import (
	"bytes"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/disintegration/imaging"
)

type Processor interface {
	Process(srcReader io.Reader, ops []string) (io.Reader, error)
}

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) Process(srcReader io.Reader, ops []string) (io.Reader, error) {
	src, err := imaging.Decode(srcReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	for _, op := range ops {
		switch op {
		case "resize":
			src = imaging.Resize(src, 800, 0, imaging.Lanczos)
		case "thumbnail":
			src = imaging.Thumbnail(src, 200, 200, imaging.Lanczos)
		}
	}

	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, src, imaging.JPEG)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf, nil
}
