package user

import (
	_ "time"
	orm "wikipay-admin/database"
)

type Error struct {
	Code      int    `json:"code" gorm:"type:int(7);primary_key"`         // 编码
	ZhCN      string `json:"zhCN" gorm:"column:zh-CN;type:varchar(300);"` // 中文
	ZhHK      string `json:"zhHK" gorm:"column:zh-HK;type:varchar(300);"` // 香港
	ZhTW      string `json:"zhTW" gorm:"column:zh-TW;type:varchar(300);"` // 台湾
	En        string `json:"en" gorm:"type:varchar(300);"`                // 英文
	Vi        string `json:"vi" gorm:"type:varchar(300);"`                // 越南语
	Th        string `json:"th" gorm:"type:varchar(300);"`                // 泰文
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

//
func (Error) TableName() string {
	return "error"
}

// 创建Error
func (e *Error) Create() (Error, error) {
	var doc Error
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取
func (e *Error) Get() (Error, error) {
	var doc Error

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, e.Code).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Error带分页
func (e *Error) GetPage(pageSize int, pageIndex int) ([]Error, int, error) {
	var doc []Error

	table := orm.Eloquent.Select("*").Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	// dataPermission := new(models.DataPermission)
	// dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	// table, err := dataPermission.GetDataScope(e.TableName(), table)
	// if err != nil {
	// 	return nil, 0, err
	// }

	var count int
	if err := table.Order("code DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).
		Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

// 更新Error
func (e *Error) Update(id int) (update Error, err error) {
	if err = orm.Eloquent.Table(e.TableName()).First(&update, "code = ?", id).Error; err != nil {
		return
	}
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Error
func (e *Error) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Delete(&Error{}, "code = ?", id).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Error) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("code in (?)", id).Delete(&Error{}).Error; err != nil {
		return
	}
	Result = true
	return
}
