package wallet

import (
	"fmt"
	"testing"
)

func TestCoin(t *testing.T) {
	c := getCoin("btc")
	fmt.Println(c)
}
