package model

type ScanMeta struct {
	At          int64             `json:"at"`
	FingerPrint string            `json:"fingerprint"`
	Links       map[string]string `json:"links"`
	PageSize    int64             `json:"page_size"`
}

type ScanRawData struct {
	Number         int64  `json:"number"`
	TxTrieRoot     string `json:"txTrieRoot"`
	WitnessAddress string `json:"witness_address"`
	ParentHash     string `json:"parentHash"`
	Version        int64  `json:"version"`
	Timestamp      int64  `json:"timestamp"`
}

type BlockHeader struct {
	RawData          ScanRawData `json:"raw_data"`
	WitnessSignature string      `json:"witness_signature"`
}

type Transactions struct {
	Ret        []RetEle `json:"ret"`
	Signature  []string `json:"signature"`
	Txid       string   `json:"txId"`
	RawDataHex string   `json:"raw_data_hex"`
	RawData    RawData  `json:"raw_data"`
}

type BlockData struct {
	BlockId      string         `json:"blockId"`
	BlockHeader  BlockHeader    `json:"block_header"`
	Transactions []Transactions `json:"transactions"`
}

type History struct {
	Id     int64
	TxHash string
}
