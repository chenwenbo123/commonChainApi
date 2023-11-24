package utils

import (
	"math"
	"math/big"
)

func BigIntToFloat64(num *big.Int, decimal int) float64 {
	str, _ := new(big.Float).SetInt(num).Float64()
	return str / math.Pow10(decimal)
}
