package main

import (
	"context"

	"github.com/adexcell/image-processor/config"
	"github.com/adexcell/image-processor/internal/app"
	"github.com/adexcell/image-processor/pkg/logger"
	"github.com/adexcell/image-processor/pkg/otel"
	"github.com/rs/zerolog/log"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logger.Init(c.Logger)

	ctx := context.Background()

	if err = otel.Init(ctx, c.OTEL); err != nil {
		log.Error().Err(err).Msg("otel.Init")
	}
	defer otel.Close()

	if err := app.Run(ctx, c); err != nil {
		log.Error().Err(err).Msg("app.Run")
	}
}
