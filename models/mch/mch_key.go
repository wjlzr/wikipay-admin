package mch

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type MchKey struct {
	Id        int    `json:"id" gorm:"type:bigint(10) unsigned;primary_key"` //
	UserId    int    `json:"userId" gorm:"type:bigint(10);"`                 // 用户序列号
	Key       string `json:"key" gorm:"type:varchar(32);"`                   // 密钥
	CreateAt  int    `json:"createAt" gorm:"type:bigint(13);"`               // 创建时间
	ExpireAt  int    `json:"expireAt" gorm:"type:bigint(13);"`               // 过期时间
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

//
func (MchKey) TableName() string {
	return "mch_key"
}

// 创建MchKey
func (e *MchKey) Create() (MchKey, error) {
	var doc MchKey
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取MchKey
func (e *MchKey) Get() (MchKey, error) {
	var doc MchKey
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取MchKey带分页
func (e *MchKey) GetPage(pageSize, pageIndex int, info string) ([]MchKey, int, error) {
	var doc []MchKey
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
		where = fmt.Sprintf("concat(user_id) like '%s%s%s'", "%", info, "%")
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

// 更新MchKey
func (e *MchKey) Update(id int) (update MchKey, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除MchKey
func (e *MchKey) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&MchKey{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *MchKey) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&MchKey{}).Error; err != nil {
		return
	}
	Result = true
	return
}
