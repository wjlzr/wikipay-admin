package monitor

import (
	"fmt"
	orm "wikipay-admin/database"
)

//
type WithDrawInfo struct {
	Coin   string  `json:"coin"`
	Amount float64 `json:"amount"`
	Usd    float64 `json:"usd"`
	Type   int     `json:"type"` //1、公司 2、个人
}

//获取当日充值、提现统计
func FindWithDrawWithNow() ([]WithDrawInfo, error) {
	sql := `SELECT ANY_VALUE(coin) AS coin,IFNULL(ANY_VALUE(SUM(amount)),0) AS amount,ANY_VALUE(type) AS type FROM withdraw_deposit
			WHERE type = 1 AND status = 4 AND FROM_UNIXTIME(close_at/1000,'%s') = DATE_FORMAT(now(),'%s')
			GROUP BY coin
			UNION ALL
			SELECT ANY_VALUE(coin) AS coin,IFNULL(ANY_VALUE(SUM(amount + fee)),0) AS amount,ANY_VALUE(type) AS type  FROM withdraw_deposit
			WHERE type = 2 AND status = 4 AND FROM_UNIXTIME(close_at/1000,'%s') = DATE_FORMAT(now(),'%s')
			GROUP BY coin`

	var infos []WithDrawInfo
	err := orm.Eloquent.Raw(fmt.Sprintf(sql,
		"%Y-%m-%d",
		"%Y-%m-%d",
		"%Y-%m-%d",
		"%Y-%m-%d")).
		Find(&infos).
		Error

	return infos, err
}

//获取资产分类
func FindAssetGroup() ([]WithDrawInfo, error) {
	sql := `
			SELECT ANY_VALUE(coin) AS coin,IFNULL(ANY_VALUE(SUM(amount)),0) AS amount,ANY_VALUE(IF(coin<>'',1,1)) AS type 
			FROM monitor_address
			WHERE deleted_at IS NULL
			GROUP BY coin
			UNION ALL
			SELECT ANY_VALUE(coin) AS coin,IFNULL(ANY_VALUE(SUM(amount )),0) AS amount,ANY_VALUE(IF(coin<>'',2,2)) AS type  
			FROM monitor_user_address
			WHERE deleted_at IS NULL
			GROUP BY coin
		`

	var infos []WithDrawInfo
	err := orm.Eloquent.Raw(sql).
		Find(&infos).
		Error

	return infos, err
}
