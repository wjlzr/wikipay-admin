package wallet

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type Account struct {
	Id        int    `json:"id" gorm:"type:int(11);primary_key"`                  // 序列号
	UserId    int64  `json:"userId" gorm:"type:bigint(10);"`                      //
	Coin      string `json:"coin" gorm:"type:varchar(20);"`                       // 货币代码
	Available string `json:"available" gorm:"type:decimal(32,8);"`                // 可用资金
	Frozen    string `json:"frozen" gorm:"type:decimal(32,8) unsigned zerofill;"` // 冻结资金
	InAddress string `json:"inAddress" gorm:"type:varchar(100);"`                 // 入金地址
	Tag       string `json:"tag" gorm:"type:varchar(100);"`                       // 入金标签
	CreateAt  int64  `json:"createAt" gorm:"type:bigint(13) unsigned zerofill;"`  // 创建时间
	Type      int    `json:"type" gorm:"type:tinyint(1);"`                        //
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

type AssetStatistic struct {
	Coin  string  `json:"coin" gorm:"type:varchar(20);"` // 货币代码
	Total float64 `json:"total" gorm:"-"`                // 货币代码
}

//表名
func (Account) TableName() string {
	return "account"
}

// 创建Account
func (e *Account) Create() (Account, error) {
	var doc Account
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//
func (e *Account) GetAccountWithAddress(address string) (*Account, error) {
	var doc Account
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "in_address = ?", address).Error; err != nil {
		return &doc, err
	}
	return &doc, nil
}

// 获取Account
func (e *Account) Get() (Account, error) {
	var doc Account
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "id = ?", e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Account带分页
func (e *Account) GetPage(pageSize int, pageIndex int, info string) ([]Account, int, error) {
	var doc []Account
	table := orm.Eloquent.Select("id,user_id,coin,CONCAT(0+CAST(available AS CHAR(32)),'') AS available,CONCAT(0+CAST(frozen AS CHAR(32)),'') AS frozen,in_address,tag,create_at,type").Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var where string
	if info != "" {
		where = fmt.Sprintf("concat(user_id) like '%s%s%s'", "%", info, "%")
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

// 更新Account
func (e *Account) Update(id int) (update Account, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

//删除Account
func (e *Account) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&Account{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Account) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Account{}).Error; err != nil {
		return
	}
	Result = true
	return
}

//资产统计
func (e Account) AssetStatistics(coin string) (total float64, err error) {

	var assetStatistic AssetStatistic

	table := orm.Eloquent.Select("*").Table(e.TableName())

	err = table.Select("any_value(coin) AS coin, IFNULL(any_value(SUM(available + frozen)),0) AS total").
		Where("coin = ?", coin).
		First(&assetStatistic).Error
	return assetStatistic.Total, err
}
