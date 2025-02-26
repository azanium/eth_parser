# Ethereum Transaction Parser

A Go-based service for parsing and tracking Ethereum transactions .

## Overview

This service provides a robust interface for monitoring and tracking Ethereum blockchain transactions. 

## Getting Started

### Prerequisites

- Go 1.19 or later
- Make (for using Makefile commands)

### Installation

1. Clone the repository

```bash
git clone [repository-url]
cd eth_parser
```

2. Install dependencies

```bash

go mod download
```

3. Build the project

```bash
make build
```

4. Run the service

```bash
make run
```

## API Documentation

### Get Current Block

```
GET /get-current-block
```

```bash
curl -X GET localhost:8080/get-current-block
```

Returns the current block number being monitored.

Response:

```json
{
    "current_block": 123456
}
```

### Subscribe to Address

```
POST /subscribe
```

```bash
curl -X POST localhost:8080/subscribe -i -d '{"address": "ADDRESS"}'
```

Subscribe to monitor transactions for a specific Ethereum address.

Request Body:

```json
{
    "address": "ADDRESS", "status": true
}
```

Response:

```json
{
    "status": true,
    "address": "ADDRESS"
}
```

### Get Transactions

```
GET /get-transaction/{ethereum_address}
```

```bash
curl -X GET localhost:8080/get-transaction/ADDRESS
```

Retrieve all transactions for a subscribed address.

Response:

```json
[
    {
        "hash": "0x...",
        "from": "0x...",
        "to": "0x...",
        "value": "0x...",
        "blockNumber": "0x..."
    }
]
```

## Error Handling

The service implements comprehensive error handling for:

- Invalid transaction hashes
- Network connectivity issues
- Invalid Ethereum addresses
- Concurrent access management

## Architecture

The project follows a clean architecture pattern with the following components:

- `internal/app/parser`: Core transaction parsing logic
- `internal/delivery/httpserver`: HTTP API implementation
- `internal/domain`: Business logic interfaces and entities
- `internal/utils`: Utility functions

## Development

### Running Tests

```bash
make test
```

### Code Style

The project follows standard Go code style guidelines. Run the following to ensure code quality:

```bash
make lint
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Ethereum JSON-RPC API
- Go standard library
