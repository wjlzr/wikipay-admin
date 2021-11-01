package wallet

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type TransferTemp struct {
	TransactionNumber string `json:"transactionNumber" gorm:"type:bigint(18);primary_key"` // 流水号
	FromId            int    `json:"fromId" gorm:"type:bigint(13);"`                       // 付款人序列号
	ToId              int    `json:"toId" gorm:"type:bigint(13);"`                         // 收款人序列号
	Money             string `json:"money" gorm:"type:decimal(32,2);"`                     // 转账金额
	Status            int    `json:"status" gorm:"type:tinyint(1);"`                       // 1、未支付 2、已完成  3、已失效   4、已取消
	Source            string `json:"source" gorm:"type:varchar(50);"`                      // 转账来源
	Comment           string `json:"comment" gorm:"type:varchar(100);"`                    // 备注
	CreateAt          int    `json:"createAt" gorm:"type:bigint(13);"`                     // 创建时间
	DataScope         string `json:"dataScope" gorm:"-"`
	Params            string `json:"params"  gorm:"-"`
}

//
func (TransferTemp) TableName() string {
	return "transfer_temp"
}

// 创建TransferTemp
func (e *TransferTemp) Create() (TransferTemp, error) {
	var doc TransferTemp
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取TransferTemp
func (e *TransferTemp) Get() (TransferTemp, error) {
	var doc TransferTemp
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取TransferTemp带分页
func (e *TransferTemp) GetPage(pageSize, pageIndex int, info string) ([]TransferTemp, int, error) {
	var doc []TransferTemp
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
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Count(&count)
	return doc, count, nil
}

// 更新TransferTemp
func (e *TransferTemp) Update(id int64) (update TransferTemp, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_number = ?", id).First(&update).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除TransferTemp
func (e *TransferTemp) Delete(id int64) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_number = ?", id).Delete(&TransferTemp{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *TransferTemp) BatchDelete(id []int64) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("transaction_number in (?)", id).Delete(&TransferTemp{}).Error; err != nil {
		return
	}
	Result = true
	return
}
