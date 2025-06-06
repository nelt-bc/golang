package signature

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"tron/utils"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type AbiArgument struct {
	Type  string
	Value any
}

type Domain struct {
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
}

const EIP721_DOMAIN_TYPE = "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"

func keccak256StructHash(typeDef string, values []AbiArgument) []byte {
	typeMapping := map[string]abi.Argument{}
	typeMapping["int256"] = abi.Argument{Type: abi.Type{T: abi.IntTy, Size: 256}}
	typeMapping["uint256"] = abi.Argument{Type: abi.Type{T: abi.UintTy, Size: 256}}
	typeMapping["bool"] = abi.Argument{Type: abi.Type{T: abi.BoolTy}}
	typeMapping["bytes"] = abi.Argument{Type: abi.Type{T: abi.BytesTy}}
	typeMapping["address"] = abi.Argument{Type: abi.Type{T: abi.AddressTy}}
	typeMapping["bytes32"] = abi.Argument{Type: abi.Type{T: abi.FixedBytesTy, Size: 32}}
	typeMapping["string"] = abi.Argument{Type: abi.Type{T: abi.FixedBytesTy, Size: 32}}

	typeHash := crypto.Keccak256Hash([]byte(typeDef))

	abiValues := []any{}
	types := abi.Arguments{}

	abiValues = append(abiValues, typeHash)
	types = append(types, abi.Argument{Type: abi.Type{T: abi.FixedBytesTy, Size: 32}})

	for _, abiArg := range values {
		types = append(types, typeMapping[abiArg.Type])

		if abiArg.Type == "string" {
			abiValues = append(abiValues, crypto.Keccak256Hash([]byte(abiArg.Value.(string))))
		} else {
			abiValues = append(abiValues, abiArg.Value)
		}
	}

	encoded, err := types.Pack(abiValues...)

	if err != nil {
		log.Fatal(err)
	}

	return crypto.Keccak256(encoded)
}

func SignTypedData(domain Domain, messageHashType string, abiArguments []AbiArgument, privateKey string) ([]byte, error) {
	domainHash := keccak256StructHash(EIP721_DOMAIN_TYPE, []AbiArgument{
		{Type: "string", Value: domain.Name},
		{Type: "string", Value: domain.Version},
		{Type: "uint256", Value: domain.ChainId},
		{Type: "address", Value: domain.VerifyingContract},
	})

	messageHash := keccak256StructHash(messageHashType, abiArguments)
	digest := crypto.Keccak256(
		[]byte("\x19\x01"),
		domainHash,
		messageHash,
	)

	pkBytes, _ := hex.DecodeString(privateKey)
	pk, _ := btcec.PrivKeyFromBytes(pkBytes)

	signature, err := crypto.Sign(digest, pk.ToECDSA())
	if err != nil {
		return []byte{}, err
	}
	return signature, nil
}

func VerifySignature(signature []byte, domain Domain, messageHashType string, abiArguments []AbiArgument, expectedSigner string) (bool, error) {
	if len(signature) != 65 {
		return false, errors.New("invalid signature length")
	}

	domainHash := keccak256StructHash(EIP721_DOMAIN_TYPE, []AbiArgument{
		{Type: "string", Value: domain.Name},
		{Type: "string", Value: domain.Version},
		{Type: "uint256", Value: domain.ChainId},
		{Type: "address", Value: domain.VerifyingContract},
	})
	messageHash := keccak256StructHash(messageHashType, abiArguments)

	digest := crypto.Keccak256(
		[]byte("\x19\x01"),
		domainHash,
		messageHash,
	)

	v := signature[64]
	if v < 27 {
		v += 27
	}
	sig := append(signature[:64], v-27)
	pubKey, err := crypto.SigToPub(digest, sig)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %w", err)
	}

	recovered := crypto.PubkeyToAddress(*pubKey)
	tronAddress, _ := utils.EthToTronAddress(recovered.String())
	return tronAddress == expectedSigner, nil
}

func SplitSignature(signature []byte) (r, s []byte, v byte, err error) {
	if len(signature) != 65 {
		return []byte{}, []byte{}, 0, err
	}

	r = signature[0:32]
	s = signature[32:64]
	v = signature[64] + 27
	return r, s, v, nil
}
