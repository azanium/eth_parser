package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func HexToInt(hexStr string) (int64, error) {
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	return strconv.ParseInt(hexStr, 16, 64)
}

func AddressToHex(address string) string {
	address = strings.ToLower(strings.TrimPrefix(address, "0x"))
	return "0x" + fmt.Sprintf("%064s", address)
}
