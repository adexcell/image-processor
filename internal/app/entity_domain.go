package app

import (
	"github.com/adexcell/image-processor/internal/entity/adapter/postgres"
	"github.com/adexcell/image-processor/internal/entity/adapter/rabbit"
	"github.com/adexcell/image-processor/internal/entity/adapter/redis"
	httprouter "github.com/adexcell/image-processor/internal/entity/controller/http_router"
	"github.com/adexcell/image-processor/internal/entity/usecase"
)

func EntityDomain(d Dependencies) {
	entityUseCase := usecase.New(
		postgres.New(d.Postgres),
		redis.New(d.Redis),
		rabbit.New(d.RabbitMQ),
	)

	httprouter.EntityRouter(d.RouterHTTP, entityUseCase, d.Metrics)
}
