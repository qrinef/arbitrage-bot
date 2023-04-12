package compute

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/qrinef/arbitrage-bot/utils"
)

type swaps []swap

type swap struct {
	pool              common.Address
	poolFee           *uint256.Int
	amountIn          *uint256.Int
	amountOut         *uint256.Int
	tokenIn           common.Address
	tokenOut          common.Address
	reserveIn         *uint256.Int
	reserveOut        *uint256.Int
	reserveInPending  *uint256.Int
	reserveOutPending *uint256.Int
}

func (s *Service) handler(data *Method, factory common.Address) (swaps swaps, err error) {
	swaps, err = s.getAmountsOut(data.AmountIn, data.Path, factory)
	if err != nil {
		return swaps, err
	}

	return swaps, err
}

func (s *Service) getAmountsOut(amountIn *uint256.Int, path []common.Address, factory common.Address) (amounts swaps, err error) {
	for i := 0; i < len(path)-1; i++ {
		reserveIn, reserveOut, fee, pool, _err := s.poolsService.GetReserves(factory, path[i], path[i+1])
		if _err != nil {
			return amounts, _err
		}

		_amountIn := amountIn
		if i > 0 {
			_amountIn = amounts[i-1].amountOut
		}

		amountOut := utils.GetAmountOut(_amountIn, reserveIn, reserveOut, fee)

		amounts = append(amounts, swap{
			pool:              pool,
			poolFee:           fee,
			amountIn:          _amountIn,
			amountOut:         amountOut,
			tokenIn:           path[i],
			tokenOut:          path[i+1],
			reserveIn:         reserveIn,
			reserveOut:        reserveOut,
			reserveInPending:  new(uint256.Int).Add(reserveIn, _amountIn),
			reserveOutPending: new(uint256.Int).Sub(reserveOut, amountOut),
		})
	}

	return amounts, err
}
