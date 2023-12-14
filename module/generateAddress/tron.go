package generate

import (
	"github.com/Francisundermoon/commonChainApi/utils/tron"
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
