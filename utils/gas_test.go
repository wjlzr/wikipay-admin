package utils

import (
	"fmt"
	"testing"
)

func TestGas(t *testing.T) {
	ethGas := GetGasPrice("BTC")
	fmt.Println(ethGas)
}
