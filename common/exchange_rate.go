package common

import (
	"strings"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/redis"
)

const (
	currencyPrefix = "currency_usd_"
	toUsdPrice     = "_to_usd_price"
)

//获取数字货币的美元价格
func GetCoinUsdPrice(coin string) float64 {
	coin = strings.ToLower(coin)
	val := redis.ClusterClient().Get(coin + toUsdPrice).Val()
	if val != "" {
		return xfloat64.StrToFloat64(val)
	}
	return 0
}

//获取法定货币价格
func GetCurrencyPrice(currency string) float64 {
	val := redis.ClusterClient().Get(currencyPrefix + currency).Val()
	if val != "" {
		return xfloat64.StrToFloat64(val)
	}
	return 0
}
