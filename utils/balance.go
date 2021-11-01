package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"wikipay-admin/models/xfloat64"
)

type balanceInfo struct {
	Address            string `json:"address"`
	TotalReceived      int    `json:"total_received"`
	TotalSent          int    `json:"total_sent"`
	Balance            int    `json:"balance"`
	UnconfirmedBalance int    `json:"unconfirmed_balance"`
	FinalBalance       int    `json:"final_balance"`
	NTx                int    `json:"n_tx"`
	UnconfirmedNTx     int    `json:"unconfirmed_n_tx"`
	FinalNTx           int    `json:"final_n_tx"`
}

//获取btc当前数量
func GetBtcBalance(address string) float64 {
	if address != "" {
		url := fmt.Sprintf("https://api.blockcypher.com/v1/btc/main/addrs/%s/balance", address)
		resp, err := http.Get(url)
		if err == nil {
			bs, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				var b balanceInfo
				err = json.Unmarshal(bs, &b)
				if err == nil {
					return xfloat64.Div(float64(b.Balance), math.Pow(10, 8))
				}
			}
		}
	}
	return 0
}
