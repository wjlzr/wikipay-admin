package wallet

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
type WithdrawDeposit struct {
	OrderNumber     string  `json:"orderNumber" gorm:"type:bigint(18) unsigned zerofill;primary_key"` // 订单号
	TransFlowNumber string  `json:"transFlowNumber" gorm:"type:bigint(18);"`                          // 交易流水号
	UserId          int64   `json:"userId" gorm:"type:bigint(10);"`                                   //
	Type            int     `json:"type" gorm:"type:tinyint(1);"`                                     // 交易类型
	Coin            string  `json:"coin" gorm:"type:varchar(100);"`                                   // 币种
	Amount          string  `json:"amount" gorm:"type:decimal(32,8);"`                                // 提现数量
	Usd             string  `json:"usd" gorm:"type:decimal(32,8);"`                                   // 数字货币汇率
	FromAddress     string  `json:"fromAddress" gorm:"type:varchar(100);"`                            // 打款地址
	ToAddress       string  `json:"toAddress" gorm:"type:varchar(100);"`                              // 提现地址
	TxHash          string  `json:"txHash" gorm:"type:varchar(120);"`                                 // 交易哈希
	TxFee           string  `json:"txFee" gorm:"type:decimal(32,8);"`                                 // 链上交易手续费
	Status          int     `json:"status" gorm:"type:tinyint(1);"`                                   // 充提状态
	Fee             string  `json:"fee" gorm:"type:decimal(32,8);"`                                   // 平台提现手续费
	Comment         string  `json:"comment" gorm:"type:varchar(200);"`                                // 备注内容
	Content         string  `json:"content" gorm:"type:varchar(300);"`                                // 其它内容
	CreateAt        int64   `json:"createAt" gorm:"type:bigint(13);"`                                 // 交易时间
	ProcessAt       int64   `json:"processAt" gorm:"type:bigint(13);"`                                // 处理时间
	ReviewAt        int64   `json:"reviewAt" gorm:"type:bigint(13) unsigned zerofill;"`               // 审核时间
	CloseAt         int64   `json:"closeAt" gorm:"type:bigint(13) unsigned zerofill;"`                // 结束时间
	ActivityAmount  string  `json:"activityAmount" gorm:"type:decimal(32,8);"`
	AccountType     int     `json:"accountType" gorm:"type:tinyint(1);"` //
	RealUsd         string  `json:"realUsd" gorm:"type:decimal(16,4);"`  //
	Spread          float64 `json:"spread" gorm:"-"`                     //点差
	DataScope       string  `json:"dataScope" gorm:"-"`
	Params          string  `json:"params"  gorm:"-"`
}

//筛选参数
type SearchParams struct {
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}

//
func (WithdrawDeposit) TableName() string {
	return "withdraw_deposit"
}

//创建
func (e *WithdrawDeposit) Create() (WithdrawDeposit, error) {
	var doc WithdrawDeposit
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *WithdrawDeposit) Get() (WithdrawDeposit, error) {
	var doc WithdrawDeposit
	table := orm.Eloquent.Table(e.TableName())

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取WithdrawDeposit带分页
func (e *WithdrawDeposit) GetPage(pageSize, pageIndex int, info, status, tradeType string, req SearchParams) ([]WithdrawDeposit, int, error) {
	var doc []WithdrawDeposit

	table := orm.Eloquent.Select("*").Table(e.TableName())

	//数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var where string = " 1 = 1 "
	if info != "" {
		where = fmt.Sprintf(" %s AND concat(user_id,order_number) like '%s%s%s'", where, "%", info, "%")
	}
	if status != "" {
		where = fmt.Sprintf(" %s AND status = %s ", where, status)
	}
	if tradeType != "" {
		where = fmt.Sprintf(" %s AND type = %s ", where, tradeType)
	}

	if req.StartTime != "" {
		where = fmt.Sprintf(" %s AND ((create_at >= %s AND create_at <= %s) OR (close_at >= %s AND close_at <= %s))", where, req.StartTime, req.EndTime, req.StartTime, req.EndTime)
	}
	if err := table.
		Select("order_number,trans_flow_number,user_id,type,coin,amount,usd,from_address,to_address,tx_hash,tx_fee,status,fee,comment,content,create_at,process_at,review_at,close_at,activity_amount,account_type,real_usd,(CASE account_type WHEN '2' THEN ABS(amount*usd-amount*real_usd) ELSE 0 END) spread").
		Where(where).
		Order("create_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Where(where).Count(&count)

	//点差换算人民币
	cnyPrice := common.GetCurrencyPrice(common.CURRENCY_CNY)
	for k, v := range doc {
		if v.AccountType == 2 {
			doc[k].Spread = xfloat64.Mul(v.Spread, cnyPrice)
		}
	}

	return doc, count, nil
}

//更新状态
func (e *WithdrawDeposit) UpdateStatus(id int64) (update WithdrawDeposit, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", id).First(&update).Error; err != nil {
		return
	}
	if err = orm.Eloquent.Table(e.TableName()).
		Model(&update).
		Where("order_number = ?", id).
		Update("status", 3).Error; err != nil {
	}
	return
}

//更新WithdrawDeposit
func (e *WithdrawDeposit) Update(id int64) (update WithdrawDeposit, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Where("order_number = ?", id).Updates(&e).Error; err != nil {
		return
	}

	return
}

// 删除WithdrawDeposit
func (e *WithdrawDeposit) Delete(id int64) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", id).Delete(&WithdrawDeposit{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *WithdrawDeposit) BatchDelete(id []int64) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number in (?)", id).Delete(&WithdrawDeposit{}).Error; err != nil {
		return
	}
	Result = true
	return
}
