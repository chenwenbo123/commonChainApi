package model

type TronGenerate struct {
	Address      string
	PrivateKey   string
	Index        int
	WatchAddress string
}

type name struct {
	Data []string
	Meta []string
}

type OneData struct {
	Ret []RetEle
}

type RetEle struct {
	ContractRet string `json:"contractRet"`
	Fee         int64  `json:"fee"`
}

type ParameterEle struct {
	Data            string  `json:"data"`
	OwnerAddress    string  `json:"owner_address"`
	ContractAddress string  `json:"contract_address"`
	Amount          float64 `json:"amount"`
	ToAddress       string  `json:"to_address"`
}

type Parameter struct {
	Value   ParameterEle `json:"value"`
	TypeUrl string       `json:"type_url"`
}
type ContractEle struct {
	Parameter Parameter `json:"parameter"`
	Type      string    `json:"type"`
}

type RawData struct {
	Data          string        `json:"data"`
	Contract      []ContractEle `json:"contract"`
	RefBlockBytes string        `json:"ref_block_bytes"`
	RefBlockHash  string        `json:"ref_block_hash"`
	Expiration    int64         `json:"expiration"`
	FeeLimit      int64         `json:"fee_limit"`
	Timestamp     int64         `json:"timestamp"`
}

type DataEle struct {
	Ret                  []RetEle    `json:"ret"`
	Signature            []string    `json:"signature"`
	Txid                 string      `json:"txId"`
	NetUsage             int64       `json:"net_usage"`
	RawDataHex           string      `json:"raw_data_hex"`
	NetFee               int64       `json:"net_fee"`
	EnergyUsage          int64       `json:"energy_usage"`
	BlockNumber          int64       `json:"blockNumber"`
	BlockTimestamp       int64       `json:"block_timestamp"`
	EnergyFee            int64       `json:"energy_fee"`
	EnergyUsageTotal     int64       `json:"energy_usage_total"`
	RawData              RawData     `json:"raw_data"`
	InternalTransactions interface{} `json:"internal_transactions"`
}

type Meta struct {
	At       int64 `json:"at"`
	PageSize int64 `json:"page_size"`
}

type UserTransactionHistory struct {
	Data    []DataEle `json:"data"`
	Success bool      `json:"success"`
	Meta    Meta      `json:"meta"`
}
