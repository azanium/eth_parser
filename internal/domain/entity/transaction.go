package entity

type Transaction struct {
	BlockHash        *string `json:"-"`
	BlockNumber      *string `json:"-"`
	TransactionIndex *string `json:"transactionIndex"`
	Hash             string  `json:"hash"`
	From             string  `json:"from"`
	To               *string `json:"to"`
	Gas              string  `json:"gas"`
	GasPrice         string  `json:"gasPrice"`
	Input            string  `json:"-"`
	Nonce            string  `json:"-"`
	Value            string  `json:"value"`
	V                string  `json:"-"`
	R                string  `json:"-"`
	S                string  `json:"-"`
}
