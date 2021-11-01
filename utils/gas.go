package utils

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"wikipay-admin/common"
)

const (
	ethAccessKey = "058e9ba70d221d068ffe9a2329ee6d066b4f9e52bedb0a5749b570bdb8d0"
)

var (
	gasUrl = map[string]string{
		common.BTC:        "https://bitcoinfees.earn.com/api/v1/fees/recommended",
		common.USDT_OMNI:  "https://bitcoinfees.earn.com/api/v1/fees/recommended",
		common.ETH:        "https://data-api.defipulse.com/api/v1/egs/api/ethgasAPI.json?api-key=" + ethAccessKey,
		common.USDT_ERC20: "https://data-api.defipulse.com/api/v1/egs/api/ethgasAPI.json?api-key=" + ethAccessKey,
	}
)

//
type EthGasPrice struct {
	Fast          int                `json:"fast"`
	Fastest       int                `json:"fastest"`
	SafeLow       int                `json:"safeLow"`
	Average       int                `json:"average"`
	BlockTime     float64            `json:"block_time"`
	BlockNum      int                `json:"blockNum"`
	Speed         float64            `json:"speed"`
	SafeLowWait   float64            `json:"safeLowWait"`
	AvgWait       float64            `json:"avgWait"`
	FastWait      float64            `json:"fastWait"`
	FastestWait   float64            `json:"fastestWait"`
	GasPriceRange map[string]float64 `json:"gasPriceRange"`
}

//
type BtcGasPrice struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
}

//获取eth链上手续率
func GetGasPrice(coin string) float64 {
	coin = strings.ToUpper(coin)
	resp, err := http.Get(gasUrl[coin])
	if err == nil {
		bs, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			switch coin {
			case common.ETH, common.USDT_ERC20:
				var gasPrice EthGasPrice
				if err = json.Unmarshal(bs, &gasPrice); err == nil {
					return float64(gasPrice.Fastest) / math.Pow(10, 10) * 22000
				}
			case common.BTC, common.USDT_OMNI:
				var gasPrice BtcGasPrice
				if err = json.Unmarshal(bs, &gasPrice); err == nil {
					price := (2*148 + 34*1 + 10) * gasPrice.FastestFee
					return float64(price) / math.Pow(10, 8)
				}
			}
		}
	}
	return 0
}
