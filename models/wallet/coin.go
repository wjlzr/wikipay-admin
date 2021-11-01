package wallet

import (
	_ "time"
	orm "wikipay-admin/database"
)

//
type Coin struct {
	Code                int    `json:"code" gorm:"type:tinyint(1) unsigned zerofill;primary_key"`        // 货币代码
	Type                int    `json:"type" gorm:"type:tinyint(1);"`                                     // 1.虚拟数字货币 2、visa 3、paypal
	Name                string `json:"name" gorm:"type:varchar(30);"`                                    // 货币名称
	Icon                string `json:"icon" gorm:"type:varchar(200);"`                                   // 图标
	IconHeight          int    `json:"iconHeight" gorm:"type:int(10);"`                                  //
	IconWidth           int    `json:"iconWidth" gorm:"type:int(10);"`                                   //
	ShowPrecision       int    `json:"showPrecision" gorm:"type:tinyint(1);"`                            // 显示精度
	WithdrawalPrecision int    `json:"withdrawalPrecision" gorm:"type:tinyint(1);"`                      // 提现精度
	PayPrecision        int    `json:"payPrecision" gorm:"type:tinyint(1);"`                             // 支付精度
	Status              int    `json:"status" gorm:"type:tinyint(1);"`                                   // 状态(0-不可用，1-可用)
	InfoUrl             string `json:"infoUrl" gorm:"type:varchar(100);"`                                // 信息链接
	WithdrawalMinNum    string `json:"withdrawalMinNum" gorm:"type:decimal(32,8) unsigned zerofill;"`    // 单次最小提现数量
	WithdrawalMaxNum    string `json:"withdrawalMaxNum" gorm:"type:decimal(32,8) unsigned zerofill;"`    // 单次最大提现数量
	WithdrawalMaxNumDay string `json:"withdrawalMaxNumDay" gorm:"type:decimal(32,8) unsigned zerofill;"` // 当日最大提现数量
	WithdrawalFee       string `json:"withdrawalFee" gorm:"type:decimal(32,8) unsigned zerofill;"`       // 单次提现费用
	DepositPriority     int    `json:"depositPriority" gorm:"type:tinyint(1);"`                          // 充值优先级
	WithdrawalPriority  int    `json:"withdrawalPriority" gorm:"type:tinyint(1);"`                       // 提现优先级
	CanDeposit          int    `json:"canDeposit" gorm:"type:tinyint(1);"`                               // 是否可充币(0-不可充，1-可充)
	CanWithdrawal       int    `json:"canWithdrawal" gorm:"type:tinyint(1);"`                            // 是否可提币(0-不可充，1-可充)
	IsShowTag           int    `json:"isShowTag" gorm:"type:tinyint(1);"`                                // 是否显示tag(0-不显示，1-显示)
	MinDepositNum       string `json:"minDepositNum" gorm:"type:decimal(32,8);"`                         // 单次最小充值金额
	MaxDepositNum       string `json:"maxDepositNum" gorm:"type:decimal(32,8);"`                         // 单次最大充值金额
	DepositConfirmNum   int    `json:"depositConfirmNum" gorm:"type:int(10);"`                           // 充币确认数
	Description         string `json:"description" gorm:"type:varchar(32);"`                             // 描述
	TxidPrefixUrl       string `json:"txidPrefixUrl" gorm:"type:varchar(1024);"`                         // 前缀
	BlockHeight         int    `json:"blockHeight" gorm:"type:int(11);"`                                 // 块高度
	ActivityFee         string `json:"activityFee" gorm:"type:decimal(10,4);"`                           //
	DataScope           string `json:"dataScope" gorm:"-"`
	Params              string `json:"params"  gorm:"-"`
}

func (Coin) TableName() string {
	return "coin"
}

// 创建Coin
func (e *Coin) Create() (Coin, error) {
	var doc Coin
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Coin
func (e *Coin) Get() (Coin, error) {
	var doc Coin
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "code = ?", e.Code).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Coin带分页
func (e *Coin) GetPage(pageSize int, pageIndex int) ([]Coin, int, error) {
	var doc []Coin
	table := orm.Eloquent.Select("code,type,name,icon,icon_height,icon_width,show_precision,withdrawal_precision,pay_precision,status,info_url,0+CAST(withdrawal_min_num AS char(32)) AS withdrawal_min_num,0+CAST(withdrawal_max_num AS char(32)) AS withdrawal_max_num,0+CAST(withdrawal_max_num_day AS char(32)) AS withdrawal_max_num_day ,0+CAST(withdrawal_fee AS char(32)) AS withdrawal_fee,deposit_priority,withdrawal_priority,can_deposit,can_withdrawal,is_show_tag,0+CAST(min_deposit_num AS char(32)) AS min_deposit_num,0+CAST(max_deposit_num AS char(32)) AS max_deposit_num,deposit_confirm_num,description,txid_prefix_url,block_height,activity_fee,alias,alias_icon,local_icon").
		Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	// dataPermission := new(models.DataPermission)
	// dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	// table, err := dataPermission.GetDataScope(e.TableName(), table)
	// if err != nil {
	// 	return nil, 0, err
	// }

	var count int
	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

// 更新Coin
func (e *Coin) Update(id int) (update Coin, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("code = ?", id).First(&update).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Coin
func (e *Coin) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("code = ?", id).Delete(&Coin{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Coin) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("code in (?)", id).Delete(&Coin{}).Error; err != nil {
		return
	}
	Result = true
	return
}

//
func (e *Coin) GetAccoutCoins() []Coin {
	var coins []Coin
	err := orm.Eloquent.Table(e.TableName()).
		Select([]string{"name", "type"}).
		Where("find_in_set(name,'USD,USDT,BTC,ETH')").
		Find(&coins)

	if err != nil {
		return nil
	}
	return coins

}
