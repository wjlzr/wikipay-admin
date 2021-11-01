package wallet

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type Profit struct {
	OrderNumber string `json:"orderNumber" gorm:"type:bigint(18);primary_key"` // 订单号
	UserId      int    `json:"userId" gorm:"type:bigint(10);"`                 //
	Amount      string `json:"amount" gorm:"type:decimal(32,8);"`              // 金额
	Profit      string `json:"profit" gorm:"type:decimal(32,8);"`              // 收益
	Interest    string `json:"interest" gorm:"type:decimal(10,4);"`            // 年化收益率
	CreateAt    int    `json:"createAt" gorm:"type:bigint(13);"`               // 创建时间
	DataScope   string `json:"dataScope" gorm:"-"`
	Params      string `json:"params"  gorm:"-"`
}

func (Profit) TableName() string {
	return "profit"
}

// 创建Profit
func (e *Profit) Create() (Profit, error) {
	var doc Profit
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Profit
func (e *Profit) Get() (Profit, error) {
	var doc Profit
	table := orm.Eloquent.Table(e.TableName())

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Profit带分页
func (e *Profit) GetPage(pageSize, pageIndex int, info string) ([]Profit, int, error) {
	var doc []Profit
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
		where = fmt.Sprintf("concat(user_id,order_number) like '%s%s%s'", "%", info, "%")
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

// 更新Profit
func (e *Profit) Update(id int) (update Profit, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", id).First(&update).Error; err != nil {
		return
	}
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Profit
func (e *Profit) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", id).Delete(&Profit{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Profit) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number in (?)", id).Delete(&Profit{}).Error; err != nil {
		return
	}
	Result = true
	return
}
