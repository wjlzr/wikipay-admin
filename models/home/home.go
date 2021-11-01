package home

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"wikipay-admin/common"
	orm "wikipay-admin/database"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/tools"
)

var (
	ErrDateTime = errors.New("时间错误")
	//coins       = []string{"USDT", "BTC", "ETH", "USD"}
	// bits        = map[string]int{
	// 	"USDT": 4,
	// 	"BTC":  8,
	// 	"USD":  2,
	// 	"ETH":  5,
	// }
)

//用户统计信息
type TotalInfo struct {
	TotalUserNum     int64 `json:"totalUserNum"`     //总用户数
	TodayUserNum     int64 `json:"todayUserNum"`     //今日增加用户
	TodayDepositNum  int64 `json:"todayDepositNum"`  //今日充值笔数
	TodayWithdrawNum int64 `json:"todayWithdrawNum"` //今日提现笔数
	TodayExchangeNum int64 `json:"todayExchangeNum"` //今日兑换笔数
}

//
type TotailProfitInfo struct {
	TotalProfit          string `json:"totalProfit"`          //总收益
	CurrMonthTotalProfit string `json:"currMonthTotalProfit"` //当月收益
}

//
type MonthProfitInfo struct {
	Date     string `json:"date"`
	Profit   string `json:"profit"`
	Interest string `json:"interest"`
	Money    string `json:"money"`
}

//
type ProfitInfo struct {
	TotailProfitInfo TotailProfitInfo  `json:"totailProfitInfo"`
	MonthProfitInfo  []MonthProfitInfo `json:"monthProfitInfo"`
}

//
type TradeInfo struct {
	Amount string `json:"amount"`
	Coin   string `json:"coin"`
	Date   string `json:"date"`
	Type   int    `json:"type"`
}

//首页财务管理
type FinancialInfo struct {
	Account              string `json:"account"`
	TotalAmount          string `json:"totalAmount"`        //账户总额
	DepositAmount        string `json:"depositAmount"`      //充值
	AvailableAmount      string `json:"availableAmount"`    //可用
	FrozenAmount         string `json:"frozenAmount"`       //提现中
	WithDrawingAmount    string `json:"withdrawAmount"`     //已提现
	ExchangeAmount       string `json:"exchangeAmount"`     //兑换
	TotalAmountCny       string `json:"totalAmountCny"`     //账户总额-人民币
	DepositAmountCny     string `json:"depositAmountCny"`   //充值-人民币
	AvailableAmountCny   string `json:"availableAmountCny"` //可用-人民币
	FrozenAmountCny      string `json:"frozenAmountCny"`    //提现中-人民币
	WithDrawingAmountCny string `json:"withdrawAmountCny"`  //已提现-人民币
	ExchangeAmountCny    string `json:"exchangeAmountCny"`  //兑换-人民币

}

//
type wdinfo struct {
	Type   int
	Amount string
	Coin   string
}

//
type accountInfo struct {
	Coin      string
	Available string
	Frozen    string
}

//兑换
type exchangeInfo struct {
	Amount string
	Coin   string
}

//
type MonthReq struct {
	Month string `form:"month"`
}

//
type TradeReq struct {
	Month string `form:"month"`
	Coin  string `form:"coin"`
}

//
type FinancialReq struct {
	StartTime int64 `form:"startTime" binding:"required"`
	EndTime   int64 `form:"endTime" binding:"required"`
}

//
//获取总用户数等信息
func GetTotalInfo() (*TotalInfo, error) {
	var info TotalInfo
	err := orm.Eloquent.Raw(fmt.Sprintf(
		`SELECT 
			(SELECT count(*) FROM user) total_user_num,
			(SELECT count(*) FROM user WHERE create_at BETWEEN UNIX_TIMESTAMP(CONCAT(CURDATE(),' 00:00:00')) * 1000 AND UNIX_TIMESTAMP(CONCAT(CURDATE(),' 23:59:59')) * 1000 ) AS today_user_num,
			(SELECT count(*) FROM withdraw_deposit WHERE status = 4 AND type = 1 AND close_at BETWEEN UNIX_TIMESTAMP(CONCAT(CURDATE(),' 00:00:00')) * 1000 AND UNIX_TIMESTAMP(CONCAT(CURDATE(),' 23:59:59')) * 1000) AS today_deposit_num,
			(SELECT count(*) FROM withdraw_deposit WHERE status = 4 AND type = 2 AND close_at BETWEEN UNIX_TIMESTAMP(CONCAT(CURDATE(),' 00:00:00')) * 1000 AND UNIX_TIMESTAMP(CONCAT(CURDATE(),' 23:59:59')) * 1000) AS today_withdraw_num,
			(SELECT count(*) FROM coin_bid_ask WHERE create_at BETWEEN UNIX_TIMESTAMP(CONCAT(CURDATE(),' 00:00:00')) * 1000 AND UNIX_TIMESTAMP(CONCAT(CURDATE(),' 23:59:59')) * 1000) AS today_exchange_num
		`)).
		Scan(&info).Error

	return &info, err
}

//获取收益及月分收益
func GetProfit(month string) *ProfitInfo {
	tInfo, err := getTotalProfit(month)
	if err != nil {
		return nil
	}

	mInfos, err := getProfitWithMonth(month)
	if err != nil {
		return nil
	}

	var info ProfitInfo
	info.MonthProfitInfo = mInfos
	info.TotailProfitInfo = tInfo

	return &info
}

//获取收益
func getTotalProfit(month string) (TotailProfitInfo, error) {
	var info TotailProfitInfo
	err := orm.Eloquent.Raw(fmt.Sprintf(`
		SELECT
			(SELECT SUM(profit) FROM profit) AS total_profit,
			(SELECT SUM(profit) FROM profit WHERE FROM_UNIXTIME(create_at/1000,'%s') = DATE_FORMAT(now(),'%s')) AS curr_month_total_profit
	`,
		"%Y-%m",
		month)).
		Scan(&info).
		Error

	return info, err
}

//获取选择月的利息
func getProfitWithMonth(month string) ([]MonthProfitInfo, error) {
	var infos []MonthProfitInfo
	err := orm.Eloquent.Raw(fmt.Sprintf(`
		SELECT FROM_UNIXTIME(create_at/1000,'%s') AS date,IFNULL(any_value(sum(amount)),0) AS money ,IFNULL(any_value(sum(profit)),0) AS profit, any_value(interest) AS interest 
		FROM profit 
		WHERE FROM_UNIXTIME(create_at/1000,'%s')= '%s'
		GROUP BY FROM_UNIXTIME(create_at/1000,'%s')`,
		"%Y-%m-%d",
		"%Y-%m",
		month,
		"%Y-%m-%d")).
		Find(&infos).
		Error

	return infos, err
}

//充值、提现分离
func GroupTradeInfo(coin, month string) ([]TradeInfo, []TradeInfo, error) {
	infos, err := GetTradeInfo(coin, month)
	if err != nil {
		return nil, nil, err
	}

	withdrawInfo := make([]TradeInfo, 0)
	depositInfo := make([]TradeInfo, 0)
	for _, info := range infos {
		switch info.Type {
		case 1:
			depositInfo = append(depositInfo, info)
		case 2:
			withdrawInfo = append(withdrawInfo, info)
		}
	}
	return depositInfo, withdrawInfo, nil
}

//获取交易信息
func GetTradeInfo(coin, month string) ([]TradeInfo, error) {
	var infos []TradeInfo
	err := orm.Eloquent.Raw(fmt.Sprintf(`
		SELECT 
			FROM_UNIXTIME(close_at/1000,'%s') AS date,
			IFNULL(any_value(sum(amount + fee)),0) AS amount,
			any_value(coin) AS coin,
			any_value(type) AS type
		FROM withdraw_deposit 
		WHERE FROM_UNIXTIME(close_at/1000,'%s') = '%s' AND coin != 'AB' AND coin = '%s' AND type = 1 AND status = 4
		GROUP BY coin,FROM_UNIXTIME(close_at/1000,'%s')
		UNION ALL
		SELECT 
			FROM_UNIXTIME(close_at/1000,'%s') AS date,
			IFNULL(any_value(sum(amount + fee)),0) AS amount,
			any_value(coin) AS coin,
			any_value(type) AS type
		FROM withdraw_deposit 
		WHERE FROM_UNIXTIME(close_at/1000,'%s') = '%s'  AND coin = '%s' AND type = 2 AND status = 4
		GROUP BY coin,FROM_UNIXTIME(close_at/1000,'%s')`,
		"%Y-%m-%d",
		"%Y-%m",
		month,
		coin,
		"%Y-%m-%d",
		"%Y-%m-%d",
		"%Y-%m",
		month,
		coin,
		"%Y-%m-%d",
	)).Find(&infos).Error
	return infos, err
}

//财务管理
func GetFinancialInfo() ([]FinancialInfo, error) {
	var (
		wdInfos       []wdinfo
		accountInfos  []accountInfo
		exchangeInfos []exchangeInfo
		infos         []FinancialInfo
	)

	err := orm.Eloquent.Raw(fmt.Sprintf(`
		SELECT any_value(type) AS type,IFNULL(any_value(sum(amount + fee)),0) AS amount, any_value(coin) AS coin
		FROM  withdraw_deposit
		WHERE account_type = 1 AND status = 4 AND type = 2
		GROUP BY coin
		UNION ALL 
		SELECT any_value(type) AS type, IFNULL(any_value(sum(amount)),0) AS amount, any_value(coin) AS coin
		FROM  withdraw_deposit
		WHERE account_type = 1 AND status= 4 AND type = 1
		GROUP BY coin
		UNION ALL
		SELECT IFNULL(any_value(type),1) AS type,IFNULL(any_value(sum(IF(coin='AB',amount,amount * usd))),0) AS amount,any_value(IF(coin<>'','USD','USD')) AS coin
		FROM  withdraw_deposit
		WHERE account_type = 2 AND status = 4 AND type = 1
		UNION ALL
		SELECT IFNULL(any_value(type),2) AS type,IFNULL(any_value(sum(amount * usd + fee * usd)),0) AS amount,any_value(IF(coin<>'','USD','USD')) AS coin
		FROM  withdraw_deposit
		WHERE account_type = 2 AND status = 4 AND type = 2
	`)).Find(&wdInfos).Error

	err = orm.Eloquent.Raw(fmt.Sprintf(`
		SELECT any_value(coin) AS coin ,IFNULL(any_value(sum(available)),0) AS available, IFNULL(any_value(sum(frozen)),0) AS frozen FROM account
		GROUP BY coin`)).Find(&accountInfos).Error

	err = orm.Eloquent.Raw(fmt.Sprintf(`
			SELECT c.coin, (if(c.ask_amount=c.max_amount,c.ask_amount-c.min_amount,c.ask_amount-c.max_amount)) AS amount
			FROM
			(
				SELECT any_value(b.coin) AS coin,IFNULL(any_value(b.amount),0) AS ask_amount,IFNULL(any_value(max(b.amount)),0) AS max_amount,IFNULL(any_value(min(b.amount)),0) AS min_amount FROM (
					SELECT any_value(ask_coin) AS coin , IFNULL(any_value(sum(ask_amount)),0) AS amount
					FROM coin_bid_ask
					GROUP BY ask_coin
					UNION ALL
					SELECT any_value(bid_coin) AS coin ,IFNULL(any_value(sum(bid_amount)),0) AS amount
					FROM coin_bid_ask
					GROUP BY bid_coin
				) b GROUP BY b.coin 
			) c`)).Find(&exchangeInfos).Error

	coinMap := make(map[string]int, 0)
	for i := 0; i < len(common.Coins); i++ {
		coinMap[common.Coins[i]] = i
		info := FinancialInfo{
			Account: common.Coins[i],
		}
		infos = append(infos, info)
	}

	for _, v := range exchangeInfos {
		k := coinMap[v.Coin]
		infos[k].ExchangeAmount = v.Amount //ConvertFloat64ToStr(mapAsk[v.Coin] - ConvertStrToFloat64(v.Amount))
	}

	for _, v := range accountInfos {
		if v.Coin != common.USDT_ERC20 && v.Coin != common.USDT_OMNI {
			k := coinMap[v.Coin]
			if k >= 0 {
				infos[k].AvailableAmount = xfloat64.TruncateStringFloat64(v.Available, common.CoinBits[v.Coin])
				infos[k].FrozenAmount = xfloat64.TruncateStringFloat64(v.Frozen, common.CoinBits[v.Coin])
			}
		}
	}

	var (
		usdtDepositAmount, usdWithDrawingAmount float64
	)

	for _, v := range wdInfos {
		if v.Coin == common.BTC || v.Coin == common.ETH || v.Coin == common.USD {
			k := coinMap[v.Coin]
			if k >= 0 {
				switch v.Type {
				case 1:
					infos[k].DepositAmount = xfloat64.TruncateStringFloat64(v.Amount, common.CoinBits[v.Coin])
				case 2:
					infos[k].WithDrawingAmount = xfloat64.TruncateStringFloat64(v.Amount, common.CoinBits[v.Coin])
				}
			}
		}
		if strings.Index(v.Coin, common.USDT) >= 0 {
			switch v.Type {
			case 1:
				usdtDepositAmount = xfloat64.Float64Truncate(xfloat64.Add(usdtDepositAmount, ConvertStrToFloat64(v.Amount)), 4)
			case 2:
				usdWithDrawingAmount = xfloat64.Float64Truncate(xfloat64.Add(usdWithDrawingAmount, ConvertStrToFloat64(v.Amount)), 4)
			}
		}
	}
	infos[coinMap[common.USDT]].DepositAmount = xfloat64.TruncateStringFloat64(ConvertFloat64ToStr(usdtDepositAmount), 4)
	infos[coinMap[common.USDT]].WithDrawingAmount = xfloat64.TruncateStringFloat64(ConvertFloat64ToStr(usdWithDrawingAmount), 4)

	//infos[coinMap["USD"]].DepositAmount = wdInfos[coinMap["AB"]].Amount
	//	infos[coinMap["USD"]].TotalAmount = wdInfos[coinMap["AB"]].Amount
	for k, v := range infos {
		//f := ConvertStrToFloat64(v.DepositAmount) + ConvertStrToFloat64(v.ExchangeAmount) - ConvertStrToFloat64(v.AvailableAmount) - ConvertStrToFloat64(v.FrozenAmount) - ConvertStrToFloat64(v.WithDrawingAmount)
		//f := ConvertStrToFloat64(v.DepositAmount) + ConvertStrToFloat64(v.ExchangeAmount)
		//v.TotalAmount = ConvertFloat64ToStr(f)
		//	fmt.Println(v.TotalAmount)
		f := xfloat64.FromStringAdd(v.AvailableAmount, v.FrozenAmount)
		infos[k].TotalAmount = ConvertFloat64ToStr(f)
	}
	return infos, err
}

//根据时间区间获取财务数据
func GetFinancialInfoWithDateTime(req *FinancialReq) ([]FinancialInfo, error) {
	if req.StartTime <= 0 {
		req.StartTime = tools.MilliSecond()
	}

	if req.EndTime <= 0 {
		req.EndTime = tools.MilliSecond()
	}

	if req.StartTime > req.EndTime {
		return nil, ErrDateTime
	}

	var (
		wdInfos       []wdinfo
		accountInfos  []accountInfo
		exchangeInfos []exchangeInfo
		infos         []FinancialInfo
	)

	err := orm.Eloquent.Raw(fmt.Sprintf(`
		SELECT IFNULL(any_value(type),2) AS type,IFNULL(any_value(sum(amount + fee)),0) AS amount, any_value(coin) AS coin
		FROM  withdraw_deposit
		WHERE account_type = 1 AND status = 4 AND type = 2 AND close_at BETWEEN %d AND %d
		GROUP BY coin
		UNION ALL 
		SELECT IFNULL(any_value(type),1) AS type, IFNULL(any_value(sum(amount)),0) AS amount, any_value(coin) AS coin
		FROM  withdraw_deposit
		WHERE account_type = 1 AND status= 4 AND coin <> 'AB' AND type = 1 AND close_at BETWEEN %d AND %d
		GROUP BY coin
		UNION ALL
		SELECT IFNULL(any_value(type),1) AS type,IFNULL(any_value(sum(IF(coin='AB',amount,amount*usd))),0) AS amount, any_value(IF(coin<>'','USD','USD')) AS coin
		FROM  withdraw_deposit
		WHERE account_type = 2 AND status = 4 AND type = 1 AND close_at BETWEEN %d AND %d
		UNION ALL
		SELECT IFNULL(any_value(type),2) AS type,IFNULL(any_value(sum(amount * usd + fee * usd)),0) AS amount,any_value(IF(coin<>'','USD','USD')) AS coin
		FROM  withdraw_deposit
		WHERE account_type = 2 AND status = 4 AND type = 2 AND close_at BETWEEN %d AND %d
	`,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime)).Find(&wdInfos).Error

	err = orm.Eloquent.Raw(fmt.Sprintf(`
		SELECT any_value(coin) AS coin ,IFNULL(any_value(sum(available)),0) AS available, IFNULL(any_value(sum(frozen)),0) AS frozen FROM account
		GROUP BY coin`)).Find(&accountInfos).Error

	// err = orm.Eloquent.Raw(fmt.Sprintf(`
	// 	SELECT any_value(ask_coin) AS coin ,any_value(sum(ask_amount)) AS amount
	// 	FROM coin_bid_ask
	// 	WHERE create_at BETWEEN %d AND %d
	// 	GROUP BY ask_coin`,
	// 	req.StartTime,
	// 	req.EndTime)).Find(&exchangeInfos).Error

	err = orm.Eloquent.Raw(fmt.Sprintf(`
			SELECT c.coin, (if(c.ask_amount=c.max_amount,c.ask_amount-c.min_amount,c.ask_amount-c.max_amount)) AS amount
			FROM
			(
				SELECT any_value(b.coin) AS coin,IFNULL(any_value(b.amount),0) AS ask_amount,IFNULL(any_value(max(b.amount)),0) AS max_amount,IFNULL(any_value(min(b.amount)),0) AS min_amount FROM (
					SELECT any_value(ask_coin) AS coin , IFNULL(any_value(sum(ask_amount)),0) AS amount
					FROM coin_bid_ask
					WHERE create_at BETWEEN %d AND %d
					GROUP BY ask_coin
					UNION ALL
					SELECT any_value(bid_coin) AS coin ,IFNULL(any_value(sum(bid_amount)),0) AS amount
					FROM coin_bid_ask
					WHERE create_at BETWEEN %d AND %d
					GROUP BY bid_coin
				) b GROUP BY b.coin 
			) c`,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime)).Find(&exchangeInfos).Error

	coinMap := make(map[string]int, 0)
	for i := 0; i < len(common.Coins); i++ {
		coinMap[common.Coins[i]] = i
		info := FinancialInfo{
			Account: common.Coins[i],
		}
		infos = append(infos, info)
	}

	for _, v := range exchangeInfos {
		k := coinMap[v.Coin]
		infos[k].ExchangeAmount = xfloat64.TruncateStringFloat64(v.Amount, common.CoinBits[v.Coin])
	}

	for _, v := range accountInfos {
		if v.Coin != common.USDT_ERC20 && v.Coin != common.USDT_OMNI {
			k := coinMap[v.Coin]
			if k >= 0 {
				infos[k].AvailableAmount = xfloat64.TruncateStringFloat64(v.Available, common.CoinBits[v.Coin])
				infos[k].FrozenAmount = xfloat64.TruncateStringFloat64(v.Frozen, common.CoinBits[v.Coin])
			}
		}
	}

	var (
		usdtDepositAmount, usdWithDrawingAmount float64
	)

	for _, v := range wdInfos {
		if v.Coin == common.BTC || v.Coin == common.ETH || v.Coin == common.USD {
			k := coinMap[v.Coin]
			if k >= 0 {
				switch v.Type {
				case 1:
					infos[k].DepositAmount = xfloat64.TruncateStringFloat64(v.Amount, common.CoinBits[v.Coin])
				case 2:
					infos[k].WithDrawingAmount = xfloat64.TruncateStringFloat64(v.Amount, common.CoinBits[v.Coin])
				}
			}
		}
		if strings.Index(v.Coin, common.USDT) >= 0 {
			switch v.Type {
			case 1:
				usdtDepositAmount = xfloat64.Float64Truncate(xfloat64.Add(usdtDepositAmount, ConvertStrToFloat64(v.Amount)), 4)
			case 2:
				usdWithDrawingAmount = xfloat64.Float64Truncate(xfloat64.Add(usdWithDrawingAmount, ConvertStrToFloat64(v.Amount)), 4)
			}
		}
	}
	infos[coinMap[common.USDT]].DepositAmount = xfloat64.TruncateStringFloat64(ConvertFloat64ToStr(usdtDepositAmount), 4)
	infos[coinMap[common.USDT]].WithDrawingAmount = xfloat64.TruncateStringFloat64(ConvertFloat64ToStr(usdWithDrawingAmount), 4)

	cnyPrice := common.GetCurrencyPrice(common.CURRENCY_CNY)
	for k, v := range infos {

		f := xfloat64.FromStringAdd(v.AvailableAmount, v.FrozenAmount)
		infos[k].TotalAmount = ConvertFloat64ToStr(f)

		//返回各个字段对应的折算人民币
		usdPrice := common.GetCoinUsdPrice(v.Account)
		infos[k].TotalAmountCny = ConvertFloat64ToStr(xfloat64.Float64Truncate(xfloat64.Mul(ConvertStrToFloat64(v.TotalAmount), xfloat64.Mul(usdPrice, cnyPrice)), 6))
		infos[k].DepositAmountCny = ConvertFloat64ToStr(xfloat64.Float64Truncate(xfloat64.Mul(ConvertStrToFloat64(v.DepositAmount), xfloat64.Mul(usdPrice, cnyPrice)), 6))
		infos[k].AvailableAmountCny = ConvertFloat64ToStr(xfloat64.Float64Truncate(xfloat64.Mul(ConvertStrToFloat64(v.AvailableAmount), xfloat64.Mul(usdPrice, cnyPrice)), 6))
		infos[k].FrozenAmountCny = ConvertFloat64ToStr(xfloat64.Float64Truncate(xfloat64.Mul(ConvertStrToFloat64(v.FrozenAmount), xfloat64.Mul(usdPrice, cnyPrice)), 6))
		infos[k].WithDrawingAmountCny = ConvertFloat64ToStr(xfloat64.Float64Truncate(xfloat64.Mul(ConvertStrToFloat64(v.WithDrawingAmount), xfloat64.Mul(usdPrice, cnyPrice)), 6))
		infos[k].ExchangeAmountCny = ConvertFloat64ToStr(xfloat64.Float64Truncate(xfloat64.Mul(ConvertStrToFloat64(v.ExchangeAmount), xfloat64.Mul(usdPrice, cnyPrice)), 6))
	}
	return infos, err
}

//浮点数转字符串
func ConvertFloat64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

//字符串转浮点数
func ConvertStrToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

//struct to map
func structToMap(infos []exchangeInfo) map[string]float64 {
	mapInfo := make(map[string]float64, 0)
	for _, v := range infos {
		mapInfo[v.Coin] = ConvertStrToFloat64(v.Amount)
	}
	return mapInfo
}
