package im

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type ContactHistory struct {
	Id        int    `json:"id" gorm:"type:int(10) unsigned;primary_key"`      //
	FromId    int64  `json:"fromId" gorm:"type:bigint(10) unsigned zerofill;"` //
	ToId      int64  `json:"toId" gorm:"type:bigint(10) unsigned zerofill;"`   //
	Created   string `json:"created" gorm:"type:timestamp;"`                   //
	Updated   string `json:"updated" gorm:"type:timestamp;"`                   //
	Deleted   string `json:"deleted" gorm:"type:timestamp;"`                   //
	Msg       string `json:"msg" gorm:"type:varchar(200);"`                    // 邀请消息
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (ContactHistory) TableName() string {
	return "contact_history"
}

// 创建ContactHistory
func (e *ContactHistory) Create() (ContactHistory, error) {
	var doc ContactHistory
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取ContactHistory
func (e *ContactHistory) Get() (ContactHistory, error) {
	var doc ContactHistory
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取ContactHistory带分页
func (e *ContactHistory) GetPage(pageSize, pageIndex int, info string) ([]ContactHistory, int, error) {
	var doc []ContactHistory
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
		where = fmt.Sprintf("concat(from_id,to_id) like '%s%s%s'", "%", info, "%")
	}
	if err := table.Where(where).
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

// 更新ContactHistory
func (e *ContactHistory) Update(id int) (update ContactHistory, err error) {
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

// 删除ContactHistory
func (e *ContactHistory) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&ContactHistory{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *ContactHistory) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&ContactHistory{}).Error; err != nil {
		return
	}
	Result = true
	return
}
