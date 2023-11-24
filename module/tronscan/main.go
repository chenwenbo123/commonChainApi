package tronScan

import (
	"chargeWithdraw/db"
	errorModel "chargeWithdraw/error"
	"chargeWithdraw/model"
	"chargeWithdraw/utils"
	"chargeWithdraw/utils/tron"
	Config "chargeWithdraw/yaml"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"sync"
	"time"
)

type TronCurrentBlock struct {
	Num int64
}

type TronScan struct {
	Con *Config.Conf
	Db  *gorm.DB
}

var (
	result        model.UserTransactionHistory
	blockIdData   model.BlockData
	blockScanData model.BlockData
	txRecord      model.History
	apiKeys       = []string{"", "", ""}
)

func (t *TronScan) InitBlock() (int64, error) {
	//随机apikey
	apikey := apiKeys[rand.Intn(len(apiKeys))]
	//获取最新区块
	body, error := utils.Get("https://api.trongrid.io/walletsolidity/getnowblock?visible=true", apikey)
	if error != nil {
		log.Fatalln("波场充提模块启动失败")
		return 0, error
	}
	json.Unmarshal(body, &blockIdData)
	//获取最新区块
	return blockIdData.BlockHeader.RawData.Number, nil
}

func (t *TronScan) LoadModule() {
	t.Con = Config.LoadConfig()
	firstBlockNum, err := t.InitBlock()
	if err != nil {
		panic(errorModel.RequestErr.String())
	}
	var tronCurrentBlock = TronCurrentBlock{Num: firstBlockNum - 3}
	var tronInformMsg = make(chan *model.Inform)
	go t.WatchTron(tronInformMsg)

	c := cron.New()
	_ = c.AddFunc("*/3 * * * * *", func() {
		t.Db = db.InitDb(t.Con)
		var g = sync.WaitGroup{}
		var copyBlockNum int64 = tronCurrentBlock.Num
		tronCurrentBlock.Num++
		g.Add(1)
		go t.ScanBlock(&g, copyBlockNum, tronInformMsg)
		g.Wait()
		d, _ := t.Db.DB()
		d.Close()
	})
	c.Start()
}

func (t *TronScan) InitTaskWatchTron(currentBlock *TronCurrentBlock, cbn int64, msg chan *model.Inform) {

}

func (t *TronScan) JudgeType(chain int, blockNum int64, coin, watchAddress, fromAddress, toAddress string, num float64, txid string) *model.Inform {
	switch watchAddress {
	case fromAddress:
		return &model.Inform{
			Chain:       chain,
			BlockNum:    blockNum,
			CoinName:    coin,
			Type:        "out",
			FromAddress: fromAddress,
			ToAddress:   toAddress,
			Num:         num,
			Txid:        txid,
		}
	case toAddress:
		return &model.Inform{
			Chain:       chain,
			BlockNum:    blockNum,
			CoinName:    coin,
			Type:        "in",
			FromAddress: fromAddress,
			ToAddress:   toAddress,
			Num:         num,
			Txid:        txid,
		}
	default:
		return &model.Inform{
			Chain:       3,
			CoinName:    "",
			Type:        "",
			FromAddress: "",
			ToAddress:   "",
			Num:         0,
			Txid:        "",
		}
	}
}

func (t *TronScan) InitStorageData() []string {
	var data []model.WatchTronAddress
	var ads []string
	t.Db.Select("Address").Model(&data).Pluck("address", &ads)
	return ads
}

func (t *TronScan) ScanBlock(g *sync.WaitGroup, cbn int64, msg chan *model.Inform) {
	startTime := time.Now()
	defer func() {
		// 延迟释放连接
		if error0 := recover(); error0 != nil {
			fmt.Println(error0)
		}
	}()
	fmt.Println("正在扫描区块:", cbn)
	body1, error := utils.TronPost("https://api.trongrid.io/walletsolidity/getblockbynum", t.Con.Tron.ApiKey, map[string]interface{}{"num": cbn, "visible": "true"})
	if error != nil {
		fmt.Println(error)
		fmt.Println("获取指定区块数据出错，", blockIdData.BlockHeader.RawData.Number)
		return
	}
	//修改本地区块完成
	json.Unmarshal(body1, &blockScanData)

	var contractCharge, trxCharge, bothNone = 0, 0, 0
	//查询数据库地址数据
	watchAddresses := t.InitStorageData()
	for _, transaction := range blockScanData.Transactions {
		switch transaction.RawData.Contract[0].Type {
		case "TriggerSmartContract":
			if transaction.RawData.Contract[0].Parameter.Value.ContractAddress == "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" {
				//转账
				if transaction.RawData.Contract[0].Parameter.Value.Data[0:8] == "a9059cbb" {
					toAddress := tron.GetBaseAddress(transaction.RawData.Contract[0].Parameter.Value.Data[30:72])
					num := tron.HexToBigint("0x"+transaction.RawData.Contract[0].Parameter.Value.Data[73:136]) / 1000000
					//fmt.Println(index, "类型：转账", "转出地址："+transaction.RawData.Contract[0].Parameter.Value.OwnerAddress, "目标地址:"+toAddress, "转账金额:"+num, "交易Id:"+transaction.Txid)
					//通知管道
					for _, address := range watchAddresses {
						msg <- t.JudgeType(2, cbn, "USDT", address, transaction.RawData.Contract[0].Parameter.Value.OwnerAddress, toAddress, num, transaction.Txid)
					}
				}
				//授权转账
				if transaction.RawData.Contract[0].Parameter.Value.Data[0:8] == "23b872dd" {
					fromAddress := tron.GetBaseAddress(transaction.RawData.Contract[0].Parameter.Value.Data[30:72])
					toAddress := tron.GetBaseAddress("41" + transaction.RawData.Contract[0].Parameter.Value.Data[96:136])
					num := tron.HexToBigint("0x"+transaction.RawData.Contract[0].Parameter.Value.Data[136:200]) / 1000000
					//fmt.Println(index, "类型：授权转账", "转出地址:"+fromAddress, "目标地址:"+toAddress, "转账金额:"+num, "交易Id:"+transaction.Txid)
					//通知管道
					for _, address := range watchAddresses {
						msg <- t.JudgeType(2, cbn, "USDT", address, fromAddress, toAddress, num, transaction.Txid)
					}
				}
			}
			contractCharge++
		case "TransferContract":
			num := transaction.RawData.Contract[0].Parameter.Value.Amount / 1000000
			//fmt.Println(index, "类型：TRX转账", "转出地址:"+transaction.RawData.Contract[0].Parameter.Value.OwnerAddress, "目标地址:"+transaction.RawData.Contract[0].Parameter.Value.ToAddress, "转账金额:"+num, "交易Id:"+transaction.Txid)
			//通知管道
			for _, address := range watchAddresses {
				msg <- t.JudgeType(2, cbn, "TRX", address, transaction.RawData.Contract[0].Parameter.Value.OwnerAddress, transaction.RawData.Contract[0].Parameter.Value.ToAddress, num, transaction.Txid)
			}
			trxCharge++
		default:
			bothNone++
		}
	}
	spend := time.Since(startTime)
	fmt.Println(spend)
	g.Done()
	//}
}

func (t *TronScan) WatchTron(msg chan *model.Inform) {
	for {
		select {
		case data, ok := <-msg:
			if ok {
				if data.Chain == 2 {

					go utils.SendInform(Config.LoadConfig().Inform.Url, 2, data.BlockNum, data.CoinName, data.ContractAddress, data.Type, data.FromAddress, data.ToAddress, data.Num, data.Txid)

					//fmt.Println("BlockNum:"+fmt.Sprint(data.BlockNum), "CoinName:"+data.CoinName, data.FromAddress, data.ToAddress, "Type:"+data.Type, "Amount:"+data.Num, "Txid:"+data.Txid)
				}
			}
		default:

		}
	}
}
