package compute

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/qrinef/arbitrage-bot/services/pools"
)

type Service struct {
	poolsService *pools.Service

	routers map[common.Address]common.Address
}

func NewService(poolsService *pools.Service) *Service {
	return &Service{
		poolsService: poolsService,
		routers: map[common.Address]common.Address{
			common.HexToAddress(""): common.HexToAddress(""),
			common.HexToAddress(""): common.HexToAddress(""),
		},
	}
}

func (s *Service) Compute(tx *types.Transaction) (err error) {
	if tx == nil || tx.To() == nil {
		return errors.New("tx / to address is incorrect")
	}

	factory, ok := s.routers[*tx.To()]
	if !ok {
		return errors.New("tx router not found")
	}

	method, err := s.decodeMethod(tx.Data(), uint256.MustFromBig(tx.Value()))
	if err != nil {
		return err
	}

	swapsData, err := s.handler(method, factory)
	if err != nil {
		return err
	}

	strategy, err := s.strategy(swapsData, factory)
	if err != nil {
		return err
	}

	_ = strategy
	return err
}
