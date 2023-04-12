package container

import (
	"github.com/qrinef/arbitrage-bot/services/compute"
	"github.com/qrinef/arbitrage-bot/services/discovery"
	"github.com/qrinef/arbitrage-bot/services/pools"
	poolsCollector "github.com/qrinef/arbitrage-bot/services/pools/collector"
	reservesCollector "github.com/qrinef/arbitrage-bot/services/pools/reserves/collector"
	reservesEvents "github.com/qrinef/arbitrage-bot/services/pools/reserves/events"
	"github.com/qrinef/arbitrage-bot/services/storage"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Service struct {
	LoggerService    *zap.SugaredLogger
	StorageService   *storage.Service
	DiscoveryService *discovery.Service
	PoolsService     *pools.Service
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start() {
	container := dig.New()

	container.Provide(func() *zap.SugaredLogger {
		logger, _ := zap.NewProduction()
		return logger.Sugar()
	})

	container.Provide(storage.NewService)
	container.Provide(discovery.NewService)
	container.Provide(pools.NewService)
	container.Provide(poolsCollector.NewService)
	container.Provide(reservesCollector.NewService)
	container.Provide(reservesEvents.NewService)
	container.Provide(compute.NewService)

	err := container.Invoke(func(
		loggerService *zap.SugaredLogger,
		storageService *storage.Service,
		discoveryService *discovery.Service,
		poolsService *pools.Service,
	) {
		s.LoggerService = loggerService
		s.StorageService = storageService
		s.DiscoveryService = discoveryService
		s.PoolsService = poolsService

		s.StorageService.Start()
		s.LoggerService.Info("Container fully prepared")
	})
	if err != nil {
		panic(err)
	}
}
