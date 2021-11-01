package financial

import (
	"fmt"
	orm "wikipay-admin/database"
)

//交易请求
type TransactionReq struct {
	Pagination
	Coin      string `form:"coin" binding:"required"`
	Type      int    `form:"type" binding:"required"`
	StartTime int64  `form:"startTime" binding:"required"`
	EndTime   int64  `form:"endTime" binding:"required"`
}

//基本信息
type TransactionBaseInfo struct {
	Type     int    `json:"type"`     //类型
	CreateAt int64  `json:"createAt"` //创建时间
	OrderId  string `json:"orderId"`  //
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	NickName string `json:"nickName"`
}

//
type TransactionCoinInfo struct {
	Coin  string `json:"coin"`  //币种
	Money string `json:"money"` //金额
}

//CARD充值
type TransactionCardInfo struct {
	TransactionBaseInfo
	TransactionCoinInfo
	CardSecret string `json:"cardSecret"` //卡密码
}

//充值
type TransactionDepositInfo struct {
	TransactionBaseInfo
	TransactionCoinInfo
	Address    string `json:"address"` //充值地址
	AccoutType string `json:"accountType"`
}

//提现
type TransactionWithdrawInfo struct {
	TransactionBaseInfo
	TransactionCoinInfo
	Fee         string `json:"fee"`     //手续费
	Address     string `json:"address"` //提现地址
	AccountType int    `json:"accountType"`
}

//兑换
type TransactionExchangeInfo struct {
	TransactionBaseInfo
	FromAccount string `json:"fromAccount"` //付款账户
	ToAccount   string `json:"toAccount"`   //收款账户
	Rate        string `json:"rate"`        //兑换汇率
	RealUsd     string `json:"realUsd"`     //实时汇率
}

//获取交易账单
func GetTransactionInfo(req *TransactionReq) (interface{}, int, error) {
	var (
		sql string
		err error
	)

	switch req.Type {
	case 1:
		var (
			depositInfos, dCountInfo []TransactionDepositInfo
			cardInfos, cCountInfo    []TransactionCardInfo
		)
		if req.Coin == "AB" {
			sql = fmt.Sprintf(` 
				SELECT a.type, a.create_at, a.order_number AS order_id, CONCAT(b.first_name,b.last_name) AS name, b.phone, b.nick_name, a.coin, a.amount AS money, CONCAT(LEFT(c.card_no,4), '****' ,RIGHT(c.card_no,4)) AS card_secret FROM withdraw_deposit a 
				LEFT JOIN user b ON a.user_id = b.id
				LEFT JOIN topup_card_list c ON a.order_number = c.order_number
				WHERE a.type = 1 AND a.account_type = 1 AND a.status = 4 AND coin = '%s' AND a.create_at BETWEEN %d AND %d
				ORDER BY a.create_at DESC
			   `,
				req.Coin,
				req.StartTime,
				req.EndTime)
			err = orm.Eloquent.Raw(sql).Find(&dCountInfo).Error

			err = orm.Eloquent.Raw(fmt.Sprintf(`%s LIMIT %d OFFSET %d `,
				sql,
				req.PageSize,
				req.PageNum,
			)).Find(&cardInfos).Error

			return cardInfos, len(dCountInfo), err
		}

		sql = fmt.Sprintf(`
			SELECT c.* FROM (
				SELECT a.type,a.close_at,a.account_type, a.create_at, a.order_number AS order_id, CONCAT(b.first_name,b.last_name) AS name, CONCAT(LEFT(b.phone,3), '****' ,RIGHT(b.phone,3)) AS phone, b.nick_name, a.coin, a.amount AS money ,a.to_address AS address 
				FROM withdraw_deposit a 
				LEFT JOIN user b ON a.user_id = b.id
				WHERE a.type = 1 AND a.account_type = 1 AND a.status = 4 AND coin = '%s' AND a.close_at BETWEEN %d AND %d
				UNION ALL
				SELECT a.type,a.close_at,a.account_type, a.create_at, a.order_number AS order_id, CONCAT(b.first_name,b.last_name) AS name, CONCAT(LEFT(b.phone,3), '****' ,RIGHT(b.phone,3)) AS phone, b.nick_name, a.coin, a.amount AS money ,a.to_address AS address 
				FROM withdraw_deposit a 
				LEFT JOIN user b ON a.user_id = b.id
				WHERE a.type = 1 AND a.account_type = 2 AND a.status = 4 AND coin = '%s' AND a.close_at BETWEEN %d AND %d
			) c ORDER BY c.close_at DESC
	  		 `,
			req.Coin,
			req.StartTime,
			req.EndTime,
			req.Coin,
			req.StartTime,
			req.EndTime)
		err = orm.Eloquent.Raw(sql).Find(&cCountInfo).Error
		err = orm.Eloquent.Raw(fmt.Sprintf(`%s LIMIT %d OFFSET %d `,
			sql,
			req.PageSize,
			req.PageNum,
		)).Find(&depositInfos).Error

		return depositInfos, len(cCountInfo), err
	case 2:
		sql = fmt.Sprintf(` 
			SELECT c.* FROM (
				SELECT a.type, a.create_at, a.account_type,a.order_number AS order_id, CONCAT(b.first_name,b.last_name) AS name, b.phone, b.nick_name, a.coin, a.amount AS money, a.fee,a.close_at,a.to_address AS address 
				FROM withdraw_deposit a 
				LEFT JOIN user b ON a.user_id = b.id
				WHERE a.type = 2 AND a.account_type = 1 AND a.status = 4 AND coin = '%s' AND a.close_at BETWEEN %d AND %d
				UNION ALL
				SELECT a.type, a.create_at,a.account_type,a.order_number AS order_id, CONCAT(b.first_name,b.last_name) AS name, b.phone,b.nick_name, a.coin, a.amount AS money, a.fee,a.close_at,a.to_address AS address 
				FROM withdraw_deposit a 
				LEFT JOIN user b ON a.user_id = b.id
				WHERE a.type = 2 AND a.account_type = 2 AND a.status = 4 AND coin = '%s' AND a.close_at BETWEEN %d AND %d
			) c ORDER BY c.close_at DESC
	   `,
			req.Coin,
			req.StartTime,
			req.EndTime,
			req.Coin,
			req.StartTime,
			req.EndTime)

		var infos, countInfo []TransactionWithdrawInfo
		err = orm.Eloquent.Raw(sql).Find(&countInfo).Error
		err = orm.Eloquent.Raw(fmt.Sprintf(`%s LIMIT %d OFFSET %d `,
			sql,
			req.PageSize,
			req.PageNum,
		)).Find(&infos).Error

		return infos, len(countInfo), err

	case 5:
		sql = fmt.Sprintf(`
		   SELECT c.* FROM (
				SELECT if(1=1,5,5) AS type, 
					a.create_at, 
					a.order_number AS order_id, 
					CONCAT(b.first_name,b.last_name) AS name, 
					b.phone,
					b.nick_name,
					CONCAT('-',0+CAST(a.bid_amount AS char(32)),a.bid_coin) AS from_account,
					CONCAT('+',0+CAST(a.ask_amount AS char(32)),a.ask_coin) AS to_account,
					IF(a.bid_coin='USD',1/a.usd,a.usd) AS rate,
					a.real_usd
				FROM coin_bid_ask a
				LEFT JOIN user b ON a.user_id = b.id
				WHERE a.status = 1 AND a.bid_coin = '%s' AND a.create_at BETWEEN %d AND %d
				UNION ALL
				SELECT if(1=1,5,5) AS type, 
					a.create_at, 
					a.order_number AS order_id, 
					CONCAT(b.first_name,b.last_name) AS name, 
					b.phone, 
					b.nick_name,
					CONCAT('-',0+CAST(a.bid_amount AS char(32)),a.bid_coin) AS from_account,
					CONCAT('+',0+CAST(a.ask_amount AS char(32)),a.ask_coin) AS to_account,
					IF(a.ask_coin='USD',1/a.usd,a.usd) AS rate,
					a.real_usd
				FROM coin_bid_ask a
				LEFT JOIN user b ON a.user_id = b.id
				WHERE a.status = 1 AND a.ask_coin = '%s' AND a.create_at BETWEEN %d AND %d
			) c  ORDER BY c.create_at DESC
		`,
			req.Coin,
			req.StartTime,
			req.EndTime,
			req.Coin,
			req.StartTime,
			req.EndTime,
		)
		var exchangeInfo, countInfo []TransactionExchangeInfo

		err = orm.Eloquent.Raw(sql).Find(&countInfo).Error
		err = orm.Eloquent.Raw(fmt.Sprintf(`%s LIMIT %d OFFSET %d `,
			sql,
			req.PageSize,
			req.PageNum,
		)).Find(&exchangeInfo).Error

		return exchangeInfo, len(countInfo), err
	}
	return nil, 0, nil
}
