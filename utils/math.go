package utils

import "github.com/holiman/uint256"

var feeFactor = uint256.NewInt(10000)

func GetAmountOut(amountIn, reserveIn, reserveOut, fee *uint256.Int) *uint256.Int {
	amountInWithFee := new(uint256.Int).Mul(amountIn, fee)
	numerator := new(uint256.Int).Mul(amountInWithFee, reserveOut)

	denominator := new(uint256.Int).Mul(reserveIn, feeFactor)
	denominator = denominator.Add(denominator, amountInWithFee)

	return numerator.Div(numerator, denominator)
}
