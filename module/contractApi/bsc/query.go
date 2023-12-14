package bsc

import (
	"context"
	"github.com/Francisundermoon/commonChainApi/contract"
	"github.com/Francisundermoon/commonChainApi/utils"
	Config "github.com/Francisundermoon/commonChainApi/yaml"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	Node = Config.LoadConfig().Node.Bsc
	Usdt = "0x55d398326f99059fF775485246999027B3197955"
)

func GetFeeBalance(userAddress string) (float64, error) {
	conn, err := ethclient.Dial(Node)
	account := common.HexToAddress(userAddress)
	balance, err := conn.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return 0, err
	}
	conn.Client().Close()
	return utils.BigIntToFloat64(balance, 18), nil
}

func GetUsdtBalance(userAddress string) (float64, error) {
	conn, err := ethclient.Dial(Node)
	account := common.HexToAddress(userAddress)
	tokenAddressParam := common.HexToAddress(Usdt)
	instance, err := contract.NewErc20(tokenAddressParam, conn)
	if err != nil {
		return 0, err
	}
	bal, err := instance.BalanceOf(&bind.CallOpts{}, account)
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}
	conn.Client().Close()
	return utils.BigIntToFloat64(bal, int(decimals)), nil
}

func GetCoinBalance(tokenAddress, userAddress string) (float64, error) {
	conn, err := ethclient.Dial(Node)
	account := common.HexToAddress(userAddress)
	tokenAddressParam := common.HexToAddress(tokenAddress)
	instance, err := contract.NewErc20(tokenAddressParam, conn)
	if err != nil {
		return 0, err
	}
	bal, err := instance.BalanceOf(&bind.CallOpts{}, account)
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}
	conn.Client().Close()
	return utils.BigIntToFloat64(bal, int(decimals)), nil
}

func CheckAllowance(tokenAddress, addressA, addressB string) (bool, float64, error) {
	conn, err := ethclient.Dial(Node)
	tokenAddressParam := common.HexToAddress(tokenAddress)
	instance, err := contract.NewErc20(tokenAddressParam, conn)
	if err != nil {
		return false, 0, err
	}
	bal, err := instance.Allowance(&bind.CallOpts{}, common.HexToAddress(addressA), common.HexToAddress(addressB))
	if err != nil {
		return false, 0, err
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	conn.Client().Close()
	return utils.BigIntToFloat64(bal, int(decimals)) > 0, utils.BigIntToFloat64(bal, int(decimals)), nil
}

func CheckUsdtAllowance(addressA, addressB string) (bool, float64, error) {
	conn, err := ethclient.Dial(Node)
	tokenAddressParam := common.HexToAddress(Usdt)
	instance, err := contract.NewErc20(tokenAddressParam, conn)
	if err != nil {
		return false, 0, err
	}
	bal, err := instance.Allowance(&bind.CallOpts{}, common.HexToAddress(addressA), common.HexToAddress(addressB))
	if err != nil {
		return false, 0, err
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	conn.Client().Close()
	return utils.BigIntToFloat64(bal, int(decimals)) > 0, utils.BigIntToFloat64(bal, int(decimals)), nil
}

func GetTotalSupply(tokenAddress string) (int64, int, error) {
	conn, err := ethclient.Dial(Node)
	tokenAddressParam := common.HexToAddress(tokenAddress)
	instance, err := contract.NewErc20(tokenAddressParam, conn)
	if err != nil {
		return 0, 0, err
	}
	bal, err := instance.TotalSupply(&bind.CallOpts{})
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, 0, err
	}
	//fmt.Println("总供应:", bal)
	conn.Client().Close()
	return bal.Int64(), int(decimals), nil
}
