package router

import (
	_ "github.com/Francisundermoon/commonChainApi/docs"
	"github.com/Francisundermoon/commonChainApi/middleware"
	Config "github.com/Francisundermoon/commonChainApi/yaml"
	"github.com/gin-gonic/gin"
)

func Init(g *gin.Engine, c *Config.Conf) {
	r := g.Group("/api")
	//跨域配置
	r.Use(middleware.Cors())
	//分组路由
	{
		r.GET("/queryTotalsupply", QueryTotalSupply)
		r.GET("/createWallet", CreateWallet)
		r.GET("/queryFeeBalance", QueryFeeBalance)
		r.GET("/queryTokenBalance", QueryTokenBalance)
		r.GET("/tranferCoin", TransferCoin)
		r.POST("/burn", Burn)
		r.GET("/multiSign", MultiSign)
		r.POST("/mintNft", MintNft)
		r.POST("/burnNft", BurnNft)
		r.GET("/informReceive", InformReceiver)
		r.GET("/getOa", GetOa)
		r.GET("/transferFeeCoin", TransferFeeCoin)
		r.GET("/checkAllowance", CheckAllowance)
		r.GET("/approveCoin", ApproveCoin)
		r.GET("/transferFromCoin", TransferFromCoin)
		r.GET("/issueSimpleCoin", IssueSimpleCoin)
		r.GET("/getBaseAddress", GetBaseAddress)
	}
}
