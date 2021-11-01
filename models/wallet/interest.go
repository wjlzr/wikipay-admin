package wallet

import (
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type Interest struct {
	Id        int    `json:"id" gorm:"type:int(10);primary_key"`  // 序列号
	Interest  string `json:"interest" gorm:"type:decimal(10,4);"` // 年化收益率
	CreateAt  int64  `json:"createAt" gorm:"type:bigint(13);"`    // 创建时间
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (Interest) TableName() string {
	return "interest"
}

// 创建Interest
func (e *Interest) Create() (Interest, error) {
	var doc Interest
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Interest
func (e *Interest) Get() (Interest, error) {
	var doc Interest
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Interest带分页
func (e *Interest) GetPage(pageSize int, pageIndex int) ([]Interest, int, error) {
	var doc []Interest
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

// 更新Interest
func (e *Interest) Update(id int) (update Interest, err error) {
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Interest
func (e *Interest) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Delete(&Interest{}, id).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Interest) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Interest{}).Error; err != nil {
		return
	}
	Result = true
	return
}
