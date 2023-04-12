package poolsCollector

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/imroc/req/v3"
	"github.com/qrinef/arbitrage-bot/entities"
	"github.com/qrinef/arbitrage-bot/services/storage"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"time"
)

type Service struct {
	storageService *storage.Service
	logger         *zap.SugaredLogger
}

type Logs struct {
	Result []struct {
		Address     string
		Data        string
		Topics      []string
		BlockNumber string
	}
}

func NewService(storageService *storage.Service, logger *zap.SugaredLogger) *Service {
	return &Service{
		storageService: storageService,
		logger:         logger,
	}
}

func (s *Service) Start() {
	totalNewPools := 0

	for {
		s.logger.With("fromBlock", s.getLatestBlock()).Info("poolsCollectors: loading new pools")

		if cnt := s.handlerLogs(s.getLogs()); cnt > 0 {
			totalNewPools += cnt
			time.Sleep(time.Millisecond * 200)
		} else {
			break
		}
	}

	s.logger.With("newPools", totalNewPools).Info("poolsCollectors: finished")
}

func (s *Service) getLogs() (logs *Logs) {
	resp, err := req.C().SetTimeout(30*time.Second).R().
		SetHeader("Content-Type", "application/json").
		SetQueryParamsAnyType(map[string]interface{}{
			"apiKey":    s.storageService.Config.ScanApiKey,
			"module":    "logs",
			"action":    "getLogs",
			"topic0":    "0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9",
			"fromBlock": s.getLatestBlock(),
		}).
		SetRetryCount(5).
		SetRetryFixedInterval(1 * time.Second).
		SetSuccessResult(&logs).
		Get(s.storageService.Config.ScanApiEndpoint)

	if err != nil {
		s.logger.With(err).Fatal("poolsCollectors: logs not received")
	}

	if !resp.IsSuccessState() {
		s.logger.With(err).Fatal("poolsCollectors: logs not received")
	}

	return logs
}

func (s *Service) handlerLogs(logs *Logs) int {
	var tokensPools []entities.Pool

	for _, log := range logs.Result {
		blockNumber, err := uint256.FromHex(log.BlockNumber)
		if err != nil {
			s.logger.With(err).Warnf("poolsCollectors: block number \"%s\" no hex", log.BlockNumber)
			continue
		}

		tokensPools = append(tokensPools, entities.Pool{
			BlockNumber: int(blockNumber.Uint64()),
			Factory:     log.Address,
			Address:     common.HexToAddress(log.Data[26:66]).Hex(),
			Token0:      common.HexToAddress(log.Topics[1][26:]).Hex(),
			Token1:      common.HexToAddress(log.Topics[2][26:]).Hex(),
		})
	}

	if len(tokensPools) > 0 {
		tx := s.storageService.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&tokensPools)
		return int(tx.RowsAffected)
	}

	return 0
}

func (s *Service) getLatestBlock() (latestBlock int) {
	s.storageService.Db.
		Model(entities.Pool{}).
		Select("COALESCE(MAX(block_number), 0) AS latest_block").
		Find(&latestBlock)

	return latestBlock + 1
}
