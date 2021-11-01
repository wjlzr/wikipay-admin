package wallet

import (
	"fmt"
	orm "wikipay-admin/database"
	"wikipay-admin/models"

	_ "time"
	"wikipay-admin/tools"
)

//
type Transaction struct {
	Id          int    `json:"id" gorm:"type:int(11);primary_key"`   //
	OrderNumber string `json:"orderNumber" gorm:"type:bigint(18);"`  //
	FromId      int    `json:"fromId" gorm:"type:bigint(10);"`       //
	ToId        int    `json:"toId" gorm:"type:bigint(10);"`         //
	TxHash      string `json:"txHash" gorm:"type:varchar(200);"`     //
	FromMoney   string `json:"fromMoney" gorm:"type:decimal(32,8);"` //
	Type        int    `json:"type" gorm:"type:tinyint(1);"`         //
	CreateAt    int    `json:"createAt" gorm:"type:bigint(13);"`     //
	Comment     string `json:"comment" gorm:"type:varchar(100);"`    //
	Trmo        string `json:"trmo" gorm:"type:varchar(1000);"`      //
	FromCoin    string `json:"fromCoin" gorm:"type:varchar(20);"`    //
	ToCoin      string `json:"toCoin" gorm:"type:varchar(20);"`      //
	ToMoney     string `json:"toMoney" gorm:"type:decimal(32,8);"`   //
	Usd         string `json:"usd" gorm:"type:decimal(16,4);"`       //
	DataScope   string `json:"dataScope" gorm:"-"`
	Params      string `json:"params"  gorm:"-"`
}

//表名
func (Transaction) TableName() string {
	return "transaction"
}

//创建
func (e *Transaction) Create() (Transaction, error) {
	var doc Transaction

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *Transaction) Get() (Transaction, error) {
	var doc Transaction

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "id = ?", e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取带分页
func (e *Transaction) GetPage(pageSize, pageIndex int, info string, req SearchParams) ([]Transaction, int, error) {
	var doc []Transaction

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
		//where = fmt.Sprintf("order_number = %s", info)
		where = fmt.Sprintf("concat(order_number) like '%s%s%s'", "%", info, "%")
	}

	if req.StartTime != "" {
		table = table.Where("create_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		table = table.Where("create_at <= ?", req.EndTime)
	}

	if err := table.Where(where).
		Order("id DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	var count int
	table.Where(where).Count(&count)

	return doc, count, nil
}

//更新
func (e *Transaction) Update(id int) (update Transaction, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id = ?", id).First(&update).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

//删除
func (e *Transaction) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id = ?", id).Delete(&Transaction{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Transaction) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id in (?)", id).Delete(&Transaction{}).Error; err != nil {
		return
	}
	Result = true
	return
}
