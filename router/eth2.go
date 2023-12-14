package router

import (
	eth2 "github.com/Francisundermoon/commonChainApi/module/contractApi/eth"
	tron2 "github.com/Francisundermoon/commonChainApi/module/contractApi/tron"
	"github.com/Francisundermoon/commonChainApi/utils"
	"github.com/Francisundermoon/commonChainApi/utils/eth"
	tron3 "github.com/Francisundermoon/commonChainApi/utils/tron"
	"github.com/gin-gonic/gin"
)

// CheckAllowance  检查授权
// @title Swagger API
// @version 1.0
// @Tags 查询Token 检查授权
// @description   查询检查授权
// @BasePath /api/checkAllowance
// @Produce  json
// @Param chain_id query int64 true "公链的ID 0代表以太 1代表币安 2代表波场"
// @Param token_address query string true "代币地址"
// @Param owner_address query string true "来源地址"
// @Param spender_address query string true "消费者地址"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/checkAllowance [get]
func CheckAllowance(ctx *gin.Context) {
	//var a QueryTotalSupplyParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	token_address := ctx.Query("token_address")
	from_address := ctx.Query("owner_address")
	to_address := ctx.Query("spender_address")
	switch chainId {
	case 0:
		var eth = eth2.GetEthClient(0)
		isApproved, num, err := eth.CheckAllowance(token_address, from_address, to_address)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isApproved": isApproved, "num": num}, Msg: "ok", Status: 1}})
	case 1:
		var eth = eth2.GetEthClient(1)
		isApproved, num, err := eth.CheckAllowance(token_address, from_address, to_address)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isApproved": isApproved, "num": num}, Msg: "ok", Status: 1}})

	case 2:
		var tron = tron2.GetTronClient(0)
		//totalSupply, err := tron.CheckAllowance(token_address, from_address, to_address)
		//if err != nil {
		//	HandleErr(err.Msg, ctx)
		//	return
		//}
		err, isApproved, num := tron3.RealCheckAllowance(tron.Client, token_address, from_address, to_address)
		tron.Client.Conn.Close()
		if err != nil {
			HandleErr("network error", ctx)
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isApproved": isApproved, "num": num}, Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "error", Status: 0}})
	}

}

// ApproveCoin  代币授权
// @title Swagger API
// @version 1.0
// @Tags  代币授权
// @description   代币授权
// @BasePath /api/approveCoin
// @Produce  json
// @Param chain_id query int64 true "公链的ID 0代表以太 1代表币安 2代表波场"
// @Param token_address query string true "代币地址"
// @Param private_key query string true "私钥"
// @Param spender_address query string true "消费者地址"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/approveCoin [get]
func ApproveCoin(ctx *gin.Context) {
	//var a QueryTotalSupplyParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	token_address := ctx.Query("token_address")
	private_key := ctx.Query("private_key")
	to_address := ctx.Query("spender_address")
	switch chainId {
	case 0:
		var eth = eth2.GetEthClient(0)
		txid, err := eth.Approve(private_key, token_address, to_address)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isApproved": true, "txid": txid}, Msg: "ok", Status: 1}})
	case 1:
		var eth = eth2.GetEthClient(1)
		txid, err := eth.Approve(private_key, token_address, to_address)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isApproved": true, "txid": txid}, Msg: "ok", Status: 1}})
	case 2:
		var tron = tron2.GetTronClient(1)
		//totalSupply, err := tron.CheckAllowance(token_address, from_address, to_address)
		//if err != nil {
		//	HandleErr(err.Msg, ctx)
		//	return
		//}
		txid, err := tron.Approve(private_key, token_address, to_address)
		tron.OtherClient.GRPC.Conn.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isApproved": true, "num": txid}, Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "error", Status: 0}})
	}

}

// TransferFromCoin  授权转账代币
// @title Swagger API
// @version 1.0
// @Tags  授权转账代币
// @description   授权转账代币
// @BasePath /api/transferFromCoin
// @Produce  json
// @Param chain_id query int64 true "公链的ID 0代表以太 1代表币安 2代表波场"
// @Param token_address query string true "代币地址"
// @Param amount query float64 true "代币数量"
// @Param from_address query string true "来源地址"
// @Param to_address query string true "目标地址"
// @Param private_key query string true "消费地址私钥"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/transferFromCoin [get]
func TransferFromCoin(ctx *gin.Context) {
	//var a QueryTotalSupplyParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	token_address := ctx.Query("token_address")
	private_key := ctx.Query("private_key")
	to_address := ctx.Query("to_address")
	from_address := ctx.Query("from_address")
	amount, _ := utils.StrToFloat64(ctx.Query("amount"))
	switch chainId {
	case 0:
		var eth = eth2.GetEthClient(0)
		txid, err := eth.TransferFromCoin(private_key, token_address, from_address, to_address, amount)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isTransfered": true, "txid": txid}, Msg: "ok", Status: 1}})
	case 1:
		var eth = eth2.GetEthClient(1)
		txid, err := eth.TransferFromCoin(private_key, token_address, from_address, to_address, amount)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isTransfered": true, "txid": txid}, Msg: "ok", Status: 1}})
	case 2:
		var tron = tron2.GetTronClient(1)
		//totalSupply, err := tron.CheckAllowance(token_address, from_address, to_address)
		//if err != nil {
		//	HandleErr(err.Msg, ctx)
		//	return
		//}
		txid, err := tron.TransferFromCoin(private_key, token_address, from_address, to_address, amount)
		tron.OtherClient.GRPC.Conn.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"isTransfered": true, "txid": txid}, Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "error", Status: 0}})
	}

}

// IssueSimpleCoin 发布普通代币
// @title Swagger API
// @version 1.0
// @Tags  发布普通代币
// @description   发布普通代币
// @BasePath /api/issueSimpleCoin
// @Produce  json
// @Param chain_id query int64 true "公链的ID 0代表以太 1代表币安"
// @Param name query string true "代币名称"
// @Param decimal query int64 true "代币精度"
// @Param symbol query string true "代币符号"
// @Param total_supply query int64 true "总供应量"
// @Param private_key query string true "私钥"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/issueSimpleCoin [get]
func IssueSimpleCoin(ctx *gin.Context) {
	//var a QueryTotalSupplyParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	private_key := ctx.Query("private_key")
	name := ctx.Query("name")
	decimal, _ := utils.StrToInt64(ctx.Query("decimal"))
	symbol := ctx.Query("symbol")
	total_supply, _ := utils.StrToInt64(ctx.Query("total_supply"))
	switch chainId {
	case 0:
		var eth = eth2.GetEthClient(0)
		coinAddr, err := eth.IssueCoin(private_key, name, symbol, decimal, total_supply)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"coinAddr": coinAddr, "isIssued": true}, Msg: "ok", Status: 1}})
	case 1:
		var eth = eth2.GetEthClient(1)
		coinAddr, err := eth.IssueCoin(private_key, name, symbol, decimal, total_supply)
		eth.Client.Close()
		if err != nil {
			HandleErr(err.Msg, ctx)
			return
		}
		if err != nil {
			return
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"coinAddr": coinAddr, "isIssued": true}, Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "invalid public chain", Status: 0}})
	}

}

// GetBaseAddress  私钥转地址
// @title Swagger API
// @version 1.0
// @Tags  私钥转地址
// @description   私钥转地址
// @BasePath /api/getBaseAddress
// @Produce  json
// @Param chain_id query int64 true "公链的ID 0代表以太 1代表币安 2代表波场"
// @Param private_key query string true "私钥"
// @Success 200 {object} Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/getBaseAddress [get]
func GetBaseAddress(ctx *gin.Context) {
	//var a QueryTotalSupplyParam
	//err := ctx.ShouldBindJSON(&a)
	//if err != nil {
	//	HandleErr(err.Error(), ctx)
	//	return
	//}
	chainId, _ := utils.StrToInt64(ctx.Query("chain_id"))
	private_key := ctx.Query("private_key")
	switch chainId {
	case 0:
		result := eth.PkToAddress(private_key)
		if result == "" {
			ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"msg": "error privateKey", "status": 0}, Msg: "ok", Status: 1}})
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"Addr": result, "status": "ok"}, Msg: "ok", Status: 1}})
	case 1:
		result := eth.PkToAddress(private_key)
		if result == "" {
			ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"msg": "error privateKey", "status": 0}, Msg: "ok", Status: 1}})
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"Addr": result, "status": "ok"}, Msg: "ok", Status: 1}})
	case 2:
		result := tron3.PkToAddress(private_key)
		if result == "" {
			ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"msg": "error privateKey", "status": 0}, Msg: "ok", Status: 1}})
		}
		ctx.JSON(200, gin.H{"data": Response{Data: map[string]interface{}{"Addr": result, "status": "ok"}, Msg: "ok", Status: 1}})
	default:
		ctx.JSON(200, gin.H{"data": Response{Data: 0, Msg: "invalid public chain", Status: 0}})
	}

}
