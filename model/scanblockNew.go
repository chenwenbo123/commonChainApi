package model

import "time"

type RealBlock struct {
	Ret        []map[string]string `json:"ret"`
	Signature  []string            `json:"signature"`
	Txid       string              `json:"txID"`
	RawData    AnotherRawData      `json:"raw_data"`
	RawDataHex string              `json:"raw_data_hex"`
}

type AnotherRawData struct {
	Contract      []AnotherContractEle `json:"contract"`
	RefBlockBytes string               `json:"ref_block_bytes"`
	RefBlockHash  string               `json:"ref_block_hash"`
	Expiration    int64                `json:"expiration"`
	FeeLimit      int64                `json:"fee_limit"`
	Timestamp     int64                `json:"timestamp"`
}

type AnotherContractEle struct {
	Parameter AnotherParameterEle `json:"parameter"`
	Type      string              `json:"type"`
}

type AnotherParameterEle struct {
	Value   ValueEle `json:"value"`
	TypeUrl string   `json:"type_url"`
}

type ValueEle struct {
	Data            string `json:"data"`
	OwnerAddress    string `json:"owner_address"`
	ContractAddress string `json:"contract_address"`
}

func FormateTime(msecond int64) string {
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Unix(int64(msecond/1000), 0).Format(timeLayout)
	//fmt.Println(timeStr)
	return timeStr
}

type TransactionInfo struct {
	Data    []Data                 `json:"data"`
	Meta    map[string]interface{} `json:"meta"`
	Success bool                   `json:"success"`
}

type Data struct {
	BlockNum            int64      `json:"block_number"`
	BlockTimestamp      int64      `json:"block_timestamp"`
	CallContractAddress string     `json:"call_contract_address"`
	ContractAddress     string     `json:"contract_address"`
	EventIndex          int64      `json:"event_index"`
	EventName           string     `json:"event_name"`
	Result              Result     `json:"result"`
	ResultType          ResultType `json:"result_type"`
	Event               string     `json:"event"`
	TransactionId       string     `json:"transaction_id"`
}

type Result struct {
	Zero  string `json:"0"`
	One   string `json:"1"`
	Two   string `json:"2"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type ResultType struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}
