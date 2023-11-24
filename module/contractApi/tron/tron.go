package tron

import (
	"chargeWithdraw/model"
	"chargeWithdraw/utils"
	"chargeWithdraw/utils/eth"
	"chargeWithdraw/utils/tron"
	Config "chargeWithdraw/yaml"
	"encoding/json"
	"fmt"
	"github.com/JFJun/trx-sign-go/grpcs"
	"github.com/JFJun/trx-sign-go/sign"
	"github.com/fbsobreira/gotron-sdk/pkg/abi"
	address2 "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"google.golang.org/grpc"
	"math"
	"math/big"
	"strings"
)

type ErrorMessage struct {
	Status int64
	Msg    string
}

type TronClient struct {
	Node        string
	ApiKey      string
	Client      *client.GrpcClient
	OtherClient *grpcs.Client
}

type ITronClient interface {
	QueryTotalSupply(contractAdress string) string
	GetFeeBalance(userAddress string) (float64, *model.ErrorMessage)
	Transfer(contractAdress, from, toAddress, pk string, num float64) string
	Burn(contractAddress, pk string, num float64) string
	MultiSign(pk, newOwner string) (string, *model.ErrorMessage)
}

func (t *TronClient) Init(index int) {
	switch index {
	case 0:
		t.Client = client.NewGrpcClient(t.Node)
		t.Client.SetAPIKey(t.ApiKey)
		t.Client.Start(grpc.WithInsecure())
		//t.Client.Start(grpc.WithInsecure())
	case 1:
		var err error
		t.OtherClient, err = grpcs.NewClient(t.Node)
		if err != nil {
			model.ErrorHandle(err, 0)
		}
		t.OtherClient.GRPC.Start(grpc.WithInsecure())
		t.OtherClient.GRPC.SetAPIKey(t.ApiKey)

		//t.OtherClient.GRPC.Start(grpc.WithInsecure())
	}
}

func (t *TronClient) QueryTotalSupply(contractAdress string) (string, *model.ErrorMessage) {
	tx, err := t.Client.TriggerConstantContract("TX8h6Df74VpJsXF6sTDz1QJsq3Ec8dABc3",
		contractAdress,
		"totalSupply()", ``)
	if err != nil {
		model.ErrorHandle(err, 0)
		return "", &model.ErrorMessage{Msg: fmt.Sprint(err)}
	}
	a, _ := t.Client.GetContractABI(contractAdress)
	arg, err := abi.GetParser(a, "totalSupply")
	if err != nil {
		model.ErrorHandle(err, 0)
		return "", &model.ErrorMessage{Msg: fmt.Sprint(err)}
	}
	result := map[string]interface{}{}
	arg.UnpackIntoMap(result, tx.ConstantResult[0])
	return fmt.Sprint(result[""]), nil
}

func (t *TronClient) GetFeeBalance(userAddress string) (float64, *model.ErrorMessage) {
	acc, errr := t.OtherClient.GetTrxBalance(userAddress)
	if errr != nil {
		return 0, model.ErrorHandle(errr, 0)
	}
	trxbalance := float64(acc.GetBalance()) / math.Pow10(6)
	return trxbalance, nil
}

func (t *TronClient) GetTrc20Balance(tokenAddress, userAddress string) (float64, *model.ErrorMessage) {
	result, err := t.OtherClient.GetTrc20Balance(userAddress, tokenAddress)
	if err != nil {
		return 0, model.ErrorHandle(err, 0)
	}
	result1 := utils.BigIntToFloat64(result, 6)
	return result1, nil
}

func (t *TronClient) TransferTrc20(contractAdress, from, toAddress, pk string, num float64) (string, *model.ErrorMessage) {
	tx, err := t.OtherClient.TransferTrc20(from, toAddress, contractAdress, big.NewInt(int64(num*math.Pow10(6))), 50000000)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	signTx, err := tron.SignTransaction(tx.Transaction, pk)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	err = t.OtherClient.BroadcastTransaction(signTx)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	txid := strings.TrimLeft(common.BytesToHexString(tx.GetTxid()), "0x")
	return txid, nil
}

func (t *TronClient) Burn(contractAddress, pk string, num float64) (string, *model.ErrorMessage) {
	decimal, _ := t.Client.TRC20GetDecimals(contractAddress)
	data := []interface{}{map[string]string{"uint256": fmt.Sprint(big.NewInt(int64(num * math.Pow10(int(decimal.Int64())))))}}
	addressJson, _ := json.Marshal(data)
	tx, err := t.Client.TriggerContract(tron.PkToAddress(pk), contractAddress, "burn(uint256)", string(addressJson), 50000000, 0, "", 0)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	signTx, err := tron.SignTransaction(tx.Transaction, pk)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	_, err = t.Client.Broadcast(signTx)
	if err != nil {
		return "", model.ErrorHandle(err, 0)
	}
	txid := strings.TrimLeft(common.BytesToHexString(tx.GetTxid()), "0x")
	return txid, nil
}

func (t *TronClient) MultiSign(pk, fromAddress, newOwner string) (string, *model.ErrorMessage) {
	ta, err := t.OtherClient.GRPC.GetAccount(fromAddress)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	newAddress, _ := address2.Base58ToAddress(newOwner)
	ta.OwnerPermission.Keys[0].Address = newAddress
	ta.ActivePermission[0].Keys[0].Address = newAddress
	transaction, err := tron.UpdateAccountPermission(t.OtherClient.GRPC, fromAddress, ta.OwnerPermission, nil, ta.ActivePermission)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	signTx, err := sign.SignTransaction(transaction.Transaction, pk)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	err = t.OtherClient.BroadcastTransaction(signTx)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	txid := strings.TrimLeft(common.BytesToHexString(transaction.GetTxid()), "0x")
	fmt.Println(txid)
	return txid, nil
}

func GetTronClient(index int64) *TronClient {
	var config = Config.LoadConfig()
	var tron = TronClient{
		Node:   config.Node.Tron,
		ApiKey: config.Tron.ApiKey,
	}
	switch index {
	case 0:
		tron.Init(0)
		return &tron
	case 1:
		tron.Init(1)
		tron.OtherClient.GRPC.Start(grpc.WithInsecure())
		return &tron
	}
	return nil
}

func (t *TronClient) TransferFeeCoin(pk, fromAddress, toAddress string, amount float64) (string, *model.ErrorMessage) {
	tx, err := t.OtherClient.Transfer(fromAddress, toAddress, int64(amount*math.Pow10(6)))
	if err != nil {
		return "", &model.ErrorMessage{Msg: "network error"}
	}
	signTx, err := sign.SignTransaction(tx.Transaction, pk)
	if err != nil {
		return "", &model.ErrorMessage{Msg: "network error"}
	}
	err = t.OtherClient.BroadcastTransaction(signTx)
	if err != nil {
		return "", &model.ErrorMessage{Msg: "network error"}
	}
	txid := strings.TrimLeft(common.BytesToHexString(tx.GetTxid()), "0x")
	return txid, nil
}

func (t *TronClient) Approve(pk, tokenAddress, toAddress string) (string, *model.ErrorMessage) {
	tx, err := t.OtherClient.GRPC.TRC20Approve(tron.PkToAddress(pk), toAddress, tokenAddress, big.NewInt(int64(math.Pow10(18)*9999999)), 300000)
	if err != nil {
		return "", &model.ErrorMessage{Msg: err.Error()}
	}
	signTx, err := sign.SignTransaction(tx.Transaction, pk)
	if err != nil {
		return "", &model.ErrorMessage{Msg: err.Error()}
	}
	err = t.OtherClient.BroadcastTransaction(signTx)
	if err != nil {
		return "", &model.ErrorMessage{Msg: err.Error()}
	}
	txid := strings.TrimLeft(common.BytesToHexString(tx.GetTxid()), "0x")
	return txid, nil
}

func (t *TronClient) TransferFromCoin(pk, tokenAddress, fromAddress, toAddress string, amount float64) (string, *model.ErrorMessage) {
	decimal, err := t.OtherClient.GRPC.TRC20GetDecimals(tokenAddress)
	if err != nil {
		return "", &model.ErrorMessage{Msg: err.Error()}
	}
	err, txid := tron.TransferFromCoin(t.OtherClient.GRPC, pk, tokenAddress, fromAddress, toAddress, eth.Float64ToBigNum(amount, int(decimal.Int64())))
	if err != nil {
		return "", &model.ErrorMessage{Msg: err.Error()}
	}
	return txid, nil
}

func (t *TronClient) IssueSimpleCoin(pk, tokenAddress, fromAddress, toAddress string, amount float64) (string, *model.ErrorMessage) {
	decimal, err := t.OtherClient.GRPC.TRC20GetDecimals(tokenAddress)
	if err != nil {
		return "", &model.ErrorMessage{Msg: err.Error()}
	}
	err, txid := tron.TransferFromCoin(t.OtherClient.GRPC, pk, tokenAddress, fromAddress, toAddress, eth.Float64ToBigNum(amount, int(decimal.Int64())))
	if err != nil {
		return "", &model.ErrorMessage{Msg: err.Error()}
	}
	return txid, nil
}
