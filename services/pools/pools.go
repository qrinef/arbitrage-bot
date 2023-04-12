package pools

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/qrinef/arbitrage-bot/entities"
	poolsCollector "github.com/qrinef/arbitrage-bot/services/pools/collector"
	reservesCollector "github.com/qrinef/arbitrage-bot/services/pools/reserves/collector"
	reservesEvents "github.com/qrinef/arbitrage-bot/services/pools/reserves/events"
	"github.com/qrinef/arbitrage-bot/services/storage"
	"sync"
)

type Service struct {
	storageService    *storage.Service
	poolsCollector    *poolsCollector.Service
	reservesEvents    *reservesEvents.Service
	reservesCollector *reservesCollector.Service

	pools sync.Map
}

func NewService(storageService *storage.Service, poolsCollector *poolsCollector.Service, reservesEvents *reservesEvents.Service, reservesCollector *reservesCollector.Service) *Service {
	return &Service{
		storageService:    storageService,
		poolsCollector:    poolsCollector,
		reservesEvents:    reservesEvents,
		reservesCollector: reservesCollector,
	}
}

func (s *Service) Start() {
	s.poolsCollector.Start()

	s.setPools()

	s.reservesEvents.Start(&s.pools)
	s.reservesCollector.Start(&s.pools)
}

func (s *Service) setPools() {
	var poolsFromDb []entities.Pool
	s.storageService.Db.
		Offset(-1).Limit(-1).
		Find(&poolsFromDb)

	for _, pool := range poolsFromDb {
		s.pools.Store(pool.AddressCommon, entities.Reserves{})
	}
}

func (s *Service) GetReserves(factory, tokenA, tokenB common.Address) (reserveA, reserveB, fee *uint256.Int, pool common.Address, err error) {
	return reserveA, reserveB, fee, pool, err
}
