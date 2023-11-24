package model

type Inform struct {
	Chain           int
	BlockNum        int64
	CoinName        string
	ContractAddress string
	Type            string
	FromAddress     string
	ToAddress       string
	Num             float64
	Txid            string
}

type CwLog struct {
	Id              int64
	Chain           int64
	BlockNum        int64
	CoinName        string
	ContractAddress string
	Type            string
	FromAddress     string
	ToAddress       string
	Num             string
	Txid            string
}

type ChargeRecord struct {
	Id          int
	IsOpen      int
	Limit       float64
	EthAddress  string
	TronAddress string
}
