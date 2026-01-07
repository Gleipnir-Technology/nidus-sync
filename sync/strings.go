package sync

import (
	"strconv"
	"strings"
)

func HexToInt64(hexStr string) (int64, error) {
	// Remove "0x" prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")
	hexStr = strings.TrimPrefix(hexStr, "0X")

	// Parse hex string to uint64
	value, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func HexToUint64(hexStr string) (uint64, error) {
	// Remove "0x" prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")
	hexStr = strings.TrimPrefix(hexStr, "0X")

	// Parse hex string to uint64
	value, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
