package main

import (
	_ "chargeWithdraw/docs"
	logs "chargeWithdraw/log"
	tronScan "chargeWithdraw/module/tronscan"
	"chargeWithdraw/router"
	Config "chargeWithdraw/yaml"
	"github.com/gin-gonic/gin"
)

// @title ChainBaby
// @version 1.0
// @description ChainBaby
// @termsOfService  http://swagger.io/terms/

// @contact.name ChainBaby
// @contact.email ChainBaby

func main() {

	//init the configfile
	config := Config.LoadConfig()
	//init the log file
	logs.InitLog()
	//load the module of scanBlock
	var a = tronScan.TronScan{
		Con: nil,
		Db:  nil,
	}
	a.LoadModule()
	//////load the api
	//var es etherScan.ScanClient
	//var esi etherScan.ScanClientInterface
	//esi = &es
	////load ethreum
	//go esi.RealInitEthreum(0)
	//////load bsc
	//go esi.RealInitEthreum(1)

	//加载数据
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.Init(r, config)

	r.Run(":" + config.System.Port)

}
