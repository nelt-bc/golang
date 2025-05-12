package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func TransformAddress(tronAddress string) (common.Address, error) {
	alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	// Validate input
	if !strings.HasPrefix(tronAddress, "T") {
		return common.Address{}, fmt.Errorf("invalid Tron address format: must start with 'T'")
	}

	// Convert Base58 to bytes
	var value big.Int
	for _, c := range tronAddress {
		value.Mul(&value, big.NewInt(58))
		index := strings.IndexRune(alphabet, c)
		if index == -1 {
			return common.Address{}, fmt.Errorf("invalid character in Tron address: %c", c)
		}
		value.Add(&value, big.NewInt(int64(index)))
	}

	// Convert to bytes
	bytes := value.Bytes()

	// A properly encoded Tron address should decode to 21 bytes (1 byte prefix + 20 bytes address)
	// We need to ensure we have the right amount of bytes
	if len(bytes) < 21 {
		return common.Address{}, fmt.Errorf("invalid Tron address: decoded to %d bytes, expected at least 21", len(bytes))
	}

	// Extract the 20-byte address part, skipping the first byte (0x41 for Tron)
	addressBytes := bytes[1:21]
	ethAddr := "0x" + hex.EncodeToString(addressBytes)

	// Convert to hex with 0x prefix
	return common.HexToAddress(ethAddr), nil
}
