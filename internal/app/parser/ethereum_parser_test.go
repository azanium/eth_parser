package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type mockHTTPClient struct {
	responses map[string][]byte
	err       error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var rpcReq rpcRequest
	if err := json.Unmarshal(body, &rpcReq); err != nil {
		return nil, err
	}

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(m.responses[rpcReq.Method])),
	}

	return response, nil
}

type mockSubscriptionRepo struct {
	subscriptions map[string]bool
}

func (m *mockSubscriptionRepo) StoreSubscription(address string) error {
	m.subscriptions[address] = true
	return nil
}

func (m *mockSubscriptionRepo) IsSubscribed(address string) bool {
	return m.subscriptions[address]
}

func TestGetCurrentBlock(t *testing.T) {
	tests := []struct {
		name          string
		chainIDResp   []byte
		blockNumResp  []byte
		expectedBlock int
		expectedError bool
		httpClientErr error
	}{
		{
			name:          "successful block retrieval",
			chainIDResp:   []byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`),
			blockNumResp:  []byte(`{"jsonrpc":"2.0","id":1,"result":"0x100"}`),
			expectedBlock: 256,
		},
		{
			name:          "error in chain ID response",
			chainIDResp:   []byte(`{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"error"}}`),
			expectedBlock: 0,
			expectedError: true,
		},
		{
			name:          "http client error",
			httpClientErr: io.ErrUnexpectedEOF,
			expectedBlock: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				responses: map[string][]byte{
					methodChainID:  tt.chainIDResp,
					methodBlockNum: tt.blockNumResp,
				},
				err: tt.httpClientErr,
			}

			mockRepo := &mockSubscriptionRepo{subscriptions: make(map[string]bool)}
			parser := NewEthereumParser(mockClient, mockRepo)

			block := parser.GetCurrentBlock()
			if block != tt.expectedBlock {
				t.Errorf("expected block %d, got %d", tt.expectedBlock, block)
			}
		})
	}
}

func TestSubscribe(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		preSubscribed  bool
		expectedResult bool
	}{
		{
			name:           "new subscription",
			address:        "0x123",
			preSubscribed:  false,
			expectedResult: true,
		},
		{
			name:           "already subscribed",
			address:        "0x456",
			preSubscribed:  true,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSubscriptionRepo{subscriptions: make(map[string]bool)}
			if tt.preSubscribed {
				mockRepo.StoreSubscription(tt.address)
			}

			parser := NewEthereumParser(nil, mockRepo)
			result := parser.Subscribe(tt.address)

			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}

			if !mockRepo.IsSubscribed(tt.address) {
				t.Error("address should be subscribed")
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	tests := []struct {
		name          string
		address       string
		subscribed    bool
		chainIDResp   []byte
		logsResp      []byte
		txResp        []byte
		expectedTxs   int
		expectedError bool
	}{
		{
			name:        "successful transaction retrieval",
			address:     "0x123",
			subscribed:  true,
			chainIDResp: []byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`),
			logsResp:    []byte(`{"jsonrpc":"2.0","id":1,"result":[{"address":"0x123","blockHash":"0x1","blockNumber":"0x1","transactionHash":"0x1"}]}`),
			txResp:      []byte(`{"jsonrpc":"2.0","id":1,"result":{"hash":"0x1","from":"0x123","to":"0x456"}}`),
			expectedTxs: 1,
		},
		{
			name:        "not subscribed",
			address:     "0x789",
			subscribed:  false,
			expectedTxs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				responses: map[string][]byte{
					methodChainID:  tt.chainIDResp,
					methodLogs:     tt.logsResp,
					methodTxByHash: tt.txResp,
				},
			}

			mockRepo := &mockSubscriptionRepo{subscriptions: make(map[string]bool)}
			if tt.subscribed {
				mockRepo.StoreSubscription(tt.address)
			}

			parser := NewEthereumParser(mockClient, mockRepo)
			txs := parser.GetTransactions(tt.address)

			if len(txs) != tt.expectedTxs {
				t.Errorf("expected %d transactions, got %d", tt.expectedTxs, len(txs))
			}
		})
	}
}

func TestSendRPCRequest(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		params        []any
		response      []byte
		expectedError bool
		httpClientErr error
	}{
		{
			name:     "successful request",
			method:   methodBlockNum,
			response: []byte(`{"jsonrpc":"2.0","id":1,"result":"0x100"}`),
		},
		{
			name:          "http client error",
			method:        methodBlockNum,
			httpClientErr: io.ErrUnexpectedEOF,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				responses: map[string][]byte{
					tt.method: tt.response,
				},
				err: tt.httpClientErr,
			}

			parser := NewEthereumParser(mockClient, nil)
			_, err := parser.sendRPCRequest(context.Background(), tt.method, 1, tt.params)

			if (err != nil) != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
