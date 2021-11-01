package financial

import (
	"fmt"
	"wikipay-admin/common"
	orm "wikipay-admin/database"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/tools"
)

type Fee struct {
	Coin     string  `json:"coin"`     //币咱
	Charge   float64 `json:"charge"`   //收取
	Spending float64 `json:"spending"` //支出
	Total    float64 `json:"total"`    //总计
	Cny      float64 `json:"cny"`      //换算成人民币
}

//查找费率
func FindFee(req *TimeReq) ([]Fee, error) {
	var fees []Fee
	if req.StartTime <= 0 {
		req.StartTime = tools.MilliSecond()
	}
	if req.EndTime <= 0 {
		req.EndTime = tools.MilliSecond()
	}
	sql := fmt.Sprintf(`
			SELECT IFNULL(coin,'BTC') AS coin, IFNULL(SUM(fee),0) AS charge, IFNULL(SUM(ABS(tx_fee)),0) AS spending,TRUNCATE(IFNULL(SUM(fee - ABS(tx_fee)),0),8) AS total
			FROM withdraw_deposit
			WHERE coin = 'BTC' AND account_type = 1 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d
			UNION ALL 
			SELECT IFNULL(coin,'ETH') AS coin, IFNULL(SUM(fee),0) AS charge, IFNULL(SUM(tx_fee),0) AS spending,TRUNCATE(IFNULL(SUM(fee - tx_fee),0),5) AS total
			FROM withdraw_deposit
			WHERE coin = 'ETH' AND account_type = 1 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d
			UNION ALL
			SELECT IF(any_value(a.coin)<>'','USDT','USDT') AS coin, TRUNCATE(IFNULL(SUM(a.charge),0),4) AS charge,TRUNCATE(IFNULL(SUM(a.spending),0),4) AS spending, TRUNCATE(IFNULL(SUM(a.total),0),4) AS total FROM (
				SELECT IF(coin<>'','USDT','USDT') AS coin, IFNULL(SUM(fee * any_value(usd)),0) AS charge, IFNULL(SUM(tx_fee * tx_fee_usd),0) AS spending,IFNULL(SUM(fee * usd -  tx_fee_usd * tx_fee),0) AS total
				FROM withdraw_deposit
				WHERE coin = 'USDT-ERC20' AND account_type = 1 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d
				UNION ALL
				SELECT IF(coin<>'','USDT','USDT') AS coin, IFNULL(SUM(fee * any_value(usd)),0) AS charge, IFNULL(SUM(tx_fee * tx_fee_usd),0) AS spending,IFNULL(SUM(fee * any_value(usd) -  tx_fee_usd * tx_fee),0) AS total
				FROM withdraw_deposit
				WHERE coin = 'USDT-OMNI' AND account_type = 1 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d
			) a
			UNION ALL
			SELECT IF(any_value(b.coin)<>'','USD','USD') AS coin, TRUNCATE(IFNULL(SUM(b.charge),0),2) AS charge,TRUNCATE(IFNULL(SUM(b.spending),0),2) AS spending, TRUNCATE(IFNULL(SUM(b.total),0),2) AS total FROM (
				SELECT IF(coin<>'','USD','USD') AS coin, IFNULL(SUM(fee * any_value(usd)),0) AS charge, IFNULL(SUM(ABS(tx_fee) * any_value(usd)),0) AS spending,IFNULL(SUM(fee - ABS(tx_fee)) * any_value(usd),0) AS total
				FROM withdraw_deposit
				WHERE coin = 'BTC' AND account_type = 2 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d 
				UNION ALL
				SELECT IF(coin<>'','USD','USD') AS coin, IFNULL(SUM(fee * any_value(usd)),0) AS charge, IFNULL(SUM(tx_fee * any_value(usd)),0) AS spending,IFNULL(SUM(fee + tx_fee) * any_value(usd),0) AS total
				FROM withdraw_deposit
				WHERE coin = 'ETH' AND account_type = 2 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d 
				UNION ALL 
				SELECT IF(coin<>'','USD','USD') AS coin, IFNULL(SUM(fee * usd),0) AS charge, IFNULL(SUM(tx_fee * tx_fee_usd),0) AS spending,IFNULL(SUM(fee * usd - tx_fee_usd * tx_fee),0) AS total
				FROM withdraw_deposit
				WHERE coin = 'USDT-ERC20' AND account_type = 2 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d
				UNION ALL
				SELECT IF(coin<>'','USD','USD') AS coin, IFNULL(SUM(fee * usd),0) AS charge, IFNULL(SUM(tx_fee * tx_fee_usd),0) AS spending,IFNULL(SUM(fee * usd - tx_fee_usd * tx_fee),0) AS total
				FROM withdraw_deposit
				WHERE coin = 'USDT-OMNI' AND account_type = 2 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d
			) b
	`,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime,
		req.StartTime,
		req.EndTime)

	err := orm.Eloquent.Raw(sql).Find(&fees).Error
	if err != nil {
		return nil, err
	}

	//换算人民币
	cnyPrice := common.GetCurrencyPrice(common.CURRENCY_CNY)
	for k, v := range fees {
		usdPrice := common.GetCoinUsdPrice(v.Coin)
		fees[k].Cny = xfloat64.Float64Truncate(xfloat64.Mul(v.Total, xfloat64.Mul(usdPrice, cnyPrice)), 6)
	}
	return fees, nil
}
