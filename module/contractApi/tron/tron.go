package tron

import (
	"fmt"

	"github.com/Francisundermoon/commonChainApi/utils"
	"github.com/Francisundermoon/commonChainApi/utils/eth"
	"github.com/Francisundermoon/commonChainApi/utils/tron"
	Config "github.com/Francisundermoon/commonChainApi/yaml"
	"github.com/JFJun/trx-sign-go/grpcs"
	"github.com/JFJun/trx-sign-go/sign"
	"github.com/fbsobreira/gotron-sdk/pkg/abi"
	address2 "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"google.golang.org/grpc"
	"math"
	"math/big"
	"strings"
)

var (
	Node   = Config.LoadConfig().Node.Tron
	ApiKey = Config.LoadConfig().Tron.ApiKey
	Usdt   = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
)

func Init() (*grpcs.Client, error) {

	conn, err := grpcs.NewClient(Node)
	conn.GRPC.Start(grpc.WithInsecure())
	conn.GRPC.SetAPIKey(ApiKey)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func QueryTotalSupply(contractAdress string) (string, error) {
	conn, err := Init()
	if err != nil {
		return "", err
	}
	tx, err := conn.GRPC.TriggerConstantContract("TX8h6Df74VpJsXF6sTDz1QJsq3Ec8dABc3",
		contractAdress,
		"totalSupply()", ``)
	if err != nil {
		return "", err
	}
	a, _ := conn.GRPC.GetContractABI(contractAdress)
	arg, err := abi.GetParser(a, "totalSupply")
	if err != nil {
		return "", err
	}
	result := map[string]interface{}{}
	arg.UnpackIntoMap(result, tx.ConstantResult[0])
	conn.GRPC.Conn.Close()
	return fmt.Sprint(result[""]), nil
}

func GetFeeBalance(userAddress string) (float64, error) {
	conn, err := Init()
	if err != nil {
		return 0, err
	}
	acc, errr := conn.GetTrxBalance(userAddress)
	if errr != nil {
		return 0, err
	}
	trxbalance := float64(acc.GetBalance()) / math.Pow10(6)
	conn.GRPC.Conn.Close()
	return trxbalance, nil
}

func GetUsdtBalance(userAddress string) (float64, error) {
	conn, err := Init()
	if err != nil {
		return 0, err
	}
	result, err := conn.GetTrc20Balance(userAddress, Usdt)
	if err != nil {
		return 0, nil
	}
	result1 := utils.BigIntToFloat64(result, 6)
	conn.GRPC.Conn.Close()
	return result1, nil
}

func TransferTrc20(from, toAddress, pk string, num float64) (string, error) {
	conn, err := Init()
	if err != nil {
		return "", nil
	}
	tx, err := conn.TransferTrc20(from, toAddress, Usdt, big.NewInt(int64(num*math.Pow10(6))), 50000000)
	if err != nil {
		return "", nil
	}
	signTx, err := tron.SignTransaction(tx.Transaction, pk)
	if err != nil {
		return "", nil
	}
	err = conn.BroadcastTransaction(signTx)
	if err != nil {
		return "", nil
	}
	txid := strings.TrimLeft(common.BytesToHexString(tx.GetTxid()), "0x")
	conn.GRPC.Conn.Close()
	return txid, nil
}

func MultiSign(pk, fromAddress, newOwner string) (string, error) {
	conn, err := Init()
	if err != nil {
		return "", nil
	}
	ta, err := conn.GRPC.GetAccount(fromAddress)
	if err != nil {
		return "", nil
	}
	newAddress, _ := address2.Base58ToAddress(newOwner)
	ta.OwnerPermission.Keys[0].Address = newAddress
	ta.ActivePermission[0].Keys[0].Address = newAddress
	transaction, err := tron.UpdateAccountPermission(conn.GRPC, fromAddress, ta.OwnerPermission, nil, ta.ActivePermission)
	if err != nil {
		return "", nil
	}
	signTx, err := sign.SignTransaction(transaction.Transaction, pk)
	if err != nil {
		return "", nil
	}
	err = conn.BroadcastTransaction(signTx)
	if err != nil {
		return "", nil
	}
	txid := strings.TrimLeft(common.BytesToHexString(transaction.GetTxid()), "0x")
	conn.GRPC.Conn.Close()
	return txid, nil
}

func TransferFeeCoin(pk, fromAddress, toAddress string, amount float64) (string, error) {
	conn, err := Init()
	if err != nil {
		return "", nil
	}
	tx, err := conn.Transfer(fromAddress, toAddress, int64(amount*math.Pow10(6)))
	if err != nil {
		return "", nil
	}
	signTx, err := sign.SignTransaction(tx.Transaction, pk)
	if err != nil {
		return "", nil
	}
	err = conn.BroadcastTransaction(signTx)
	if err != nil {
		return "", nil
	}
	txid := strings.TrimLeft(common.BytesToHexString(tx.GetTxid()), "0x")
	conn.GRPC.Conn.Close()
	return txid, nil
}

func ApproveUsdt(pk, toAddress string) (string, error) {
	conn, err := Init()
	if err != nil {
		return "", nil
	}
	tx, err := conn.GRPC.TRC20Approve(tron.PkToAddress(pk), toAddress, Usdt, big.NewInt(int64(math.Pow10(18)*9999999)), 300000)
	if err != nil {
		return "", nil
	}
	signTx, err := sign.SignTransaction(tx.Transaction, pk)
	if err != nil {
		return "", nil
	}
	err = conn.BroadcastTransaction(signTx)
	if err != nil {
		return "", nil
	}
	txid := strings.TrimLeft(common.BytesToHexString(tx.GetTxid()), "0x")
	conn.GRPC.Conn.Close()
	return txid, nil
}

func TransferFromUsdt(pk, tokenAddress, fromAddress, toAddress string, amount float64) (string, error) {
	conn, err := Init()
	if err != nil {
		return "", nil
	}
	decimal, err := conn.GRPC.TRC20GetDecimals(tokenAddress)
	if err != nil {
		return "", nil
	}
	err, txid := tron.TransferFromCoin(conn.GRPC, pk, tokenAddress, fromAddress, toAddress, eth.Float64ToBigNum(amount, int(decimal.Int64())))
	if err != nil {
		return "", nil
	}
	conn.GRPC.Conn.Close()
	return txid, nil
}
