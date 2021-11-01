package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	goDateTime = "2006-01-02 15:04:05"
	goDate     = "2006-01-02"
	goTime     = "15:04:05"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

//生成订单号
func MustGenerateOrderId(orderType int) int64 {
	orderId, err := GenerateOrderId(orderType)
	if err != nil {
		panic(err)
	}
	return orderId
}

//获取订单号
//订单号规则
//年月日(8位) + 毫秒取余(3位) + 7位随机数
func GenerateOrderId(orderType int) (int64, error) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	strNow := time.Now().Format(goDateTime)
	formatDate, err := time.Parse(goDateTime, strNow)
	if err != nil {
		return 0, err
	}

	var (
		//currTime = strings.Replace(formatDate.Format(goTime), ":", "", -1)
		//id = subLastString(userId, 3)
		currDate = strings.Replace(formatDate.Format(goDate), "-", "", -1)
		orderId  string
	)

	r := randInt(10000000, 99999999)
	orderId = fmt.Sprintf("%s0%d%d", currDate, orderType, r)
	return strconv.ParseInt(orderId, 10, 64)
}

//随机数-int
func randInt(min, max int) int {
	return min + r.Intn(max-min)
}

func MilliSecond() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}
