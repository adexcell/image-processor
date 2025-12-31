package main

import (
	"github.com/adexcell/image-processor/internal/image"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {

	zlog.Init()

	zlog.Logger.Info().Msg("httprouter create")
	httprouter := ginext.New("debug")

	zlog.Logger.Info().Msg("add image routes")
	imageHandler := image.New()
	imageHandler.Register(httprouter)
}
