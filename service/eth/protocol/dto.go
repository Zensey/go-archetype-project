package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"unsafe"
)

const (
	endpoint = "https://ethereum-rpc.publicnode.com/"
)

// Transaction - transaction object
type Transaction struct {
	Hash             string
	Nonce            int
	BlockHash        string
	BlockNumber      *int
	TransactionIndex *int
	From             string
	To               string
	Value            big.Int
	Gas              int
	GasPrice         big.Int
	Input            string
}

func (t *Transaction) Parse(data []byte) error {
	proxy := new(proxyTransaction)
	if err := json.Unmarshal(data, proxy); err != nil {
		return err
	}

	*t = *(*Transaction)(unsafe.Pointer(proxy))
	return nil
}

type proxyTransaction struct {
	Hash             string  `json:"hash"`
	Nonce            hexInt  `json:"nonce"`
	BlockHash        string  `json:"blockHash"`
	BlockNumber      *hexInt `json:"blockNumber"`
	TransactionIndex *hexInt `json:"transactionIndex"`
	From             string  `json:"from"`
	To               string  `json:"to"`
	Value            hexBig  `json:"value"`
	Gas              hexInt  `json:"gas"`
	GasPrice         hexBig  `json:"gasPrice"`
	Input            string  `json:"input"`
}

type EthRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// EthError - ethereum error
type EthError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EthResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *EthError       `json:"error"`
}

type BlockWithoutTransactions struct {
	Number           hexInt   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       hexBig   `json:"difficulty"`
	TotalDifficulty  hexBig   `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             hexInt   `json:"size"`
	GasLimit         hexInt   `json:"gasLimit"`
	GasUsed          hexInt   `json:"gasUsed"`
	Timestamp        hexInt   `json:"timestamp"`
	Uncles           []string `json:"uncles"`
	// Transactions     []proxyTransaction `json:"transactions"`
}

type BlockWithTransactions struct {
	Number           hexInt             `json:"number"`
	Hash             string             `json:"hash"`
	ParentHash       string             `json:"parentHash"`
	Nonce            string             `json:"nonce"`
	Sha3Uncles       string             `json:"sha3Uncles"`
	LogsBloom        string             `json:"logsBloom"`
	TransactionsRoot string             `json:"transactionsRoot"`
	StateRoot        string             `json:"stateRoot"`
	Miner            string             `json:"miner"`
	Difficulty       hexBig             `json:"difficulty"`
	TotalDifficulty  hexBig             `json:"totalDifficulty"`
	ExtraData        string             `json:"extraData"`
	Size             hexInt             `json:"size"`
	GasLimit         hexInt             `json:"gasLimit"`
	GasUsed          hexInt             `json:"gasUsed"`
	Timestamp        hexInt             `json:"timestamp"`
	Uncles           []string           `json:"uncles"`
	Transactions     []proxyTransaction `json:"transactions"`
}

func GetLastFinalizedBlockID() (int, error) {
	log.Printf("GetLastFinalizedBlockID!")

	req := EthRequest{
		ID:      0,
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			"finalized",
			false,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		return 0, err
	}
	resp, err := http.Post(endpoint, "application/json", &buf)
	if err != nil {
		return 0, err
	}

	ethResp := EthResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ethResp)
	if err != nil {
		return 0, err
	}
	if ethResp.Error != nil {
		return 0, fmt.Errorf("ethereum error: Code=%d Message=%s", ethResp.Error.Code, ethResp.Error.Message)
	}

	block := BlockWithoutTransactions{}
	err = json.Unmarshal(ethResp.Result, &block)
	if err != nil {
		return 0, err
	}

	return int(block.Number), err
}

func GetBlock(blockID int) (BlockWithTransactions, error) {
	log.Printf("GetBlock! 0x%x", blockID)
	block := BlockWithTransactions{}

	req := EthRequest{
		ID:      0,
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			IntToHex(blockID),
			true,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		return block, err
	}

	res, err := http.Post(endpoint, "application/json", &buf)
	if err != nil {
		return block, err
	}

	ethResp := EthResponse{}
	err = json.NewDecoder(res.Body).Decode(&ethResp)
	if err != nil {
		return block, err
	}
	if ethResp.Error != nil {
		return block, fmt.Errorf("ethereum error: Code=%d Message=%s", ethResp.Error.Code, ethResp.Error.Message)
	}

	err = json.Unmarshal(ethResp.Result, &block)
	return block, err
}

// Eth1 returns 1 ethereum value (10^18 wei)
func Eth1() *big.Int {
	return big.NewInt(1000000000000000000)
}
