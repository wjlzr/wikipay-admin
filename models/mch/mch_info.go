package mch

import (
	"fmt"
	"time"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type MchInfo struct {
	Id        int    `json:"id" gorm:"type:int(10) unsigned;primary_key"` //
	UserId    int    `json:"userId" gorm:"type:bigint(13);"`              // 商户序列号
	MchId     string `json:"mchId" gorm:"type:varchar(32);"`              // 商户号
	Company   string `json:"company" gorm:"type:varchar(128);"`           // 公司名称
	Telephone string `json:"telephone" gorm:"type:varchar(50);"`          // 座机电话
	Website   string `json:"website" gorm:"type:varchar(256);"`           // 公司网址
	Location  string `json:"location" gorm:"type:varchar(100);"`          // 所在地
	Corporate string `json:"corporate" gorm:"type:varchar(64);"`          // 公司法人
	Logo      string `json:"logo" gorm:"type:varchar(256);"`              // 公司logo
	Country   string `json:"country" gorm:"type:varchar(100);"`           // 国家
	Comment   string `json:"comment" gorm:"type:varchar(256);"`           // 备注
	CreateAt  int64  `json:"createAt" gorm:"type:bigint(13);"`            // 创建j时间
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

//账户资产明细
type Account struct {
	Coin        string  `json:"coin"`      //货币代码
	Available   float64 `json:"available"` //可用资金
	Frozen      float64 `json:"frozen"`    //冻结资金
	Total       float64 `json:"total"`     //总额
	Transaction Transaction
}

type Transaction struct {
	Income string `json:"income"` //总收入
	Refund string `json:"refund"` //退款
}

//
func (MchInfo) TableName() string {
	return "mch_info"
}

// 创建MchInfo
func (e *MchInfo) Create() (MchInfo, error) {
	var doc MchInfo
	e.CreateAt = time.Now().UnixNano() / 1e6
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}

	doc = *e
	return doc, nil
}

// 获取MchInfo
func (e *MchInfo) Get() (accountAssets []Account, err error) {

	var mchInfo MchInfo
	if err := orm.Eloquent.Table(e.TableName()).First(&mchInfo, e.Id).Error; err != nil {
		return accountAssets, err
	}

	sql := fmt.Sprintf(`
		SELECT coin,available,frozen,(frozen+available) total FROM account WHERE user_id = %d
	`, mchInfo.UserId)

	if err = orm.Eloquent.Raw(sql).Find(&accountAssets).Error; err != nil {
		return accountAssets, err
	}
	for k, v := range accountAssets {
		if v.Coin == "USD" {
			accountAssets[k].Transaction.Income, err = MchTransaction{}.GetAccountDetail("USD", "SUCCESS", mchInfo.UserId)
			if err != nil {
				continue
			}
			accountAssets[k].Transaction.Refund, err = MchTransaction{}.GetAccountDetail("USD", "REFUND", mchInfo.UserId)
			if err != nil {
				continue
			}
		}
	}

	return accountAssets, nil
}

// 获取MchInfo带分页
func (e *MchInfo) GetPage(pageSize, pageIndex int, info string) ([]MchInfo, int, error) {
	var doc []MchInfo
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
		where = fmt.Sprintf("concat(user_id,mch_id) like '%s%s%s'", "%", info, "%")
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

// 更新MchInfo
func (e *MchInfo) Update(id int) (update MchInfo, err error) {

	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}
	//struct to map update
	m, err := tools.StructToMap(&e)
	if err != nil {
		return MchInfo{}, err
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(m).Error; err != nil {
		return
	}
	return
}

// 删除MchInfo
func (e *MchInfo) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&MchInfo{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *MchInfo) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&MchInfo{}).Error; err != nil {
		return
	}
	Result = true
	return
}
