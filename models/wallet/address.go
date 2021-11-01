package wallet

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type Address struct {
	Id        int    `json:"id" gorm:"type:int(11);primary_key"`        //
	UserId    int64  `json:"userId" gorm:"type:bigint(10);"`            //
	Coin      string `json:"coin" gorm:"type:varchar(100);"`            // 币种名称
	Address   string `json:"address" gorm:"type:varchar(200);"`         // 地址
	Name      string `json:"name" gorm:"type:varchar(100);"`            // 姓名
	Comment   string `json:"comment" gorm:"type:varchar(200);"`         // 备注
	CreateAt  int64  `json:"createAt" gorm:"type:bigint(13) unsigned;"` // 创建时间
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (Address) TableName() string {
	return "address"
}

// 创建Address
func (e *Address) Create() (Address, error) {
	var doc Address
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Address
func (e *Address) Get() (Address, error) {
	var doc Address
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Address带分页
func (e *Address) GetPage(pageSize, pageIndex int, info string) ([]Address, int, error) {
	var doc []Address
	table := orm.Eloquent.Select("*").Table(e.TableName())
	//数据权限控制(如果不需要数据权限请将此处去掉)
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
		Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	var count int
	table.Count(&count)
	return doc, count, nil
}

// 更新Address
func (e *Address) Update(id int) (update Address, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Address
func (e *Address) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id = ?", id).Delete(&Address{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Address) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id in (?)", id).Delete(&Address{}).Error; err != nil {
		return
	}
	Result = true
	return
}
