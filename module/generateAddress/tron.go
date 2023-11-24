package generate

import (
	"chargeWithdraw/utils/tron"
)

func ProduceTron() *SingleData {
	sourceData := ProduceEth()
	address := tron.PkToAddress(sourceData.PrivateKey[2:len(sourceData.PrivateKey)])
	return &SingleData{
		Mnemonic:   sourceData.Mnemonic,
		PrivateKey: sourceData.PrivateKey,
		Address:    address,
	}
}
