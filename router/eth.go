package router

import (
	"chargeWithdraw/db"
	"chargeWithdraw/model"
	eth2 "chargeWithdraw/module/contractApi/eth"
	tron2 "chargeWithdraw/module/contractApi/tron"
	generate "chargeWithdraw/module/generateAddress"
	"chargeWithdraw/utils"
	"chargeWithdraw/utils/tron"
	Config "chargeWithdraw/yaml"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

var config = Config.LoadConfig()

func HandleErr(err string, ctx *gin.Context) {
	//model.ErrorHandle(err)
	ctx.JSON(200, gin.H{"response": Response{
		Status: 0,
		Msg:    err,
		Data:   nil,
	}})
	return
}

type Response struct {
	Status int64       `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

type QueryTotalSupplyParam struct {
	ChainId      int64  `json:"chain_id"`
	TokenAddress string `json:"token_address"`
}

// QueryTotalSupply 查询Token总发行量
// @title Swagger API
// @version 1.0
// @Tags 查询Token总发行量
// @description  查询Token总发行量
// @BasePath /api/queryTotalsupply
// @Produce  json
// @Param chain_id query int64 true "公链的ID 0代表以太 1代表币安 2代表波场"
// @Param token_address query string true "智能合约地址"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/queryTotalsupply [get]
func QueryTotalSupply(ctx *gin.Context) {
	//var a QueryTotalSupplyParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	token_address := ctx.Query("token_address")
	switch chainId {
	case 0:
		var eth = eth2.GetEthClient(0)
		totalSupply, _, err := eth.GetTotalSupply(token_address)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: totalSupply, Msg: "ok", Status: 1}})
	case 1:
		var eth = eth2.GetEthClient(1)
		totalSupply, _, err := eth.GetTotalSupply(token_address)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: totalSupply, Msg: "ok", Status: 1}})
	case 2:
		var tron = tron2.GetTronClient(0)
		totalSupply, err := tron.QueryTotalSupply(token_address)
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		tron.Client.Conn.Close()
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: totalSupply, Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "error", Status: 0}})
	}

}

type CreateWalletParam struct {
	ChainId int64 `json:"chain_id"`
}

// CreateWallet 创建钱包
// @title Swagger API
// @version 1.0
// @Tags 创建钱包
// @description  创建钱包
// @id 1
// @BasePath /api/createWallet
// @Accept multipart/form-data
// @Produce application/json
// @Param chain_id query int64 true "钱包的ID 0代表以太 1代表币安 2代表波场"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/createWallet [get]
func CreateWallet(ctx *gin.Context) {
	var a CreateWalletParam
	//err := ctx.ShouldBind(&a)
	//if err != nil {
	//	fmt.Println(err)
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	a.ChainId, _ = utils.StrToInt64(ctx.Query("chain_id"))
	fmt.Println(a)
	chainId := a.ChainId
	switch chainId {
	case 0:
		ctx.JSON(200, gin.H{"data": Response{Data: generate.ProduceEth(), Msg: "ok", Status: 1}})
	case 1:
		ctx.JSON(200, gin.H{"data": Response{Data: generate.ProduceEth(), Msg: "ok", Status: 1}})
	case 2:
		ctx.JSON(200, gin.H{"data": Response{Data: generate.ProduceTron(), Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "error", Status: 0}})
	}
}

type QueryFeeBalanceParam struct {
	ChainId     int64  `json:"chain_id"`
	UserAddress string `json:"user_address"`
}

// QueryFeeBalance 查询手续费余额
// @title Swagger API
// @version 1.0
// @Tags 查询手续费币余额
// @description  查询手续费币余额
// @id 0
// @BasePath /api/queryFeeBalance
// @Produce  json
// @Param chain_id query int64 true "钱包的ID 0代表以太 1代表币安 2代表波场"
// @Param user_address query string true "用户地址"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/queryFeeBalance [get]
func QueryFeeBalance(ctx *gin.Context) {
	//var a QueryFeeBalanceParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	userAddress := ctx.Query("user_address")
	switch chainId {
	case 0:
		var eth = eth2.GetEthClient(0)
		data, err := eth.GetFeeBalance(userAddress)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: data, Msg: "ok", Status: 1}})
	case 1:
		var eth = eth2.GetEthClient(1)
		data, err := eth.GetFeeBalance(userAddress)
		if err != nil {
			eth.Client.Close()
			HandleErr(err.Msg, ctx)
			return
		}
		eth.Client.Close()
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: data, Msg: "ok", Status: 1}})
	case 2:
		var tron = tron2.GetTronClient(1)
		data, err := tron.GetFeeBalance(userAddress)
		if err != nil {
			tron.OtherClient.GRPC.Conn.Close()
			HandleErr(err.Msg, ctx)
			return
		}
		tron.OtherClient.GRPC.Conn.Close()
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: data, Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "error", Status: 0}})
	}
}

type QueryTokenBalanceParam struct {
	ChainId      int64  `json:"chain_id"`
	UserAddress  string `json:"user_address"`
	TokenAddress string `json:"token_address"`
}

// QueryTokenBalance 查询代币余额
// @title Swagger API
// @version 1.0
// @Tags 查询合约余额
// @description  查询合约余额
// @BasePath /api/queryTokenBalance
// @Produce  json
// @Param chain_id query int64 true "钱包的ID 0代表以太 1代表币安 2代表波场"
// @Param token_address query string true "合约地址"
// @Param user_address query string true "用户地址"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/queryTokenBalance [get]
func QueryTokenBalance(ctx *gin.Context) {
	//var a QueryTokenBalanceParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	tokenAddress := ctx.Query("token_address")
	userAddress := ctx.Query("user_address")
	//if err != nil {
	//	return
	//}
	switch chainId {
	case 0:
		var eth = eth2.GetEthClient(0)
		data, err := eth.GetTokenBalance(tokenAddress, userAddress)
		if err != nil {
			HandleErr(err.Msg, ctx)
			eth.Client.Close()
			return
		}
		eth.Client.Close()
		ctx.JSON(200, gin.H{"data": Response{Data: data, Msg: "ok", Status: 1}})
	case 1:
		var eth = eth2.GetEthClient(1)
		data, err := eth.GetTokenBalance(tokenAddress, userAddress)
		if err != nil {
			HandleErr(err.Msg, ctx)
			eth.Client.Close()
			return
		}
		eth.Client.Close()
		ctx.JSON(200, gin.H{"data": Response{Data: data, Msg: "ok", Status: 1}})
	case 2:
		var tron = tron2.GetTronClient(1)
		data, err := tron.GetTrc20Balance(tokenAddress, userAddress)
		if err != nil {
			HandleErr(err.Msg, ctx)
			tron.OtherClient.GRPC.Conn.Close()
			return
		}
		tron.OtherClient.GRPC.Conn.Close()
		ctx.JSON(200, gin.H{"data": Response{Data: data, Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "error", Status: 0}})
	}
}

type WithdrawCoinParam struct {
	ChainId      int64   `json:"chain_id"`
	TokenAddress string  `json:"token_address"`
	PrivateKey   string  `json:"private_key"`
	FromAddress  string  `json:"from_address"`
	ToAddress    string  `json:"to_address"`
	Amount       float64 `json:"amount"`
}

// TransferCoin  转账代币
// @title Swagger API
// @version 1.0
// @Tags 转账代币
// @description  转账代币
// @BasePath /api/tranferCoin
// @Produce  json
// @Param chain_id query int64 true "钱包的ID 0代表以太 1代表币安 2代表波场"
// @Param token_address query string true "合约地址"
// @Param private_key query string true "签名私钥"
// @Param from_address query string true "转出地址"
// @Param to_address query string true "目标地址"
// @Param amount query float64 true "提币数量"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/tranferCoin [get]
func TransferCoin(ctx *gin.Context) {
	//var a WithdrawCoinParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	tokenAddress := ctx.Query("token_address")
	privateKey := ctx.Query("private_key")
	toAddress := ctx.Query("to_address")
	amount, _ := utils.StrToFloat64(ctx.Query("amount"))
	fmt.Println("shu", amount)
	fromAddress := ctx.Query("from_address")
	if chainId < 2 {
		var eth = eth2.GetEthClient(int(chainId))
		//加载中间件开始
		//加载中间件结束
		err, txid := eth.TransferErc20(privateKey, fromAddress, toAddress, tokenAddress, amount)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txid}, Msg: "ok", Status: 1}})
	} else if chainId == 2 {
		var tron1 = tron2.GetTronClient(1)
		//加载中间件结束
		txid, err := tron1.TransferTrc20(tokenAddress, fromAddress, toAddress, privateKey, amount)
		tron1.OtherClient.GRPC.Conn.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txid}, Msg: "ok", Status: 1}})
	} else {
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": ""}, Msg: "ok", Status: 0}})
	}
}

type BurnParam struct {
	ChainId      int64   `json:"chain_id"`
	TokenAddress string  `json:"token_address"`
	PrivateKey   string  `json:"private_key"`
	Amount       float64 `json:"amount"`
}

// Burn  销毁
// @title Swagger API
// @version 1.0
// @id 10
// @Tags 销毁
// @description  销毁
// @BasePath /api/burn
// @Produce  json
// @Param chain_id query int64 true "钱包的ID 0代表以太 1代表币安 2代表波场"
// @Param token_address query string true "合约地址"
// @Param private_key query string true "钱包私钥"
// @Param amount query float64 true "销毁数量"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/burn [post]
func Burn(ctx *gin.Context) {
	var a BurnParam
	err := ctx.ShouldBindJSON(&a)
	if err != nil {
		HandleErr(err.Error(), ctx)
		return
	}
	chainId := a.ChainId
	tokenAddress := a.TokenAddress
	privateKey := a.PrivateKey
	amount := a.Amount

	if chainId < 2 {
		var eth = eth2.GetEthClient(int(chainId))
		txid, err := eth.Burn(privateKey, tokenAddress, amount)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txid}, Msg: "ok", Status: 1}})
	} else if chainId == 2 {
		var tron1 = tron2.GetTronClient(0)
		txid, err := tron1.Burn(tokenAddress, privateKey, amount)
		tron1.Client.Conn.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txid}, Msg: "ok", Status: 1}})
	} else {
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": ""}, Msg: "ok", Status: 0}})
	}
}

type MultiSignParam struct {
	CurrentAddress string `json:"current_address"`
	PrivateKey     string `json:"private_key"`
	NewOwner       string `json:"new_owner"`
}

// MultiSign  多签
// @title Swagger API
// @version 1.0
// @Tags 多签
// @description  多签
// @BasePath /api/multiSign
// @Produce  json
// @Param current_address query string true "操作地址"
// @Param private_key query string true "签名私钥"
// @Param new_owner query string true "新权限拥有者"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/multiSign [get]
func MultiSign(ctx *gin.Context) {
	//var a MultiSignParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	currentAddress := ctx.Query("current_address")
	privateKey := ctx.Query("private_key")
	toAddress := ctx.Query("new_owner")
	var tron2 = tron2.GetTronClient(1)
	txid, err1 := tron2.MultiSign(privateKey, currentAddress, toAddress)
	tron2.OtherClient.GRPC.Conn.Close()
	if err1 != nil {
		HandleErr(err1.Msg, ctx)
		return
	}
	ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txid}, Msg: "ok", Status: 1}})
}

// GetOa 获取权限
// @title Swagger API
// @version 1.0
// @Tags 波场获取账户权限Owner和Active
// @description  获取权限所属
// @BasePath /api/getOa
// @Produce  json
// @Param from_address query string true "来源地址"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/getOa [get]
func GetOa(ctx *gin.Context) {
	//var a MultiSignParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	currentAddress := ctx.Query("from_address")
	var tron2 = tron2.GetTronClient(1)
	ac, err1 := tron2.OtherClient.GRPC.GetAccount(currentAddress)
	tron2.OtherClient.GRPC.Conn.Close()
	if err1 != nil {
		return
	}

	ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"owner": tron.GetBaseAddress(common.BytesToAddress(ac.OwnerPermission.Keys[0].Address).String()), "active": tron.GetBaseAddress(common.BytesToAddress(ac.ActivePermission[0].Keys[0].Address).String())}, Msg: "ok", Status: 1}})
}

type MintNftParam struct {
	ChainId    int64   `json:"chain_id"`
	Price      float64 `json:"price"`
	Num        int64   `json:"num"`
	PrivateKey string  `json:"private_key"`
	NftAddress string  `json:"nft_address"`
}

// MintNft  铸造Nft
// @title Swagger API
// @version 1.0
// @Tags 铸造Nft
// @description  铸造Nft
// @BasePath /api/mintNft
// @Produce  json
// @Param chain_id query int64 true "钱包的ID 0代表以太 1代表币安"
// @Param nft_address query string true "nft地址"
// @Param private_key query string true "签名私钥"
// @Param num query int64 true "数量"
// @Param price query float64 true "价格"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/mintNft [post]
func MintNft(ctx *gin.Context) {
	var a MintNftParam
	err := ctx.ShouldBindJSON(&a)
	if err != nil {
		HandleErr(err.Error(), ctx)
		return
	}
	chainId := a.ChainId
	price := a.Price
	num := a.Num
	privateKey := a.PrivateKey
	nftAddress := a.NftAddress
	var eth = eth2.GetEthClient(int(chainId))
	txid, err1 := eth.MintNft(privateKey, nftAddress, num, price)
	eth.Client.Close()
	if err1 != nil {
		HandleErr(err.Error(), ctx)
		return
	}
	ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txid}, Msg: "ok", Status: 1}})
}

type BurnNftParam struct {
	ChainId    int64  `json:"chain_id"`
	PrivateKey string `json:"private_key"`
	NftAddress string `json:"nft_address"`
	TokenId    int64  `json:"token_id"`
}

// BurnNft  销毁Nft
// @title Swagger API
// @version 1.0
// @Tags 销毁Nft
// @description  销毁Nft
// @BasePath /api/burnNft
// @Produce  json
// @Param chain_id query int64 true "钱包的ID 0代表以太 1代表币安"
// @Param nft_address query string true "nft地址"
// @Param private_key query string true "签名私钥"
// @Param token_id query int64 true "nft id"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/burnNft [post]
func BurnNft(ctx *gin.Context) {
	var a BurnNftParam
	err := ctx.ShouldBindJSON(&a)
	if err != nil {
		HandleErr(err.Error(), ctx)
		return
	}
	chainId := a.ChainId
	tokenId := a.TokenId
	privateKey := a.PrivateKey
	nftAddress := a.NftAddress
	var eth = eth2.GetEthClient(int(chainId))
	txid, err1 := eth.BurnNft(privateKey, nftAddress, tokenId)
	eth.Client.Close()
	if err1 != nil {
		HandleErr(err.Error(), ctx)
		return
	}
	ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txid}, Msg: "ok", Status: 1}})
}

// InformReceiver  接收充值提现通知
// @title Swagger API
// @version 1.0
// @Tags 接收充值提现通知
// @description  接收充值提现通知
// @BasePath /api/informReceive
// @Produce  json
// @Param chain query int64 true "钱包的ID 0代表以太 1代表币安"
// @Param block_num query int64 true "区块号"
// @Param coin_name query string true "币名"
// @Param contract_address query string true "合约地址"
// @Param type query string true "类型"
// @Param from_address query string true "发起地址"
// @Param to_address query string true "接收地址"
// @Param num query float64 true "数量"
// @Param txid query string true "txid"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/informReceive [post]
func InformReceiver(ctx *gin.Context) {
	chainId, _ := utils.StrToInt64(ctx.Query("chain"))
	blockNum, _ := utils.StrToInt64(ctx.Query("block_num"))
	num, _ := utils.StrToFloat64(ctx.Query("num"))
	coinName := ctx.Query("coin_name")
	contractAddress := ctx.Query("contract_address")
	typeInfo := ctx.Query("type")
	fromAddress := ctx.Query("from_address")
	toAddress := ctx.Query("to_address")
	txid := ctx.Query("txid")
	var newData = model.CwLog{
		//Id:              0,
		Chain:           chainId,
		BlockNum:        blockNum,
		CoinName:        coinName,
		ContractAddress: contractAddress,
		Type:            typeInfo,
		FromAddress:     fromAddress,
		ToAddress:       toAddress,
		Num:             fmt.Sprint(num),
		Txid:            txid,
	}
	db := db.InitDb(config)
	db.Create(&newData)
	d, _ := db.DB()
	d.Close()
}

// TransferFeeCoin   转账手续费币
// @title Swagger API
// @version 1.0
// @Tags 转账手续费币
// @description  转账手续费币
// @BasePath /api/transferFeeCoin
// @Produce  json
// @Param chain_id query int64 true "钱包的ID 0代表以太 1代表币安 2代表波场"
// @Param private_key query string true "私钥"
// @Param from_address query string true "转出地址"
// @Param to_address query string true "转入地址"
// @Param amount query float64 true "数量"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/transferFeeCoin [get]
func TransferFeeCoin(ctx *gin.Context) {
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	num, _ := utils.StrToFloat64(ctx.Query("amount"))
	pk := ctx.Query("private_key")
	fromAddress := ctx.Query("from_address")
	toAddress := ctx.Query("to_address")
	if chainId < 2 {
		cli := eth2.GetEthClient(int(chainId))
		bal, err := cli.GetFeeBalance(fromAddress)
		if err != nil {
			HandleErr("network error", ctx)
			return
		}
		if bal < num {
			HandleErr("insufficent balance", ctx)
			return
		}
		txId, err := cli.TransferFeeCoin(pk, fromAddress, toAddress, num)
		if err != nil {
			HandleErr("错误", ctx)
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txId}, Msg: "ok", Status: 1}})
	} else {
		cli := tron2.GetTronClient(1)
		bal, err := cli.GetFeeBalance(fromAddress)
		if err != nil {
			HandleErr("network error", ctx)
			return
		}
		if bal < num {
			HandleErr("insufficient balance", ctx)
			return
		}
		txid, err := cli.TransferFeeCoin(pk, fromAddress, toAddress, num)
		if err != nil {
			HandleErr("network error", ctx)
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]string{"txid": txid}, Msg: "ok", Status: 1}})
	}
}
