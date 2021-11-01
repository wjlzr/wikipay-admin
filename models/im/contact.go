package im

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type Contact struct {
	Id        int    `json:"id" gorm:"type:int(11);primary_key"`        // 序列号
	UserId    int64  `json:"userId" gorm:"type:bigint(10);"`            //
	ToUserId  int64  `json:"toUserId" gorm:"type:int(10) unsigned;"`    // 被联系人序列号
	ToComment string `json:"toComment" gorm:"type:varchar(100);"`       // 被联系人备注
	CreateAt  int64  `json:"createAt" gorm:"type:bigint(13) unsigned;"` // 创建时间
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (Contact) TableName() string {
	return "contact"
}

//创建
func (e *Contact) Create() (Contact, error) {
	var doc Contact
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *Contact) Get() (Contact, error) {
	var doc Contact
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Contact带分页
func (e *Contact) GetPage(pageSize, pageIndex int, info string) ([]Contact, int, error) {
	var doc []Contact
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
		where = fmt.Sprintf("concat(user_id,to_user_id) like '%s%s%s'", "%", info, "%")
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

// 更新Contact
func (e *Contact) Update(id int) (update Contact, err error) {
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

// 删除Contact
func (e *Contact) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&Contact{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Contact) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Contact{}).Error; err != nil {
		return
	}
	Result = true
	return
}
