package eth

import (
	"chargeWithdraw/contract"
	"chargeWithdraw/contract/nft"
	"chargeWithdraw/model"
	"chargeWithdraw/utils"
	"chargeWithdraw/utils/eth"
	Config "chargeWithdraw/yaml"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
)

type EthClient struct {
	Node   string
	Client *ethclient.Client
}

type IEthClient interface {
	Init()
	GetTotalSupply(tokenAddress string) (int64, int, *model.ErrorMessage)
	TransferErc20(pKey string, to string, token string, num float64) (*model.ErrorMessage, string)
	Burn(pk, tokenAddress string, num float64) (string, *model.ErrorMessage)
	GetFeeBalance(userAddress string) (float64, *model.ErrorMessage)
	GetTokenBalance(tokenAddress, userAddress string) (float64, *model.ErrorMessage)
	MintNft(pk, nftAddress string, num int64, price float64) (string, *model.ErrorMessage)
	BurnNft(pk, nftAddress string, tokenid int64) (string, *model.ErrorMessage)
}

func (e *EthClient) Init() {
	conn, err := ethclient.Dial(e.Node)
	if err != nil {
		model.ErrorHandle(err, 0)
	}
	e.Client = conn
}
func (e *EthClient) GetTotalSupply(tokenAddress string) (int64, int, *model.ErrorMessage) {
	tokenAddressParam := common.HexToAddress(tokenAddress)
	instance, err := contract.NewErc20(tokenAddressParam, e.Client)
	if err != nil {
		return 0, 0, model.ErrorHandle(err, 1)
	}
	bal, err := instance.TotalSupply(&bind.CallOpts{})
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, 0, model.ErrorHandle(err, 1)
	}
	//fmt.Println("总供应:", bal)
	return bal.Int64(), int(decimals), nil
}

func (e *EthClient) GetFeeBalance(userAddress string) (float64, *model.ErrorMessage) {
	account := common.HexToAddress(userAddress)
	balance, err := e.Client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return 0, model.ErrorHandle(err, 0)
	}
	return utils.BigIntToFloat64(balance, 18), nil
}

func (e *EthClient) GetTokenBalance(tokenAddress, userAddress string) (float64, *model.ErrorMessage) {
	tokenAddressParam := common.HexToAddress(tokenAddress)
	instance, err := contract.NewErc20(tokenAddressParam, e.Client)
	if err != nil {
		return 0, model.ErrorHandle(err, 1)
	}
	bal, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(userAddress))
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, model.ErrorHandle(err, 1)
	}
	return utils.BigIntToFloat64(bal, int(decimals)), nil
}

func (e *EthClient) TransferErc20(pKey string, from, to string, token string, num float64) (*model.ErrorMessage, string) {
	privateKey, err := crypto.HexToECDSA(pKey)
	var ctx = context.Background()
	cd, err := e.Client.ChainID(ctx)
	if err != nil {
		fmt.Println("报错1")
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {
		fmt.Println("报错2")
		return model.ErrorHandle(err, 1), ""
	}
	//部署合约
	addr, err := contract.NewErc20(common.HexToAddress(token), e.Client)
	if err != nil {
		fmt.Println("报错3")
		return model.ErrorHandle(err, 1), ""
	}
	_, decimal, _ := e.GetTotalSupply(token)
	//这是处理位数的代码段
	tenDecimal := big.NewFloat(math.Pow(10, float64(decimal)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, big.NewFloat(num)).Int(&big.Int{})

	tx, err := addr.Transfer(auth, common.HexToAddress(to), convertAmount)
	if err != nil {
		fmt.Println("报错4")
		return model.ErrorHandle(err, 1), ""
	}
	return nil, tx.Hash().String()
}

func (e *EthClient) Burn(pk, tokenAddress string, num float64) (string, *model.ErrorMessage) {
	//格式化
	//创建身份，需要私钥
	privateKey, err := crypto.HexToECDSA(pk)
	var ctx = context.Background()
	cd, err := e.Client.ChainID(ctx)
	if err != nil {
		fmt.Println("报错1")
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {
		fmt.Println("报错2")
		return "", model.ErrorHandle(err, 1)
	}
	//部署合约
	addr, err := contract.NewErc20(common.HexToAddress(tokenAddress), e.Client)
	if err != nil {
		fmt.Println("报错3")
		return "", model.ErrorHandle(err, 1)
	}
	_, decimal, _ := e.GetTotalSupply(tokenAddress)
	tx, err := addr.Burn(auth, eth.Float64ToBigNum(num, decimal))
	if err != nil {
		fmt.Println("报错4")
		return "", model.ErrorHandle(err, 1)
	}
	return tx.Hash().String(), nil
}

func (e *EthClient) MintNft(pk, nftAddress string, num int64, price float64) (string, *model.ErrorMessage) {
	//格式化
	//创建身份，需要私钥
	privateKey, err := crypto.HexToECDSA(pk)
	var ctx = context.Background()
	cd, _ := e.Client.ChainID(ctx)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(int64(price * math.Pow10(18)))
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	//部署合约
	addr, err := nft.NewNft(common.HexToAddress(nftAddress), e.Client)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	tx, err := addr.Mint(auth, big.NewInt(num))
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	return tx.Hash().String(), nil
}

func (e *EthClient) BurnNft(pk, nftAddress string, tokenid int64) (string, *model.ErrorMessage) {
	//格式化
	//创建身份，需要私钥
	privateKey, err := crypto.HexToECDSA(pk)
	var ctx = context.Background()
	cd, _ := e.Client.ChainID(ctx)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	//部署合约
	addr, err := nft.NewNft(common.HexToAddress(nftAddress), e.Client)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	tx, err := addr.TransferFrom(auth, common.HexToAddress(eth.PkToAddress(pk)), common.HexToAddress("0xc3fdde08ff90197d6504ff790c6402efbd4bd28f"), big.NewInt(tokenid))
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}
	return tx.Hash().String(), nil
}

func GetEthClient(index int) *EthClient {
	var config = Config.LoadConfig()
	switch index {
	case 0:
		var eth = EthClient{Node: config.Node.Ethereum}
		eth.Init()
		return &eth
	case 1:
		var eth = EthClient{Node: config.Node.Bsc}
		eth.Init()
		return &eth
	}
	return nil
}

func (e *EthClient) TransferFeeCoin(pk, fromAddress, toAddress string, num float64) (string, *model.ErrorMessage) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {

	}
	nonce, err := e.Client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {

	}
	value := big.NewInt(int64(math.Pow10(18) * num)) // in wei (1 eth) 	// in units
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {

	}
	var data []byte
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), value, 80000, gasPrice, data)

	chainID, err := e.Client.NetworkID(context.Background())
	if err != nil {

	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {

	}

	err = e.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {

	}
	return signedTx.Hash().String(), nil
}

func (e *EthClient) CheckAllowance(tokenAddress, addressA, addressB string) (bool, float64, *model.ErrorMessage) {
	tokenAddressParam := common.HexToAddress(tokenAddress)
	instance, err := contract.NewErc20(tokenAddressParam, e.Client)
	if err != nil {
		return false, 0, model.ErrorHandle(err, 1)
	}
	bal, err := instance.Allowance(&bind.CallOpts{}, common.HexToAddress(addressA), common.HexToAddress(addressB))
	if err != nil {
		return false, 0, model.ErrorHandle(err, 1)
	}
	//fmt.Println("总供应:", bal)
	return utils.BigIntToFloat64(bal, 18) > 0, utils.BigIntToFloat64(bal, 18), nil
}

func (e *EthClient) Approve(pk, tokenAddress, toAddress string) (string, *model.ErrorMessage) {
	//格式化
	//创建身份，需要私钥
	privateKey, err := crypto.HexToECDSA(pk)
	var ctx = context.Background()
	cd, err := e.Client.ChainID(ctx)
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	//部署合约
	addr, err := contract.NewErc20(common.HexToAddress(tokenAddress), e.Client)
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	_, decimal, _ := e.GetTotalSupply(tokenAddress)
	tx, err := addr.Approve(auth, common.HexToAddress(toAddress), eth.Float64ToBigNum(99999999999999999999999999, decimal))
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	return tx.Hash().String(), nil
}

func (e *EthClient) TransferFromCoin(pk, tokenAddress, fromAddress, toAddress string, amount float64) (string, *model.ErrorMessage) {
	//格式化
	//创建身份，需要私钥
	privateKey, err := crypto.HexToECDSA(pk)
	var ctx = context.Background()
	cd, err := e.Client.ChainID(ctx)
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	//部署合约
	addr, err := contract.NewErc20(common.HexToAddress(tokenAddress), e.Client)
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	_, decimal, _ := e.GetTotalSupply(tokenAddress)
	tx, err := addr.TransferFrom(auth, common.HexToAddress(fromAddress), common.HexToAddress(toAddress), eth.Float64ToBigNum(amount, decimal))
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	return tx.Hash().String(), nil
}

func (e *EthClient) IssueCoin(pk, name, symbol string, decimal, amount int64) (string, *model.ErrorMessage) {
	//格式化
	//创建身份，需要私钥
	privateKey, err := crypto.HexToECDSA(pk)
	var ctx = context.Background()
	cd, err := e.Client.ChainID(ctx)
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetInt64(cd.Int64()))
	auth.Value = new(big.Int).SetInt64(0)
	if err != nil {

		return "", model.ErrorHandle(err, 1)
	}
	//部署合约
	addr, _, _, err := contract.DeployErc20(auth, e.Client, name, uint8(decimal), big.NewInt(amount), symbol)
	if err != nil {
		return "", model.ErrorHandle(err, 1)
	}

	return addr.String(), nil
}
