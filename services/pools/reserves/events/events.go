package reservesEvents

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/qrinef/arbitrage-bot/entities"
	"github.com/qrinef/arbitrage-bot/services/storage"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Service struct {
	storageService *storage.Service
	logger         *zap.SugaredLogger
	pools          *sync.Map
}

func NewService(storageService *storage.Service, logger *zap.SugaredLogger) *Service {
	return &Service{
		storageService: storageService,
		logger:         logger,
	}
}

func (s *Service) Start(pools *sync.Map) {
	s.pools = pools

	ch := make(chan types.Log, 1000)
	s.subscribe(s.storageService.Config.RpcWsEndpoint, ch)

	go s.handlerSubscribes(ch)
}

func (s *Service) subscribe(endpoint string, channel chan types.Log) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		s.logger.With(err).With("endpoint", endpoint).Fatal("PoolsService: connect to RPC failed")
	}

	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{{
			common.HexToHash("0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"),
		}},
	}

	sub, err := client.SubscribeFilterLogs(context.Background(), query, channel)
	if err != nil {
		s.logger.With(err).With("endpoint", endpoint).Fatal("PoolsService: subscribe to RPC failed")
	}

	s.logger.With("endpoint", endpoint).Info("PoolsService: subscribe success")
	go s.handlerErrors(endpoint, sub)
}

func (s *Service) handlerSubscribes(ch chan types.Log) {
	for {
		select {
		case it := <-ch:
			if values, ok := s.pools.Load(it.Address); ok {
				pool := values.(entities.Reserves)

				if it.BlockNumber > pool.LastBlock || (it.BlockNumber == pool.LastBlock && it.Index > pool.LastIndex) {
					s.pools.Store(it.Address, entities.Reserves{
						Reserve0:  new(uint256.Int).SetBytes(it.Data[:32]),
						Reserve1:  new(uint256.Int).SetBytes(it.Data[32:]),
						LastBlock: it.BlockNumber,
						LastIndex: it.Index,
					})
				}
			}
		}
	}
}

func (s *Service) handlerErrors(endpoint string, sub ethereum.Subscription) {
	for {
		select {
		case err := <-sub.Err():
			s.logger.With(err).With("endpoint", endpoint).Info("PoolsService: subscribe reported an error")
		}
	}
}
