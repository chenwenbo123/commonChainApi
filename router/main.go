package router

import (
	_ "chargeWithdraw/docs"
	"chargeWithdraw/middleware"
	Config "chargeWithdraw/yaml"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	//加载接口文档
	url := ginSwagger.URL(c.System.Router + ":" + c.System.Port + "/api/docs/doc.json") // The url pointing to API definition
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

}
