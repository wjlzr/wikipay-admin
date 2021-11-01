package financial

import (
	"fmt"
	"wikipay-admin/common"
	orm "wikipay-admin/database"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/tools"
)

//
type TimeReq struct {
	StartTime int64 `form:"startTime" json:"startTime"`
	EndTime   int64 `form:"endTime" json:"endTime"`
}

//
type BidAsk struct {
	Coin string  `json:"coin"`
	Usd  float64 `json:"usd"`
	//	RealUsd float64 `json:"real_usd"`
	Type int     `json:"type"`
	Cny  float64 `json:"cny"` //换算成人民币
}

//获取点差
func FindBidAsk(req *TimeReq) ([]BidAsk, error) {
	var bidasks []BidAsk
	if req.StartTime <= 0 {
		req.StartTime = tools.MilliSecond()
	}
	if req.EndTime <= 0 {
		req.EndTime = tools.MilliSecond()
	}

	sql := fmt.Sprintf(`
			SELECT any_value(coin) AS coin, TRUNCATE(IFNULL(any_VALUE(SUM(amount * ABS(real_usd - usd))),0),8) AS usd,type FROM withdraw_deposit 
			WHERE account_type = 2 AND type = 1 AND status = 4 AND close_at BETWEEN %d AND %d
			GROUP BY coin 
			UNION ALL 
			SELECT any_value(coin) AS coin, TRUNCATE(IFNULL(any_VALUE(SUM(amount * ABS(real_usd - usd))),0),8) AS usd,type FROM withdraw_deposit 
			WHERE account_type = 2 AND type = 2 AND status = 4 AND close_at BETWEEN %d AND %d
			GROUP BY coin
			UNION ALL
			SELECT IF(any_value(c.coin)<>'','USDT','USDT') AS coin, TRUNCATE(IFNULL(SUM(c.usd) ,0),4) AS usd, any_value(c.type) FROM (
				SELECT IF(bid_coin<>'','USDT','USDT') AS coin, IFNULL(SUM(bid_amount * ABS(real_usd - usd)),0) AS usd,if(bid_coin<>'','5','5') AS type
				FROM coin_bid_ask
				WHERE bid_coin='USDT' AND create_at BETWEEN %d AND %d
				UNION ALL
				SELECT IF(ask_coin<>'','USDT','USDT') AS coin, IFNULL(SUM(ask_amount * ABS(real_usd - usd)),0) AS usd,if(ask_coin<>'','5','5') AS type
				FROM coin_bid_ask
				WHERE ask_coin='USDT' AND create_at BETWEEN %d AND %d
			) c
			UNION ALL
			SELECT IF(any_value(b.coin)<>'','BTC','BTC') AS coin, TRUNCATE(IFNULL(SUM(b.usd),0),8) AS usd, any_value(b.type) FROM (
				SELECT IF(bid_coin<>'','BTC','BTC') AS coin,IFNULL(SUM(bid_amount * ABS(real_usd - usd)),0) AS usd,if(bid_coin<>'','5','5') AS type
				FROM coin_bid_ask
				WHERE bid_coin='BTC' AND create_at BETWEEN %d AND %d
				UNION ALL
				SELECT IF(ask_coin<>'','BTC','BTC') AS coin,IFNULL(SUM(ask_amount * ABS(real_usd - usd)),0) AS usd,if(ask_coin<>'','5','5') AS type
				FROM coin_bid_ask
				WHERE ask_coin='BTC' AND create_at BETWEEN %d AND %d
			) b
			UNION ALL
			SELECT IF(any_value(a.coin)<>'','ETH','ETH') AS coin, TRUNCATE(IFNULL(SUM(a.usd),0),5) AS usd, any_value(a.type) FROM (
				SELECT IF(bid_coin<>'','ETH','ETH') AS coin,IFNULL(SUM(bid_amount * ABS(real_usd - usd)),0) AS usd,if(bid_coin<>'','5','5') AS type
				FROM coin_bid_ask
				WHERE bid_coin='ETH' AND create_at BETWEEN %d AND %d
				UNION ALL
				SELECT IF(ask_coin<>'','ETH','ETH') AS coin,IFNULL(SUM(ask_amount * ABS(real_usd - usd)),0) AS usd,if(ask_coin<>'','5','5') AS type
				FROM coin_bid_ask
				WHERE ask_coin='ETH' AND create_at BETWEEN %d AND %d ) a
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
		req.EndTime,
	)

	err := orm.Eloquent.Raw(sql).Find(&bidasks).Error
	if err != nil {
		return nil, err
	}
	var (
		infos   []BidAsk
		deposit = BidAsk{
			Coin: common.USDT,
			Usd:  0,
			Type: 1,
		}
		withdraw = BidAsk{
			Coin: common.USDT,
			Usd:  0,
			Type: 2,
		}
	)
	for _, b := range bidasks {
		if b.Coin == common.USDT_ERC20 || b.Coin == common.USDT_OMNI {
			switch b.Type {
			case 1:
				deposit.Usd = xfloat64.Add(deposit.Usd, b.Usd)
			case 2:
				withdraw.Usd = xfloat64.Add(withdraw.Usd, b.Usd)
			}
		} else {
			infos = append(infos, b)
		}
	}
	infos = append(infos, deposit)
	infos = append(infos, withdraw)

	//换算人民币
	cnyPrice := common.GetCurrencyPrice(common.CURRENCY_CNY)
	for k, v := range infos {
		usdPrice := common.GetCoinUsdPrice(v.Coin)
		infos[k].Cny = xfloat64.Float64Truncate(xfloat64.Mul(v.Usd, xfloat64.Mul(usdPrice, cnyPrice)), 6)
	}
	return infos, nil
}
