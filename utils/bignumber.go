package utils

import "math/big"

func BigNumber(num, decimals int64) *big.Int {
	bigNum := big.NewInt(num)
	bigDecimals := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil)

	return big.NewInt(0).Mul(bigNum, bigDecimals)
}
