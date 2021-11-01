package user

import (
	orm "wikipay-admin/database"
	"wikipay-admin/models"

	_ "time"
	"wikipay-admin/tools"
)

type Setting struct {
	Id        int    `json:"id" gorm:"type:int(10);primary_key"` //
	Name      string `json:"name" gorm:"type:varchar(100);"`     // 名称
	SetKey    string `json:"setKey" gorm:"type:varchar(100);"`   // 键
	SetValue  string `json:"setValue" gorm:"type:varchar(200);"` // 值
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (Setting) TableName() string {
	return "setting"
}

//创建
func (e *Setting) Create() (Setting, error) {
	var doc Setting

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *Setting) Get() (Setting, error) {
	var doc Setting

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "id = ?", e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

//获取带分页
func (e *Setting) GetPage(pageSize int, pageIndex int) ([]Setting, int, error) {
	var doc []Setting

	table := orm.Eloquent.Select("*").Table(e.TableName())
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
func (e *Setting) Update(id int) (update Setting, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

//删除
func (e *Setting) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&Setting{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Setting) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Setting{}).Error; err != nil {
		return
	}
	Result = true
	return
}
