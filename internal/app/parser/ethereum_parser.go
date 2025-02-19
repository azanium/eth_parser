package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"eth_parser/internal/domain/entity"
	httpclient "eth_parser/internal/domain/http_client"
	"eth_parser/internal/domain/parser"
	"eth_parser/internal/domain/repository"
	"eth_parser/internal/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const (
	rpcURL         = "https://ethereum-rpc.publicnode.com/" //"https://cloudflare-eth.com"
	rpcVersion     = "2.0"
	erc20Transfer  = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	methodBlockNum = "eth_blockNumber"
	methodChainID  = "eth_chainId"
	methodLogs     = "eth_getLogs"
	methodTxByHash = "eth_getTransactionByHash"
)

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int64  `json:"id"`
}

type rpcResponse struct {
	ID      int64           `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type logEntry struct {
	Address         string `json:"address"`
	BlockHash       string `json:"blockHash"`
	BlockNumber     string `json:"blockNumber"`
	TransactionHash string `json:"transactionHash"`
}

type logsResponse struct {
	Logs []logEntry `json:"result"`
}

type EthereumParser struct {
	mutex  sync.RWMutex
	client httpclient.HTTPClient
	repo   repository.SubscriptionRepo
}

var _ parser.Parser = (*EthereumParser)(nil)

func NewEthereumParser(client httpclient.HTTPClient, repo repository.SubscriptionRepo) *EthereumParser {
	return &EthereumParser{
		client: client,
		repo:   repo,
	}
}

func (ep *EthereumParser) GetCurrentBlock() int {
	ep.mutex.Lock()
	defer ep.mutex.Unlock()

	ctx := context.Background()

	chainID, err := ep.getChainID(ctx)
	if err != nil {
		log.Println(fmt.Errorf("failed to get chain ID: %w", err))
		return 0
	}

	blockNumResp, err := ep.sendRPCRequest(ctx, methodBlockNum, chainID, nil)
	if err != nil {
		log.Println(fmt.Errorf("failed to get block number: %w", err))
		return 0
	}

	var response rpcResponse
	if err = json.Unmarshal(blockNumResp, &response); err != nil {
		log.Println(fmt.Errorf("failed to unmarshal response: %w", err))
		return 0
	}

	if response.Error != nil {
		log.Println(fmt.Errorf("RPC request failed: %s", response.Error.Message))
		return 0
	}

	var hex string
	if err = json.Unmarshal(response.Result, &hex); err != nil {
		log.Println(fmt.Errorf("failed to unmarshal response: %w", err))
		return 0
	}

	blockNum, err := utils.HexToInt(hex)
	if err != nil {
		log.Println(fmt.Errorf("failed to parse block number: %w", err))
		return 0
	}
	return int(blockNum)
}

func (ep *EthereumParser) Subscribe(address string) bool {
	if !ep.repo.IsSubscribed(address) {
		ep.repo.StoreSubscription(address)
		return true
	}

	return true
}

func (ep *EthereumParser) GetTransactions(address string) []entity.Transaction {
	if !ep.repo.IsSubscribed(address) {
		log.Println(fmt.Errorf("address %s is not subscribed", address))
		return []entity.Transaction{}
	}

	ep.mutex.Lock()
	defer ep.mutex.Unlock()

	ctx := context.Background()

	chainID, err := ep.getChainID(ctx)
	if err != nil {
		log.Println(fmt.Errorf("failed to get chain ID: %w", err))
		return []entity.Transaction{}
	}

	params := []any{
		map[string]any{
			"fromBlock": "0x0",
			"toBlock":   "latest",
			"topics": []any{
				erc20Transfer, []string{utils.AddressToHex(address), utils.AddressToHex(address)},
			},
		},
	}

	logs, err := ep.fetchLogs(ctx, chainID, params)

	var transactions []entity.Transaction
	for _, logEntry := range logs {
		tx, err := ep.getTransaction(ctx, chainID, logEntry.TransactionHash)
		if err != nil {
			log.Println(fmt.Errorf("failed to get transaction by hash: %w", err))
			continue
		}
		transactions = append(transactions, *tx)
	}

	return transactions
}

func (ep *EthereumParser) getTransaction(ctx context.Context, chainID int64, hash string) (*entity.Transaction, error) {
	transactionRaw, err := ep.sendRPCRequest(ctx, methodTxByHash, chainID, []any{hash})
	if err != nil {
		return nil, err
	}

	var tx entity.Transaction
	if err := json.Unmarshal([]byte(transactionRaw), &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	return &tx, nil
}

func (ep *EthereumParser) fetchLogs(ctx context.Context, chainID int64, params []any) ([]logEntry, error) {
	logsRaw, err := ep.sendRPCRequest(ctx, methodLogs, chainID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	var logs logsResponse
	if err = json.Unmarshal(logsRaw, &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs: %w", err)
	}
	return logs.Logs, nil
}

func (ep *EthereumParser) getChainID(ctx context.Context) (int64, error) {
	chainIDResp, err := ep.sendRPCRequest(ctx, methodChainID, 1, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get chain ID: %w", err)
	}
	var response rpcResponse
	if err = json.Unmarshal(chainIDResp, &response); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return 0, fmt.Errorf("RPC request failed: %s", response.Error.Message)
	}

	if len(chainIDResp) == 0 || chainIDResp == nil {
		return 0, fmt.Errorf("empty response")
	}

	var hex string
	if err = json.Unmarshal(response.Result, &hex); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	id, err := utils.HexToInt(hex)
	if err != nil {
		return 0, fmt.Errorf("failed to parse chain ID: %w", err)
	}
	return id, nil
}

func (ep *EthereumParser) sendRPCRequest(ctx context.Context, method string, id int64, params []any) ([]byte, error) {
	rpcReq := rpcRequest{JSONRPC: rpcVersion, Method: method, Params: params, ID: id}
	body, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := ep.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}
