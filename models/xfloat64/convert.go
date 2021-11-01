package xfloat64

import (
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

//浮点数截取指定位数
func Float64Truncate(f float64, bit int) float64 {
	ff, err := strconv.ParseFloat(TruncateFloat64(f, bit), 64)
	if err != nil {
		return 0
	}
	return ff
}

//浮点数截取
func TruncateFloat64(f float64, m int) string {
	str := decimal.NewFromFloat(f).String()
	newn := strings.Split(str, ".")
	if len(newn) < 2 || m >= len(newn[1]) {
		return str
	}
	if m == 0 {
		return newn[0]
	}
	return newn[0] + "." + newn[1][:m]
}

///
func TruncateStringFloat64(str string, m int) string {
	newn := strings.Split(str, ".")
	if len(newn) < 2 || m >= len(newn[1]) {
		return str
	}
	return newn[0] + "." + newn[1][:m]
}
