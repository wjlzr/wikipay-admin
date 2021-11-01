package mch

import (
	"fmt"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"

	_ "time"
)

type MchTransaction struct {
	TransactionNo  string `json:"transactionNo" gorm:"type:varchar(32);primary_key"` // 交易订单号
	OutTradeNo     string `json:"outTradeNo" gorm:"type:varchar(32);"`               // 商家订单号
	CurrencyType   string `json:"currencyType" gorm:"type:varchar(16);"`             // 货币类型
	TotalAmount    string `json:"totalAmount" gorm:"type:varchar(9);"`               // 金额
	UserId         string `json:"userId" gorm:"type:bigint(13);"`                    // 付款人序列号
	TradeState     string `json:"tradeState" gorm:"type:varchar(32);"`               // 交易状态
	Attach         string `json:"attach" gorm:"type:varchar(127);"`                  //
	TradeType      string `json:"tradeType" gorm:"type:varchar(16);"`                // 交易类型
	TradeStateDesc string `json:"tradeStateDesc" gorm:"type:varchar(256);"`          // 交易状态描述
	CreateAt       string `json:"createAt" gorm:"type:bigint(13);"`                  // 创建时间
	NotifyUrl      string `json:"notifyUrl" gorm:"type:varchar(256);"`               //
	DataScope      string `json:"dataScope" gorm:"-"`
	Params         string `json:"params"  gorm:"-"`
}

//
func (MchTransaction) TableName() string {
	return "mch_transaction"
}

//创建
func (e *MchTransaction) Create() (MchTransaction, error) {
	var doc MchTransaction

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *MchTransaction) Get() (MchTransaction, error) {
	var doc MchTransaction

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取带分页
func (e *MchTransaction) GetPage(pageSize, pageIndex int, info string) ([]MchTransaction, int, error) {
	var doc []MchTransaction

	table := orm.Eloquent.Select("*").Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var where string
	if info != "" {
		where = fmt.Sprintf("concat(transaction_no,out_trade_no,user_id) like '%s%s%s'", "%", info, "%")
	}

	if err := table.Where(where).
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Count(&count)
	return doc, count, nil
}

//更新
func (e *MchTransaction) Update(id string) (update MchTransaction, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_no = ?", id).First(&update).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除
func (e *MchTransaction) Delete(id int64) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_no = ?", id).Delete(&MchTransaction{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *MchTransaction) BatchDelete(id []int64) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_no in (?)", id).Delete(&MchTransaction{}).Error; err != nil {
		return
	}
	Result = true
	return
}

//获取资产明细
func (e MchTransaction) GetAccountDetail(currencyType, tradeState string, userId int) (res string, err error) {
	var mchTransaction MchTransaction
	sql := fmt.Sprintf(`
				SELECT sum(total_amount) total_amount FROM mch_transaction WHERE user_id = %d AND currency_type = '%s' AND trade_state = '%s'
			`, userId, currencyType, tradeState)

	if err = orm.Eloquent.Raw(sql).Find(&mchTransaction).Error; err != nil {
		return mchTransaction.TotalAmount, err
	}

	return mchTransaction.TotalAmount, nil
}
