package mch

import (
	_ "time"
	orm "wikipay-admin/database"
)

//
type MchComment struct {
	Id        int    `json:"id" gorm:"type:int(10) unsigned;primary_key"` //
	UserId    int    `json:"userId" gorm:"type:bigint(10);"`              // 用户序列号
	Key       string `json:"key" gorm:"type:varchar(256);"`               // 附言键
	Value     string `json:"value" gorm:"type:varchar(256);"`             // 附言值
	Required  int    `json:"required" gorm:"type:tinyint(1);"`            // 选择项
	Priority  int    `json:"priority" gorm:"type:tinyint(1);"`            // 排序码
	CreateAt  int    `json:"createAt" gorm:"type:bigint(13);"`            // 创建时间
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

//
func (MchComment) TableName() string {
	return "mch_comment"
}

// 创建MchComment
func (e *MchComment) Create() (MchComment, error) {
	var doc MchComment
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取MchComment
func (e *MchComment) Get() (MchComment, error) {
	var doc MchComment
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取MchComment带分页
func (e *MchComment) GetPage(pageSize int, pageIndex int) ([]MchComment, int, error) {
	var doc []MchComment
	table := orm.Eloquent.Select("*").Table(e.TableName())
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

// 更新MchComment
func (e *MchComment) Update(id int) (update MchComment, err error) {
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

// 删除MchComment
func (e *MchComment) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&MchComment{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *MchComment) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&MchComment{}).Error; err != nil {
		return
	}
	Result = true
	return
}
