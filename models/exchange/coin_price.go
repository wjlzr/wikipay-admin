package exchange

import (
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type CoinPrice struct {
	Id        int     `json:"id" gorm:"type:int(10) unsigned;primary_key"` // 序列号
	Coin      string  `json:"coin" gorm:"type:varchar(30);"`               // 币种名称
	Type      int     `json:"type" gorm:"type:tinyint(1);"`                // 类型
	Ask       float64 `json:"ask" gorm:"type:decimal(16,2);"`              // 卖价
	Bid       float64 `json:"bid" gorm:"type:decimal(16,2);"`              // 买价
	CreateAt  int64   `json:"createAt" gorm:"type:bigint(13);"`            // 创建时间
	DataScope string  `json:"dataScope" gorm:"-"`
	Params    string  `json:"params"  gorm:"-"`
}

//表名
func (CoinPrice) TableName() string {
	return "coin_price"
}

//创建
func (e *CoinPrice) Create() (CoinPrice, error) {
	var doc CoinPrice
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *CoinPrice) Get() (CoinPrice, error) {
	var doc CoinPrice
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取带分页
func (e *CoinPrice) GetPage(pageSize int, pageIndex int) ([]CoinPrice, int, error) {
	var doc []CoinPrice
	table := orm.Eloquent.Select("*").Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}
	var count int
	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	table.Count(&count)
	return doc, count, nil
}

//更新
func (e *CoinPrice) Update(id int) (update CoinPrice, err error) {
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除CoinPrice
func (e *CoinPrice) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Delete(&CoinPrice{}, id).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *CoinPrice) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&CoinPrice{}).Error; err != nil {
		return
	}
	Result = true
	return
}
