package wallet

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"

	"wikipay-admin/tools"
)

type WithdrawTemp struct {
	TransactionNumber string `json:"transactionNumber" gorm:"type:bigint(18);primary_key"` // 流水号
	UserId            int    `json:"userId" gorm:"type:bigint(10);"`                       // 用户序列号
	Coin              string `json:"coin" gorm:"type:varchar(20);"`                        // 币种
	Amount            string `json:"amount" gorm:"type:decimal(32,8);"`                    // 数量
	ActivityAmount    string `json:"activityAmount" gorm:"type:decimal(32,8);"`            // 奖励数量
	Usd               string `json:"usd" gorm:"type:decimal(10,2);"`                       // 币当前的美元价格
	ToAddress         string `json:"toAddress" gorm:"type:varchar(255);"`                  // 提现到的地址
	Fee               string `json:"fee" gorm:"type:decimal(10,2);"`                       // 手续费
	Status            int    `json:"status" gorm:"type:tinyint(1);"`                       // 1、未支付 2、完成 3、已过期  4、取消
	Comment           string `json:"comment" gorm:"type:varchar(200);"`                    // 备注
	CreateAt          int    `json:"createAt" gorm:"type:bigint(13);"`                     // 创建时间
	RealUsd           string `json:"realUsd" gorm:"type:decimal(16,4);"`                   //
	AccountType       int    `json:"accountType" gorm:"type:tinyint(1);"`                  //
	DataScope         string `json:"dataScope" gorm:"-"`
	Params            string `json:"params"  gorm:"-"`
}

func (WithdrawTemp) TableName() string {
	return "withdraw_temp"
}

// 创建WithdrawTemp
func (e *WithdrawTemp) Create() (WithdrawTemp, error) {
	var doc WithdrawTemp

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取WithdrawTemp
func (e *WithdrawTemp) Get() (WithdrawTemp, error) {
	var doc WithdrawTemp

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取WithdrawTemp带分页
func (e *WithdrawTemp) GetPage(pageSize, pageIndex int, info string) ([]WithdrawTemp, int, error) {
	var doc []WithdrawTemp

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
		where = fmt.Sprintf("concat(transaction_number) like '%s%s%s'", "%", info, "%")
	}

	if err := table.Where(where).
		Order("create_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Count(&count)
	return doc, count, nil
}

//更新
func (e *WithdrawTemp) Update(id int) (update WithdrawTemp, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_number = ?", id).First(&update).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

//删除
func (e *WithdrawTemp) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_number = ?", id).Delete(&WithdrawTemp{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *WithdrawTemp) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_number in (?)", id).Delete(&WithdrawTemp{}).Error; err != nil {
		return
	}
	Result = true
	return
}
