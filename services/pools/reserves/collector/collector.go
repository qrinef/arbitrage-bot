package reservesCollector

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/qrinef/arbitrage-bot/entities"
	"github.com/qrinef/arbitrage-bot/services/storage"
	"github.com/qrinef/arbitrage-bot/utils/multicall"
	"go.uber.org/zap"
	"sync"
)

var (
	ReservesSignature = common.FromHex("0902f1ac")
	MultiCallAddress  = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")
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

	clientETH, err := ethclient.Dial(s.storageService.Config.RpcWsEndpoint)
	if err != nil {
		s.logger.With(err).Fatal("EthClient initial failed")
	}

	multiCall, err := multicall.NewMulticall(MultiCallAddress, clientETH)
	if err != nil {
		s.logger.With(err).Fatal("Liquidity multiCall initial")
	}

	var tokensPools []entities.Pool
	s.storageService.Db.
		Offset(-1).Limit(-1).
		Find(&tokensPools)

	batches := s.getBatches(tokensPools)

	var wg sync.WaitGroup

	for _, batch := range batches {
		wg.Add(1)

		go func(wg *sync.WaitGroup, batch []entities.Pool) {
			defer wg.Done()

			var requests []multicall.Multicall3Call
			for _, pool := range batch {
				requests = append(requests, multicall.Multicall3Call{
					Target:   common.HexToAddress(pool.Address),
					CallData: ReservesSignature,
				})
			}

			callAggregate, err := multiCall.TryAggregate(&bind.CallOpts{}, false, requests)
			if err != nil {
				return
			}

			var amounts []entities.Reserves
			for _, call := range callAggregate {
				if len(call.ReturnData) == 96 && call.Success == true {
					amounts = append(amounts, entities.Reserves{
						Reserve0: new(uint256.Int).SetBytes(call.ReturnData[:32]),
						Reserve1: new(uint256.Int).SetBytes(call.ReturnData[32:64]),
					})
				} else {
					amounts = append(amounts, entities.Reserves{})
				}
			}

			for i2, pool := range batch {
				if values, ok := s.pools.Load(common.HexToAddress(pool.Address)); ok {
					_pool := values.(entities.Reserves)

					if _pool.LastBlock == 0 {
						s.pools.Store(common.HexToAddress(pool.Address), amounts[i2])
					}
				}
			}
		}(&wg, batch)
	}

	wg.Wait()
	s.logger.Info("Reserves loaded")
}

func (s *Service) getBatches(tokens []entities.Pool) (batches [][]entities.Pool) {
	for i := 0; i < len(tokens); i += 100 {
		j := i + 100
		if j > len(tokens) {
			j = len(tokens)
		}

		batches = append(batches, tokens[i:j])
	}

	return batches
}
