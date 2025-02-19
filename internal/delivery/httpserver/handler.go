package httpserver

import (
	"encoding/json"
	"eth_parser/internal/domain/entity"
	"eth_parser/internal/domain/parser"
	"net/http"
	"strings"
)

type TransactionHandler struct {
	Parser parser.Parser
}

func NewTransactionHandler(parser parser.Parser) *TransactionHandler {
	return &TransactionHandler{
		Parser: parser,
	}
}

func (h *TransactionHandler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentBlock := h.Parser.GetCurrentBlock()

	json.NewEncoder(w).Encode(map[string]int{"current_block": currentBlock})
}

func (h *TransactionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	address, ok := requestBody["address"]
	if !ok || address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	subscribed := h.Parser.Subscribe(address)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": subscribed, "address": address})
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := strings.TrimPrefix(r.URL.Path, "/get-transaction/")
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	transactions := h.Parser.GetTransactions(address)
	if transactions == nil {
		json.NewEncoder(w).Encode([]entity.Transaction{})
		return
	}
	json.NewEncoder(w).Encode(transactions)
}
