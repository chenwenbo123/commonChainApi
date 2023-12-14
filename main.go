package main

import (
	_ "github.com/Francisundermoon/commonChainApi/docs"
	logs "github.com/Francisundermoon/commonChainApi/log"
	router "github.com/Francisundermoon/commonChainApi/router"
	Config "github.com/Francisundermoon/commonChainApi/yaml"
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

	//加载数据
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.Init(r, config)

	r.Run(":" + config.System.Port)

}
