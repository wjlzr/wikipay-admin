package user

import (
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type Content struct {
	Id        int    `json:"id" gorm:"type:int(10);"`
	Code      int    `json:"code" gorm:"type:int(8);"`                     //
	ZhCN      string `json:"zhCN" gorm:"column:zh-CN;type:varchar(1000);"` // 中文
	ZhHK      string `json:"zhHK" gorm:"column:zh-HK;type:varchar(1000);"` // 香港
	ZhTW      string `json:"zhTW" gorm:"column:zh-TW;type:varchar(1000);"` // 台湾
	En        string `json:"en" gorm:"type:varchar(1000);"`                // 英文
	Vi        string `json:"vi" gorm:"type:varchar(1000);"`                // 越语
	Type      int    `json:"type" gorm:"type:tinyint(1);"`                 // 1、标题  2、内容
	Th        string `json:"th" gorm:"type:varchar(1000);"`                // 泰语
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

//列名
func (Content) TableName() string {
	return "content"
}

//创建
func (e *Content) Create() (Content, error) {
	var doc Content

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *Content) Get() (Content, error) {
	var doc Content

	table := orm.Eloquent.Table(e.TableName())
	if err := table.Find(&doc, e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取带分页
func (e *Content) GetPage(pageSize int, pageIndex int) ([]Content, int, error) {
	var doc []Content
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
func (e *Content) Update(id int) (update Content, err error) {
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

//删除
func (e *Content) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Delete(&Content{}, id).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Content) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where(" id in (?)", id).Delete(&Content{}).Error; err != nil {
		return
	}
	Result = true
	return
}
