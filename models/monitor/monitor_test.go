package monitor

import (
	"fmt"
	"testing"
	"wikipay-admin/models/wallet/client"
)

func TestMonitor(t *testing.T) {
	price := GetCurrency(&client.BaseReq{
		Coin: "CNY",
	})
	fmt.Println(price)
}
