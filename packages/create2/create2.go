package create2

import (
	"encoding/hex"
	"log"
	"math/big"
	"strings"
	"tron/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func CalculateAddress(tronDeployer string, salt int64, bytecode string) (common.Address, error) {
	hexString := strings.TrimPrefix(bytecode, "0x")
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}

	hashByteCode := crypto.Keccak256(bytes)
	ethDeployer, err := utils.TransformAddress(tronDeployer)

	if err != nil {
		return common.Address{}, err
	}

	payload := []byte("\x41")
	payload = append(payload, ethDeployer.Bytes()...)
	payload = append(payload, padTo32Bytes(big.NewInt(salt).Bytes())...)
	payload = append(payload, hashByteCode...)

	address := crypto.Keccak256(payload)
	return common.BytesToAddress(address), nil
}

func padTo32Bytes(b []byte) []byte {
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b) // right-align (left-pad with zeros)
	return padded
}
