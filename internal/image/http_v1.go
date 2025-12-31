package image

import (
	"github.com/adexcell/image-processor/internal/controller"
	"github.com/wb-go/wbf/ginext"
)

type handler struct{

}

func New() controller.Handler {
	return &handler{}
}

func (h *handler) Register(router *ginext.Engine) {
	router.POST("/upload", h.UploadImage)
	router.GET("/image/:id", h.GetImage)
	router.DELETE("/image/:id", h.DeleteImage)
}

func (h *handler) UploadImage(c *ginext.Context) {

}

func (h *handler) GetImage(c *ginext.Context) {

}

func (h *handler) DeleteImage(c *ginext.Context) {

}
