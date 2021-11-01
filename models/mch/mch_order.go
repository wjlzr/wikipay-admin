package mch

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type MchOrder struct {
	TradeNo      string `json:"tradeNo" gorm:"type:varchar(32);primary_key"` // 预付订单号
	MchId        string `json:"mchId" gorm:"type:varchar(32);"`              // 商户号
	Subject      string `json:"subject" gorm:"type:varchar(128);"`           // 标题
	Body         string `json:"body" gorm:"type:varchar(256);"`              // 商品描述
	OutTradeNo   string `json:"outTradeNo" gorm:"type:varchar(32);"`         // 商户订单号
	CurrencyType string `json:"currencyType" gorm:"type:varchar(16);"`       // 货币类型
	TotalAmount  string `json:"totalAmount" gorm:"type:varchar(9);"`         // 总金额
	NotifyUrl    string `json:"notifyUrl" gorm:"type:varchar(256);"`         // 通知地址
	Attach       string `json:"attach" gorm:"type:varchar(127);"`            // 附加数据
	TradeType    string `json:"tradeType" gorm:"type:varchar(16);"`          // 支付类型
	StartAt      string `json:"startAt" gorm:"type:varchar(14);"`            // 开始时间
	ExpireAt     string `json:"expireAt" gorm:"type:varchar(14);"`           // 结束时间
	CreateIp     string `json:"createIp" gorm:"type:varchar(64);"`           // 终端IP
	CreateAt     int    `json:"createAt" gorm:"type:bigint(13);"`            // 创建时间
	DataScope    string `json:"dataScope" gorm:"-"`
	Params       string `json:"params"  gorm:"-"`
}

//
func (MchOrder) TableName() string {
	return "mch_order"
}

// 创建MchOrder
func (e *MchOrder) Create() (MchOrder, error) {
	var doc MchOrder
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取MchOrder
func (e *MchOrder) Get() (MchOrder, error) {
	var doc MchOrder
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取MchOrder带分页
func (e *MchOrder) GetPage(pageSize, pageIndex int, info string) ([]MchOrder, int, error) {
	var doc []MchOrder
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
		where = fmt.Sprintf("concat(trade_no,out_trade_no,mch_id) like '%s%s%s'", "%", info, "%")
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

// 更新MchOrder
func (e *MchOrder) Update(tradeNo string) (update MchOrder, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("trade_no = ?", tradeNo).First(&update).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除MchOrder
func (e *MchOrder) Delete(tradeNo string) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("trade_no = ?", tradeNo).Delete(&MchOrder{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *MchOrder) BatchDelete(tradeNos []string) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("trade_no in (?)", tradeNos).Delete(&MchOrder{}).Error; err != nil {
		return
	}
	Result = true
	return
}
