package etherScan

import (
	"context"
	"fmt"
	"github.com/Francisundermoon/commonChainApi/db"
	"github.com/Francisundermoon/commonChainApi/model"
	"github.com/Francisundermoon/commonChainApi/utils"
	Config "github.com/Francisundermoon/commonChainApi/yaml"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/robfig/cron"
	//"github.com/ethereum/go-ethereum/core/types"
	ethereum_watcher "github.com/HydroProtocol/ethereum-watcher"
	"github.com/HydroProtocol/ethereum-watcher/plugin"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/gommon/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"io"
	"math"
	"math/big"
	"os"
	"strings"
	"sync"
)

type ScanClient struct {
	Db   *gorm.DB
	Con  *Config.Conf
	EthC *ethclient.Client
}

type ScanClientInterface interface {
	InitEthWatchModule(index int)
	//ReceiveMsgModule(msg chan *TransactionDetails)
	InitStorageData(index int) []string
	Main(index int, storageBlockNum *StorageBlockNum, msg chan *TransactionDetails, cbn int64, wg *sync.WaitGroup)
	ScanBlock(index int, blockNum int64, client *ethclient.Client, msg chan *TransactionDetails)
	GetCoinList(index int) []Coin
	RealInitEthreum(index int)
}

type StorageBlockNum struct {
	Eth int64
}

func (s *ScanClient) InitEthWatchModule(index int) {
	con := Config.LoadConfig()
	s.Con = con
	var storageBlockNum = StorageBlockNum{Eth: InitBlockNum(0)}
	//storageBlockNum.Eth = InitBlockNum(0)
	var a = make(chan *TransactionDetails)
	//监听通知输出
	//go s.ReceiveMsgModule(a)
	//创建扫块定时任务
	c := cron.New()
	_ = c.AddFunc("*/14 * * * * *", func() {
		//初始化以太坊客户端
		s.EthC = NewClient(GetNodeInfo(index))
		s.Db = db.InitDb(con)
		var wg = sync.WaitGroup{}
		wg.Add(1)
		if index == 0 {
			go s.Main(index, &storageBlockNum, a, storageBlockNum.Eth, &wg)
		} else {
			go s.Main(index, &storageBlockNum, a, storageBlockNum.Eth, &wg)
		}
		wg.Wait()
		d, _ := s.Db.DB()
		d.Close()
		s.EthC.Client().Close()
	})
	c.Start()
	fmt.Println("以太坊监听初始化成功")
}

func (s *ScanClient) ScanBlock(index int, blockNum int64, client *ethclient.Client, msg chan *TransactionDetails) {
	fmt.Println("正在扫描块:", blockNum)
	//计时
	//defer func() {
	//	if err := recover(); err != nil {
	//		fmt.Println("错误是", err)
	//		fmt.Println("以太坊网络错误")
	//	}
	//}()
	//startT := time.Now()
	blockNumber := big.NewInt(blockNum)
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		model.ErrorHandle(err, 0)
	}
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		model.ErrorHandle(err, 0)
		fmt.Println("网络错误")
		return
	}
	//fmt.Println("成功获取交易数:", block.Transactions().Len())
	coinList := s.GetCoinList(index)
	for _, tx := range block.Transactions() {
		if tx.To() != nil {
			for _, coin := range coinList {
				if tx.To().Hex() == coin.ContractAddress {
					validRecord := DecodeData(tx.Data(), coin.Decimals)
					if validRecord == nil {
						continue
					}
					if validRecord.FromAddress == "" {
						from, err := types.Sender(types.LatestSignerForChainID(chainId), tx)
						if err != nil {
							fmt.Println(err) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
						}
						validRecord.FromAddress = from.Hex()
					}
					WriteInform(msg, &TransactionDetails{
						BlockNum:        int64(block.Number().Uint64()),
						CoinName:        coin.Coin,
						ContractAddress: coin.ContractAddress,
						FromAddress:     validRecord.FromAddress,
						ToAddress:       validRecord.ToAddress,
						Amount:          validRecord.Amount,
						Txid:            tx.Hash().Hex(),
					})
				}
			}
		}
	}
	//tc := time.Since(startT) //计算耗时
	//fmt.Println(blockNum, "扫描完成:", block.Transactions().Len())
	//fmt.Printf("耗时= %v\n", tc)
	//fmt.Println("当前携程数量", runtime.NumGoroutine())
}

// WriteInform 写入监听通知
func WriteInform(msg chan *TransactionDetails, t *TransactionDetails) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("错误是", err)
		}
	}()
	msg <- t
}

func DecodeData(data []byte, decimals int) *TransactionDetails {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("错误是", err)
		}
	}()
	contractABI, err := abi.JSON(strings.NewReader(GetLocalABI("contract/erc20.abi")))
	if err != nil {
		model.ErrorHandle(err, 0)
		return nil
	}
	//fmt.Println(contractABI.)
	methodSigData := data[:4]
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		model.ErrorHandle(err, 0)
		return nil
	}

	inputsSigData := data[4:]
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		model.ErrorHandle(err, 0)
		return nil
	}
	switch method.Name {
	case "transfer":
		return &TransactionDetails{
			FromAddress: "",
			ToAddress:   inputsMap["_to"].(common.Address).String(),
			Amount:      float64(inputsMap["_value"].(*big.Int).Int64()) / math.Pow10(decimals),
		}
	case "transferFrom":
		return &TransactionDetails{
			FromAddress: inputsMap["_from"].(common.Address).String(),
			ToAddress:   inputsMap["_to"].(common.Address).String(),
			Amount:      float64(inputsMap["_value"].(*big.Int).Int64()) / math.Pow10(decimals),
		}
	default:
		return nil
	}
}

func (s *ScanClient) GetCoinList(index int) []Coin {
	var (
		data     []CoinList
		coinList []Coin
	)
	s.Db.Model(&data).Where("chain=?", index).Select("coin", "contract_address", "decimals").Scan(&coinList)
	return coinList
}

func (s *ScanClient) Main(index int, storageBlockNum *StorageBlockNum, msg chan *TransactionDetails, cbn int64, wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("报错是", err)
		}
	}()
	s.ScanBlock(index, cbn, s.EthC, msg)
	storageBlockNum.Eth++
	wg.Done()
}

func GetLocalABI(path string) string {
	abiFile, err := os.Open(path)
	if err != nil {
		model.ErrorHandle(err, 0)
	}
	defer abiFile.Close()

	result, err := io.ReadAll(abiFile)
	if err != nil {
		model.ErrorHandle(err, 0)
	}
	return string(result)
}

func GetNodeInfo(index int) string {
	switch index {
	case 0:
		return Config.LoadConfig().Node.Ethereum
	case 1:
		return Config.LoadConfig().Node.Ethereum
	}
	return ""
}

// InitBlockNum 初始化区块号
func InitBlockNum(index int) int64 {
	C := NewClient(GetNodeInfo(index))
	//firstStart := Config.LoadConfig()
	header, err := C.HeaderByNumber(context.Background(), nil)
	if err != nil {
		panic("以太坊模块调用失败")
	}
	return header.Number.Int64() - 3
}

// WatchEthAddress 加载监控地址列表
type WatchEthAddress struct {
	Id      int
	Address string
	Type    int
}

func (s *ScanClient) InitStorageData(index int) []string {
	var data []WatchEthAddress
	var ads []string
	s.Db.Select("Address").Model(&data).Where("type=?", index).Pluck("address", &ads)
	//d, _ := db.DB()
	//d.Close()
	return ads
}

// ReceiveMsgModule 输出充提通知
//func (s *ScanClient) ReceiveMsgModule(msg chan *TransactionDetails) {
//	var op Output
//	for {
//		select {
//		case data, ok := <-msg:
//			if ok {
//				ads := s.InitStorageData()
//				for _, a := range ads {
//					if data.FromAddress == a {
//						op.ContratcAddress = data.ContractAddress
//						op.Num = data.Amount
//						op.FromAddress = data.FromAddress
//						op.ToAddress = data.ToAddress
//						op.CoinName = data.CoinName
//						op.Txid = data.Txid
//						op.BlockNum = data.BlockNum
//						op.Type = "out"
//						var countG = sync.WaitGroup{}
//						countG.Add(1)
//						go utils.SendInform(&countG, Config.LoadConfig().Inform.Url, 0, op.BlockNum, op.CoinName, op.ContratcAddress, op.Type, op.FromAddress, op.ToAddress, op.Num, op.Txid)
//						countG.Wait()
//					}
//					if data.ToAddress == a {
//						op.ContratcAddress = data.ContractAddress
//						op.Num = data.Amount
//						op.FromAddress = data.FromAddress
//						op.ToAddress = data.ToAddress
//						op.CoinName = data.CoinName
//						op.Txid = data.Txid
//						op.BlockNum = data.BlockNum
//						op.Type = "in"
//						var countG = sync.WaitGroup{}
//						countG.Add(1)
//						go utils.SendInform(&countG, Config.LoadConfig().Inform.Url, 0, op.BlockNum, op.CoinName, op.ContratcAddress, op.Type, op.FromAddress, op.ToAddress, op.Num, op.Txid)
//						countG.Wait()
//					}
//				}
//			}
//		default:
//
//		}
//	}
//}

func NewClient(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func GetLatestBlock(client *ethclient.Client) int64 {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(header.Number.String())
	return header.Number.Int64()
}

type CoinList struct {
	Id              int
	Chain           int
	Coin            string
	ContractAddress string
}

type Coin struct {
	Coin            string
	ContractAddress string
	Decimals        int
}

// Output 输出消息
type Output struct {
	BlockNum        int64
	Type            string
	CoinName        string
	ContratcAddress string
	FromAddress     string
	ToAddress       string
	Num             float64
	Txid            string
}

type TransactionDetails struct {
	BlockNum        int64
	CoinName        string
	ContractAddress string
	FromAddress     string
	ToAddress       string
	Amount          float64
	Txid            string
}

func (s *ScanClient) RealInitEthreum(index int) {
	//initTheConfigFileAndDatabase
	con := Config.LoadConfig()
	s.Con = con
	s.Db = db.InitDb(con)

	api := s.Con.Node.Ethereum
	watchContract := "0xdac17f958d2ee523a2206206994597c13d831ec7"
	if index == 1 {
		api = s.Con.Node.Bsc
		watchContract = "0x55d398326f99059ff775485246999027b3197955"
	}

	addresses := s.InitStorageData(index)
	w := ethereum_watcher.NewHttpBasedEthWatcher(context.Background(), api)
	fmt.Println("当前节点:", api, watchContract, addresses[0] == "0x0d0707963952f2fba59dd06f2b425ace40b492fe")
	w.RegisterTxReceiptPlugin(plugin.NewERC20TransferPlugin(
		func(token, from, to string, amount decimal.Decimal, isRemove bool) {

			if token == watchContract {
				fromAddress, toAddress := common.HexToAddress(from), common.HexToAddress(to)
				//fmt.Println(token, fromAddress.String(), toAddress.String())
				//if fromAddress.String() == "0x0d0707963952f2fba59dd06f2b425ace40b492fe" || toAddress.String() == "0x0d0707963952f2fba59dd06f2b425ace40b492fe" {
				//	fmt.Println("有记录")
				//}
				//logrus.Infof("New ERC20 Transfer >> token(%s), %s -> %s, amount: %s, isRemoved: %t",
				//	token, from, to, amount, isRemove)
				//var countG sync.WaitGroup
				//countG.Add(1)
				//go func() {
				newAmount, _ := amount.Float64()
				var decimalNew int
				if index == 0 {
					decimalNew = 6
				} else {
					decimalNew = 18
				}
				realAmount := decimal.NewFromFloat(newAmount / math.Pow10(decimalNew)).Round(4)
				amount, _ := realAmount.Float64()

				for _, currentAddr := range addresses {
					//fmt.Println("当前地址", currentAddr)
					//transferOut
					//fmt.Println(indexAddress, currentAddr == fromAddress.String(), currentAddr == toAddress.String())
					//if currentAddr == fromAddress.String() {
					if common.HexToAddress(currentAddr).String() == fromAddress.String() {
						fmt.Println(currentAddr, "out")
						utils.SendInform(con.Inform.Url, int64(index), 0, "USDT", "USDT", "out", fromAddress.String(), toAddress.String(), amount, "")
					}

					//	transferIn
					//if currentAddr == toAddress.String() {
					if common.HexToAddress(currentAddr).String() == toAddress.String() {
						fmt.Println(currentAddr, "in")
						utils.SendInform(con.Inform.Url, int64(index), 0, "USDT", "USDT", "in", fromAddress.String(), toAddress.String(), amount, "")
					}
				}
				//}()
				//countG.Wait()
			}
		},
	))

	fmt.Println("ETH BSC SCAN BLOCK MODULE STARTED")
	w.RunTillExit()
	//	endTheLoop
	//d, _ := s.Db.DB()
	//d.Close()
}
