package compute

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type Strategy struct {
	tokenIn     common.Address
	tokenOut    common.Address
	pools       []common.Address
	fees        []*uint256.Int
	tradeAmount *uint256.Int
	dirtyProfit *uint256.Int
}

func (s *Service) strategy(swapsData swaps, factory common.Address) (strategy Strategy, err error) {
	return strategy, err
}
