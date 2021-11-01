package utils

import "strings"

var (
	coins        = []string{"btc", "eth", "omni", "erc20"}
	UpperCoins   = []string{"BTC", "ETH", "USDT-OMNI", "USDT-ERC20"}
	OperateCoins = []string{"BTC", "ETH", "USDT"}
)

//获取大写币种
func GetUpperCoin(coin string) string {
	coin = strings.ToUpper(coin)
	for _, v := range UpperCoins {
		if v == coin {
			return v
		}
	}
	return ""
}

//获取币种
func GetCoin(coin string) string {
	coin = strings.ToLower(coin)

	for _, v := range coins {
		if strings.Index(coin, v) >= 0 {
			return v
		}
	}
	return ""
}

//获取大小写币名称
func GetCoinName(coin string, isUpper bool) string {
	var (
		usdt = "USDT"
	)

	coin = strings.ToUpper(coin)
	if !isUpper {
		usdt = "usdt"
		coin = strings.ToLower(coin)
	}

	if strings.Index(coin, usdt) >= 0 {
		return usdt
	}
	return coin
}
