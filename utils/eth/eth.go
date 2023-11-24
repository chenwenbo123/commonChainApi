package eth

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"math"
	"math/big"
)

func PkToAddress(pKey string) string {
	privateKey, err := crypto.HexToECDSA(pKey)
	if err != nil {
		return ""
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return ""
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress.String()
}

func Float64ToBigNum(num float64, decimal int) *big.Int {
	//这是处理位数的代码段
	tenDecimal := big.NewFloat(math.Pow(10, float64(decimal)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, big.NewFloat(num)).Int(&big.Int{})
	return convertAmount
}
