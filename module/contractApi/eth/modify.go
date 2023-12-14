package eth

import (
	"context"
	"github.com/Francisundermoon/commonChainApi/contract"
	"github.com/Francisundermoon/commonChainApi/utils/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
)

func TransferFeeCoin(pk, fromAddress, toAddress string, num float64) (string, error) {
	conn, err := ethclient.Dial(Node)
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return "", err
	}
	nonce, err := conn.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		return "", err
	}
	value := big.NewInt(int64(math.Pow10(18) * num)) // in wei (1 eth) 	// in units
	gasPrice, err := conn.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	var data []byte
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), value, 80000, gasPrice, data)

	chainID, err := conn.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = conn.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	conn.Client().Close()
	return signedTx.Hash().String(), nil
}

func TransferErc20(pKey string, from, to string, token string, num float64) (string, error) {
	conn, err := ethclient.Dial(Node)
	privateKey, err := crypto.HexToECDSA(pKey)
	var ctx = context.Background()
	cd, err := conn.ChainID(ctx)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {
		return "", err
	}
	//部署合约
	addr, err := contract.NewErc20(common.HexToAddress(token), conn)
	if err != nil {
		return "", err
	}
	_, decimal, _ := GetTotalSupply(token)
	//这是处理位数的代码段
	tenDecimal := big.NewFloat(math.Pow(10, float64(decimal)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, big.NewFloat(num)).Int(&big.Int{})

	tx, err := addr.Transfer(auth, common.HexToAddress(to), convertAmount)
	if err != nil {
		return "", err
	}
	conn.Client().Close()
	return tx.Hash().String(), nil
}

func TransferUsdt(pKey string, from, to string, num float64) (string, error) {
	conn, err := ethclient.Dial(Node)
	privateKey, err := crypto.HexToECDSA(pKey)
	var ctx = context.Background()
	cd, err := conn.ChainID(ctx)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {
		return "", err
	}
	//部署合约
	addr, err := contract.NewErc20(common.HexToAddress(Usdt), conn)
	if err != nil {
		return "", err
	}
	_, decimal, _ := GetTotalSupply(Usdt)
	//这是处理位数的代码段
	tenDecimal := big.NewFloat(math.Pow(10, float64(decimal)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, big.NewFloat(num)).Int(&big.Int{})

	tx, err := addr.Transfer(auth, common.HexToAddress(to), convertAmount)
	if err != nil {
		return "", err
	}
	conn.Client().Close()
	return tx.Hash().String(), nil
}

func ApproveUsdt(pk, toAddress string) (string, error) {
	//格式化
	//创建身份，需要私钥
	conn, err := ethclient.Dial(Node)
	privateKey, err := crypto.HexToECDSA(pk)
	var ctx = context.Background()
	cd, err := conn.ChainID(ctx)
	if err != nil {

		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {

		return "", err
	}
	//部署合约
	addr, err := contract.NewErc20(common.HexToAddress(Usdt), conn)
	if err != nil {

		return "", err
	}
	_, decimal, _ := GetTotalSupply(Usdt)
	tx, err := addr.Approve(auth, common.HexToAddress(toAddress), eth.Float64ToBigNum(99999999999999999999999999, decimal))
	if err != nil {

		return "", err
	}
	conn.Client().Close()
	return tx.Hash().String(), nil
}

func TransferFromUsdt(pk, fromAddress, toAddress string, amount float64) (string, error) {
	//格式化
	//创建身份，需要私钥
	conn, err := ethclient.Dial(Node)
	privateKey, err := crypto.HexToECDSA(pk)
	var ctx = context.Background()
	cd, err := conn.ChainID(ctx)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {
		return "", err
	}
	//部署合约
	addr, err := contract.NewErc20(common.HexToAddress(Usdt), conn)
	if err != nil {
		return "", err
	}
	_, decimal, _ := GetTotalSupply(Usdt)
	tx, err := addr.TransferFrom(auth, common.HexToAddress(fromAddress), common.HexToAddress(toAddress), eth.Float64ToBigNum(amount, decimal))
	if err != nil {
		return "", err
	}
	conn.Client().Close()
	return tx.Hash().String(), nil
}
