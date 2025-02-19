package parser

import "eth_parser/internal/domain/entity"

type Parser interface {
	// GetCurrentBlock last parsed block
	GetCurrentBlock() int
	// Subscribe add address to observer
	Subscribe(address string) bool
	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) []entity.Transaction
}
