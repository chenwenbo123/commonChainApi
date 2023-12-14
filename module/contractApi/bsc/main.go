package bsc

import (
	"context"
	"github.com/Francisundermoon/commonChainApi/model"
	"github.com/Francisundermoon/commonChainApi/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetFeeBalance(userAddress string) (float64, *model.ErrorMessage) {
	conn, err := ethclient.Dial("https://bsc-dataseed1.ninicoin.io")
	account := common.HexToAddress(userAddress)
	balance, err := conn.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return 0, model.ErrorHandle(err, 0)
	}
	return utils.BigIntToFloat64(balance, 18), nil
}
