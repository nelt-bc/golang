package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
)

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func TransformAddress(tronAddress string) (common.Address, error) {
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

func EthToTronAddress(ethAddress string) (string, error) {
	ethAddr := common.HexToAddress(ethAddress)
	// 1. Add 0x41 prefix (TRON Mainnet)
	prefix := []byte{0x41}
	addressBytes := ethAddr.Bytes()
	payload := append(prefix, addressBytes...) // 21 bytes

	// 2. Calculate SHA256 checksum
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	checksum := hash2[:4]

	// 3. Append checksum
	fullPayload := append(payload, checksum...) // 25 bytes

	// 4. Base58 encode
	return base58.Encode(fullPayload), nil
}
