package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

func BigNumber(num, decimals int64) *big.Int {
	bigNum := big.NewInt(num)
	bigDecimals := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil)

	return big.NewInt(0).Mul(bigNum, bigDecimals)
}

func Int64ToBytes(i int64) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, i) // or binary.LittleEndian
	return buf.Bytes()
}
