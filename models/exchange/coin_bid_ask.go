package exchange

import (
	"fmt"
	_ "time"
	"wikipay-admin/common"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/tools"
)

//
type CoinBidAsk struct {
	OrderNumber string  `json:"orderNumber" gorm:"type:bigint(18) unsigned;primary_key"` //
	UserId      int     `json:"userId" gorm:"type:bigint(13);"`                          // 用户序列号
	BidCoin     string  `json:"bidCoin" gorm:"type:varchar(30);"`                        // 持有名称
	AskCoin     string  `json:"askCoin" gorm:"type:varchar(30);"`                        // 换算名称
	BidAmount   string  `json:"bidAmount" gorm:"type:decimal(32,8);"`                    // 持有数量
	AskAmount   string  `json:"askAmount" gorm:"type:decimal(32,8);"`                    //
	Usd         string  `json:"usd" gorm:"type:decimal(16,4);"`                          //
	Status      int     `json:"status" gorm:"type:tinyint(1);"`                          //
	CreateAt    int     `json:"createAt" gorm:"type:bigint(13);"`                        // 创建时间
	RealUsd     string  `json:"realUsd" gorm:"type:decimal(16,4);"`                      //
	Spread      float64 `json:"spread" gorm:"-"`
	DataScope   string  `json:"dataScope" gorm:"-"`
	Params      string  `json:"params"  gorm:"-"`
}

//参数
type SearchParams struct {
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}

//列名
func (CoinBidAsk) TableName() string {
	return "coin_bid_ask"
}

//创建
func (e *CoinBidAsk) Create() (CoinBidAsk, error) {
	var doc CoinBidAsk

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *CoinBidAsk) Get() (CoinBidAsk, error) {
	var doc CoinBidAsk

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取带分页
func (e *CoinBidAsk) GetPage(pageSize int, pageIndex int, info string, req SearchParams) ([]CoinBidAsk, int, error) {
	var doc []CoinBidAsk
	table := orm.Eloquent.Select("*").Table(e.TableName())

	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var where string
	if info != "" {
		where = fmt.Sprintf("concat(order_number,user_id) like '%s%s%s'", "%", info, "%")
	}

	if req.StartTime != "" {
		table = table.Where("create_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		table = table.Where("create_at <= ?", req.EndTime)
	}

	var count int
	if err := table.
		Select("order_number,user_id,bid_coin,ask_coin,bid_amount,ask_amount,usd,status,create_at,real_usd,(CASE bid_coin WHEN 'USD' THEN ABS(ask_amount*usd-ask_amount*real_usd) ELSE ABS(bid_amount*usd-bid_amount*real_usd) END) spread").
		Where(where).
		Offset((pageIndex - 1) * pageSize).
		Order("create_at DESC").
		Limit(pageSize).
		Find(&doc).
		Error; err != nil {
		return nil, 0, err
	}
	table.Where(where).Count(&count)

	//点差换算人民币
	cnyPrice := common.GetCurrencyPrice(common.CURRENCY_CNY)
	for k, v := range doc {
		doc[k].Spread = xfloat64.Mul(v.Spread, cnyPrice)
	}

	return doc, count, nil
}

//更新
func (e *CoinBidAsk) Update(id string) (update CoinBidAsk, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", id).First(&update).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

//删除
func (e *CoinBidAsk) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", id).Delete(&CoinBidAsk{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *CoinBidAsk) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number in (?)", id).Delete(&CoinBidAsk{}).Error; err != nil {
		return
	}
	Result = true
	return
}
