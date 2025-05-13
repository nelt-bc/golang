package signature

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"tron/utils"
)

func TestSignature() {
	tronTokenAddr := "TNPFZPscg3sXjRLRq5pj2KuWZRBP3E7815"
	tronAaAddr := "TNS6uN2UWkkapzc4iDhC7W5xH82kkzBYgC"
	tronOwnerAddr := "TJPncMxDwoApkXjVU4oU6T28W5aUbWnGRG"
	tronSpenderAddr := "TCvRUR6dqDRjwtJihAbaBwwh9wM6zsV3aL"

	ethTokenAddr, _ := utils.TransformAddress(tronTokenAddr)
	ethAaAddress, _ := utils.TransformAddress(tronAaAddr)
	ethOwnerAddr, _ := utils.TransformAddress(tronOwnerAddr)
	ethSpenderAddr, _ := utils.TransformAddress(tronSpenderAddr)

	name := "AA"
	version := "1"
	nonce := int64(1)
	deadline := time.Now().Unix() + 3600
	chainId := big.NewInt(3448148188)

	signature, err := SignTypedData(Domain{
		Name:              name,
		Version:           version,
		ChainId:           chainId,
		VerifyingContract: ethAaAddress,
	}, "Permit(address token,address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)",
		[]AbiArgument{
			{Type: "address", Value: ethTokenAddr},
			{Type: "address", Value: ethOwnerAddr},
			{Type: "address", Value: ethSpenderAddr},
			{Type: "uint256", Value: utils.BigNumber(1, 18)},
			{Type: "uint256", Value: big.NewInt(nonce)},
			{Type: "uint256", Value: big.NewInt(deadline)},
		},
		os.Getenv("PK"),
	)
	if err != nil {
		log.Fatalf("Invalid signature %v ::", err)
	}
	fmt.Printf("Deadline: %d\n", deadline)

	r, s, v, err := SplitSignature(signature)
	if err != nil {
		log.Fatalf("Invalid signature %v ::", err)
	}
	fmt.Printf("r: 0x%s\n", hex.EncodeToString(r))
	fmt.Printf("s: 0x%s\n", hex.EncodeToString(s))
	fmt.Printf("v: %d\n", v)
}

func TestVerifySignature() {
	tronTokenAddr := "TNPFZPscg3sXjRLRq5pj2KuWZRBP3E7815"
	tronAaAddr := "TNS6uN2UWkkapzc4iDhC7W5xH82kkzBYgC"
	tronOwnerAddr := "TJPncMxDwoApkXjVU4oU6T28W5aUbWnGRG"
	tronSpenderAddr := "TCvRUR6dqDRjwtJihAbaBwwh9wM6zsV3aL"

	ethTokenAddr, _ := utils.TransformAddress(tronTokenAddr)
	ethAaAddress, _ := utils.TransformAddress(tronAaAddr)
	ethOwnerAddr, _ := utils.TransformAddress(tronOwnerAddr)
	ethSpenderAddr, _ := utils.TransformAddress(tronSpenderAddr)

	name := "AA"
	version := "1"
	nonce := int64(1)
	deadline := time.Now().Unix() + 3600
	chainId := big.NewInt(3448148188)

	domain := Domain{
		Name:              name,
		Version:           version,
		ChainId:           chainId,
		VerifyingContract: ethAaAddress,
	}
	messageHashType := "Permit(address token,address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"
	abiArgs := []AbiArgument{
		{Type: "address", Value: ethTokenAddr},
		{Type: "address", Value: ethOwnerAddr},
		{Type: "address", Value: ethSpenderAddr},
		{Type: "uint256", Value: utils.BigNumber(1, 18)},
		{Type: "uint256", Value: big.NewInt(nonce)},
		{Type: "uint256", Value: big.NewInt(deadline)},
	}

	signature, err := SignTypedData(
		domain,
		messageHashType,
		abiArgs,
		os.Getenv("PK"),
	)
	if err != nil {
		log.Fatalf("Error %v ::", err)
	}

	_, err = VerifySignature(signature, domain, messageHashType, abiArgs, tronOwnerAddr)
	if err != nil {
		log.Fatalf("Error %v ::", err)
	}

	fmt.Print("Valid signature")
}
